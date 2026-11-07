package hetzner

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"codedock/internal/models"
)

const (
	StatusRunning = "running"
	StatusOff     = "off"
)

type Adapter struct {
	client *Client
}

func NewProvider(token string) *Adapter {
	return &Adapter{client: NewClient(token)}
}

func NewProviderWithBaseURL(token, baseURL string) *Adapter {
	return &Adapter{client: NewClientWithBaseURL(token, baseURL)}
}

func (a *Adapter) CreateServer(ctx context.Context, spec models.ManagedServerSpec) (models.ManagedInstance, error) {
	created, err := a.client.CreateServer(ctx, spec.Name, spec.ServerType, spec.Image, spec.Region, spec.SSHKeyName, spec.UserData, spec.Labels)
	if err != nil {
		return models.ManagedInstance{}, err
	}
	return parseInstance(created), nil
}

func (a *Adapter) GetServer(ctx context.Context, externalID string) (models.ManagedInstance, error) {
	if externalID == "" {
		return models.ManagedInstance{}, fmt.Errorf("hetzner server id is required")
	}
	found, err := a.client.GetServer(ctx, externalID)
	if err != nil {
		return models.ManagedInstance{}, err
	}
	if found == nil {
		return models.ManagedInstance{}, fmt.Errorf("hetzner server %s not found", externalID)
	}
	return parseInstance(found), nil
}

func (a *Adapter) FindServerByLabel(ctx context.Context, selector string) (models.ManagedInstance, error) {
	found, err := a.client.FindServerByLabel(ctx, selector)
	if err != nil {
		return models.ManagedInstance{}, err
	}
	if found == nil {
		return models.ManagedInstance{}, nil
	}
	return parseInstance(found), nil
}

func (a *Adapter) DeleteServer(ctx context.Context, externalID string) error {
	if externalID == "" {
		return fmt.Errorf("hetzner server id is required")
	}
	return a.client.DeleteServer(ctx, externalID)
}

func (a *Adapter) PowerOn(ctx context.Context, externalID string) error {
	if externalID == "" {
		return fmt.Errorf("hetzner server id is required")
	}
	return a.client.PowerOn(ctx, externalID)
}

func (a *Adapter) PowerOff(ctx context.Context, externalID string) error {
	if externalID == "" {
		return fmt.Errorf("hetzner server id is required")
	}
	_, err := a.client.call(ctx, http.MethodPost, "/servers/"+externalID+"/actions/poweroff", map[string]any{})
	return err
}

func (a *Adapter) ChangeType(ctx context.Context, externalID, serverType string, upgradeDisk bool) error {
	if externalID == "" {
		return fmt.Errorf("hetzner server id is required")
	}
	return a.client.ChangeType(ctx, externalID, serverType, upgradeDisk)
}

func (a *Adapter) WaitForStatus(ctx context.Context, externalID string, wanted []string, timeout time.Duration) (models.ManagedInstance, error) {
	deadline := time.Now().Add(timeout)
	for {
		instance, err := a.GetServer(ctx, externalID)
		if err != nil {
			return models.ManagedInstance{}, err
		}
		for _, status := range wanted {
			if instance.Status == status {
				return instance, nil
			}
		}
		if time.Now().After(deadline) {
			return instance, fmt.Errorf("hetzner server %s did not reach %v within %s", externalID, wanted, timeout)
		}
		select {
		case <-ctx.Done():
			return instance, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

func (a *Adapter) CreateSSHKey(ctx context.Context, name, publicKey string) error {
	if name == "" || publicKey == "" {
		return fmt.Errorf("hetzner ssh key name and public key are required")
	}
	_, err := a.client.call(ctx, http.MethodPost, "/ssh_keys", map[string]any{
		"name":       name,
		"public_key": publicKey,
		"labels":     map[string]string{"codedock-managed": "true"},
	})
	return err
}

func (a *Adapter) DeleteSSHKey(ctx context.Context, name string) error {
	if name == "" {
		return nil
	}
	decoded, err := a.client.call(ctx, http.MethodGet, "/ssh_keys?name="+name, nil)
	if err != nil {
		return err
	}
	items, _ := decoded["ssh_keys"].([]any)
	for _, item := range items {
		key, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if _, err := a.client.call(ctx, http.MethodDelete, "/ssh_keys/"+stringify(key["id"]), nil); err != nil {
			return err
		}
	}
	return nil
}

func parseInstance(server map[string]any) models.ManagedInstance {
	if server == nil {
		return models.ManagedInstance{}
	}
	instance := models.ManagedInstance{
		ExternalID: stringify(server["id"]),
		Name:       stringify(server["name"]),
		Status:     stringify(server["status"]),
		Labels:     map[string]string{},
	}
	if publicNet, ok := server["public_net"].(map[string]any); ok {
		if ipv4, ok := publicNet["ipv4"].(map[string]any); ok {
			instance.PublicIP = stringify(ipv4["ip"])
		}
	}
	if serverType, ok := server["server_type"].(map[string]any); ok {
		instance.ServerType = stringify(serverType["name"])
	} else {
		instance.ServerType = stringify(server["server_type"])
	}
	if datacenter, ok := server["datacenter"].(map[string]any); ok {
		if location, ok := datacenter["location"].(map[string]any); ok {
			instance.Region = stringify(location["name"])
		}
	}
	if labels, ok := server["labels"].(map[string]any); ok {
		for key, value := range labels {
			if text, ok := value.(string); ok {
				instance.Labels[key] = text
			}
		}
	}
	if instance.Region == "" {
		instance.Region = instance.Labels["codedock-region"]
	}
	return instance
}

func stringify(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	default:
		return fmt.Sprintf("%v", typed)
	}
}
