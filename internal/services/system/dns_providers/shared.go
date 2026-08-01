package dnsproviders

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	"codedock.run/codedock/internal/models"
)

const (
	dnsProviderCloudflare = "cloudflare"
	dnsProviderNamecheap  = "namecheap"
	dnsProviderSpaceship  = "spaceship"
)

func newProviderHTTPClient() *http.Client {
	return &http.Client{Timeout: 15 * time.Second}
}

func getRootDomain(domain string) string {
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	rootDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err == nil {
		return rootDomain
	}

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return domain
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

func splitNamecheapDomain(domain string) (string, string, error) {
	rootDomain := getRootDomain(domain)
	parts := strings.SplitN(rootDomain, ".", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid root domain for namecheap")
	}
	return parts[0], parts[1], nil
}

func checkCloudflareRecordExists(ctx context.Context, client *http.Client, token, zoneID, domain, recordType, value string) (bool, []string, error) {
	matchingIDs := make([]string, 0)
	for page := 1; ; page++ {
		query := url.Values{}
		query.Set("name", domain)
		query.Set("type", recordType)
		query.Set("per_page", "100")
		query.Set("page", fmt.Sprintf("%d", page))
		u := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?%s", url.PathEscape(zoneID), query.Encode())
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return false, nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
			return false, nil, err
		}
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			return false, nil, fmt.Errorf("cloudflare returned status %d", resp.StatusCode)
		}
		var res struct {
			Result []struct {
				ID      string `json:"id"`
				Content string `json:"content"`
			} `json:"result"`
			ResultInfo struct {
				Page       int `json:"page"`
				TotalPages int `json:"total_pages"`
			} `json:"result_info"`
		}
		decodeErr := json.NewDecoder(resp.Body).Decode(&res)
		resp.Body.Close()
		if decodeErr != nil {
			return false, nil, decodeErr
		}
		for _, rec := range res.Result {
			if value == "" || rec.Content == value {
				matchingIDs = append(matchingIDs, rec.ID)
			}
		}
		if res.ResultInfo.TotalPages <= page || len(res.Result) == 0 {
			break
		}
	}
	return len(matchingIDs) > 0, matchingIDs, nil
}

type namecheapHost struct {
	Name    string
	Type    string
	Address string
	MXPref  string
	TTL     string
}

func parseNamecheapHosts(body []byte) ([]namecheapHost, error) {
	decoder := xml.NewDecoder(bytes.NewReader(body))
	hosts := make([]namecheapHost, 0)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		start, ok := token.(xml.StartElement)
		if !ok || !strings.EqualFold(start.Name.Local, "host") {
			continue
		}

		host := namecheapHost{}
		for _, attr := range start.Attr {
			switch strings.ToLower(attr.Name.Local) {
			case "name":
				host.Name = attr.Value
			case "type":
				host.Type = attr.Value
			case "address":
				host.Address = attr.Value
			case "mxpref":
				host.MXPref = attr.Value
			case "ttl":
				host.TTL = attr.Value
			}
		}
		hosts = append(hosts, host)
	}

	return hosts, nil
}

func namecheapRecordExists(body []byte, subDomain, recordType, value string) bool {
	hosts, err := parseNamecheapHosts(body)
	if err != nil {
		return false
	}
	for _, host := range hosts {
		if host.Name == subDomain && host.Type == recordType && (value == "" || host.Address == value) {
			return true
		}
	}
	return false
}

func buildNamecheapSetHostsQuery(parts []string, cfg *models.ServerSettings, hosts []namecheapHost) url.Values {
	values := url.Values{}
	values.Set("ApiUser", cfg.NamecheapAPIUser)
	values.Set("ApiKey", cfg.NamecheapAPIKey)
	values.Set("UserName", cfg.NamecheapAPIUser)
	values.Set("Command", "namecheap.domains.dns.setHosts")
	values.Set("ClientIp", cfg.NamecheapClientIP)
	values.Set("SLD", parts[0])
	values.Set("TLD", parts[1])

	for i, host := range hosts {
		index := i + 1
		values.Set(fmt.Sprintf("HostName%d", index), host.Name)
		values.Set(fmt.Sprintf("RecordType%d", index), host.Type)
		values.Set(fmt.Sprintf("Address%d", index), host.Address)
		if host.TTL != "" {
			values.Set(fmt.Sprintf("TTL%d", index), host.TTL)
		}
		if host.MXPref != "" {
			values.Set(fmt.Sprintf("MXPref%d", index), host.MXPref)
		}
	}

	return values
}
