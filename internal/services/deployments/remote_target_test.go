package deployments

import (
	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"context"
	"testing"
)

type remoteProjects struct{ repositories.ProjectRepository }

func (remoteProjects) Get(context.Context, string) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{OrganizationID: "owner", ServerID: "server"}, nil
}

type remoteServers struct{ repositories.ServerRepository }

func (remoteServers) GetByID(context.Context, string) (*models.Server, error) {
	return &models.Server{OrganizationID: "foreign"}, nil
}
func TestRemoteDeploymentRejectsServerWhoseOrganizationChanged(t *testing.T) {
	service := &DeploymentService{projectRepo: remoteProjects{}, Servers: remoteServers{}, HostGate: deploy.NewVolumeGate(), sshManager: ssh.NewSSHManager(remoteServers{})}
	if _, _, err := service.dockerTarget(context.Background(), &models.AppService{ProjectID: "project"}); err == nil {
		t.Fatal("foreign remote Docker target accepted")
	}
}
