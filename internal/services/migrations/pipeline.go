package migrations

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/dockerprobe"
	"codedock.run/codedock/internal/models"
)

type RuntimeKinds interface {
	Get(ctx context.Context, serviceID string) (*models.ServiceRuntime, error)
}

type StartMigrationResult struct {
	RunID        string `json:"runId"`
	Confirmation string `json:"confirmation"`
}

func (s *Service) StartMigration(ctx context.Context, userID, orgID string, req models.StartMigrationRequest) (*StartMigrationResult, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	if len(req.ContainerIDs) == 0 {
		return nil, fmt.Errorf("select at least one service")
	}
	if req.ProjectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	source, err := s.orgSource(ctx, orgID, req.SourceID)
	if err != nil {
		return nil, err
	}
	if _, _, err := s.targetServer(ctx, orgID, req.TargetServerID); err != nil {
		return nil, err
	}
	if active, err := s.migrations.ActiveRunForSource(ctx, source.ID); err != nil {
		return nil, err
	} else if active != nil {
		return nil, fmt.Errorf("source already has an active migration")
	}
	mode := req.Mode
	if mode == "" {
		mode = models.MigrationModeMove
	}
	projectID := uuid.NewString()
	if err := s.createAdoptProject(ctx, orgID, projectID, req.ProjectName, source); err != nil {
		return nil, err
	}
	token, hash, err := newConfirmationToken()
	if err != nil {
		return nil, err
	}
	selection := models.MigrationSelection{
		ContainerIDs: req.ContainerIDs, KillOriginals: req.KillOriginals,
		ProjectName: req.ProjectName, ImportEnv: req.ImportEnv,
		Overrides: map[string]string{}, Decisions: map[string]string{},
	}
	run := &models.MigrationRun{
		ID:             uuid.NewString(),
		OrganizationID: orgID,
		UserID:         userID,
		SourceID:       source.ID,
		SourceKind:     "external",
		ProjectID:      projectID,
		TargetServerID: req.TargetServerID,
		Mode:           mode,
		Status:         models.MigrationStatusPending,
		Selection:      encodeSelection(selection),
		TokenHash:      hash,
	}
	if err := s.migrations.CreateRun(ctx, run); err != nil {
		return nil, err
	}
	s.launch(run.ID)
	return &StartMigrationResult{RunID: run.ID, Confirmation: token}, nil
}

func (s *Service) StartProjectMove(ctx context.Context, userID, orgID string, req models.StartProjectMoveRequest) (*StartMigrationResult, error) {
	if err := s.requireOrg(ctx, userID, orgID, false); err != nil {
		return nil, err
	}
	project, err := s.projects.GetByOrganization(ctx, req.ProjectID, orgID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}
	if _, _, err := s.targetServer(ctx, orgID, req.TargetServerID); err != nil {
		return nil, err
	}
	services, err := s.appRepo.ListByProject(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	wanted := map[string]bool{}
	for _, id := range req.ServiceIDs {
		wanted[id] = true
	}
	var progress []models.MigrationServiceProgress
	for _, service := range services {
		if len(wanted) > 0 && !wanted[service.ID] {
			continue
		}
		if err := s.rejectClusterService(ctx, service.ID); err != nil {
			return nil, err
		}
		progress = append(progress, models.MigrationServiceProgress{
			Service: service.Name, ServiceID: service.ID,
			ContainerID: service.ContainerID, Adopted: true,
		})
	}
	if len(progress) == 0 {
		return nil, fmt.Errorf("select at least one service")
	}
	mode := req.Mode
	if mode == "" {
		mode = models.MigrationModeMove
	}
	token, hash, err := newConfirmationToken()
	if err != nil {
		return nil, err
	}
	selection := models.MigrationSelection{
		KillOriginals: req.KillOriginals, SourceServerID: project.ServerID,
		Overrides:     map[string]string{}, Decisions: map[string]string{},
	}
	run := &models.MigrationRun{
		ID:             uuid.NewString(),
		OrganizationID: orgID,
		UserID:         userID,
		SourceKind:     "project",
		ProjectID:      req.ProjectID,
		TargetServerID: req.TargetServerID,
		Mode:           mode,
		Status:         models.MigrationStatusPending,
		Selection:      encodeSelection(selection),
		Progress:       encodeProgress(progress),
		TokenHash:      hash,
	}
	if err := s.migrations.CreateRun(ctx, run); err != nil {
		return nil, err
	}
	s.launch(run.ID)
	return &StartMigrationResult{RunID: run.ID, Confirmation: token}, nil
}

func (s *Service) rejectClusterService(ctx context.Context, serviceID string) error {
	if s.runtimes == nil {
		return nil
	}
	runtime, err := s.runtimes.Get(ctx, serviceID)
	if err != nil {
		return err
	}
	if runtime.Target.Kind == "kubernetes" {
		return fmt.Errorf("project moves do not cover kubernetes workloads")
	}
	return nil
}

func (s *Service) launch(runID string) {
	ctx, handle := s.runs.track(runID)
	go func() {
		defer s.runs.release(runID)
		s.execute(ctx, handle, runID)
	}()
}

func (s *Service) execute(ctx context.Context, handle *runHandle, runID string) {
	run, err := s.migrations.GetRun(ctx, runID)
	if err != nil || run == nil {
		return
	}
	executor := &pipelineExecutor{service: s, run: run, handle: handle}
	if err := executor.drive(ctx); err != nil {
		done := context.Background()
		if executor.cancelled {
			s.setStatus(done, runID, models.MigrationStatusCancelled, executor.phase, "cancelled by request")
			s.appendLog(done, runID, "CANCELLED: rolling back target changes")
			executor.rollback(done)
			return
		}
		s.setStatus(done, runID, models.MigrationStatusFailed, executor.phase, err.Error())
		s.appendLog(done, runID, "FAILED: "+err.Error())
		return
	}
}

type pipelineExecutor struct {
	service   *Service
	run       *models.MigrationRun
	handle    *runHandle
	phase     string
	cancelled bool
	source    dockerprobe.Runner
	target    dockerprobe.Runner
	selection models.MigrationSelection
	progress  []models.MigrationServiceProgress
}

func (e *pipelineExecutor) drive(ctx context.Context) error {
	e.selection = decodeSelection(e.run.Selection)
	e.progress = decodeProgress(e.run.Progress)
	sourceRunner, targetRunner, err := e.runners(ctx)
	if err != nil {
		return err
	}
	e.source = sourceRunner
	e.target = targetRunner
	if err := e.checkCancel(ctx); err != nil {
		return err
	}
	if e.run.SourceKind == "external" {
		if err := e.adoptPhase(ctx); err != nil {
			return err
		}
	}
	if err := e.transferPhase(ctx); err != nil {
		return err
	}
	if err := e.deployPhase(ctx); err != nil {
		return err
	}
	if err := e.verifyPhase(ctx); err != nil {
		return err
	}
	return e.finishPhase(ctx)
}

func (e *pipelineExecutor) runners(ctx context.Context) (dockerprobe.Runner, dockerprobe.Runner, error) {
	_, targetRunner, err := e.service.targetServer(ctx, e.run.OrganizationID, e.run.TargetServerID)
	if err != nil {
		return nil, nil, err
	}
	if e.run.SourceKind == "project" {
		sourceServerID := e.selection.SourceServerID
		if sourceServerID == "" {
			project, err := e.service.projects.GetByOrganization(ctx, e.run.ProjectID, e.run.OrganizationID)
			if err != nil {
				return nil, nil, err
			}
			if project == nil {
				return nil, nil, fmt.Errorf("project not found")
			}
			sourceServerID = project.ServerID
		}
		_, sourceRunner, err := e.service.targetServer(ctx, e.run.OrganizationID, sourceServerID)
		if err != nil {
			return nil, nil, err
		}
		return sourceRunner, targetRunner, nil
	}
	source, err := e.service.migrations.GetSource(ctx, e.run.SourceID)
	if err != nil {
		return nil, nil, err
	}
	if source == nil {
		return nil, nil, fmt.Errorf("migration source not found")
	}
	return e.service.sources(source), targetRunner, nil
}

func (e *pipelineExecutor) checkCancel(ctx context.Context) error {
	if ctx.Err() != nil {
		e.cancelled = true
		return ctx.Err()
	}
	run, err := e.service.migrations.GetRun(ctx, e.run.ID)
	if err != nil {
		return err
	}
	if run.CancelRequested {
		e.cancelled = true
		return fmt.Errorf("cancel requested")
	}
	e.run = run
	return nil
}

func (e *pipelineExecutor) saveProgress(ctx context.Context) error {
	e.run.Selection = encodeSelection(e.selection)
	e.run.Progress = encodeProgress(e.progress)
	return e.service.migrations.UpdateRun(ctx, e.run)
}

func (e *pipelineExecutor) adoptPhase(ctx context.Context) error {
	e.phase = "ADOPT"
	e.service.setStatus(ctx, e.run.ID, models.MigrationStatusTransferring, e.phase, "")
	source, err := e.service.migrations.GetSource(ctx, e.run.SourceID)
	if err != nil {
		return err
	}
	containers, err := e.service.selectContainers(ctx, source, e.selection.ContainerIDs)
	if err != nil {
		return err
	}
	byID := map[string]*models.MigrationServiceProgress{}
	for index := range e.progress {
		byID[e.progress[index].ContainerID] = &e.progress[index]
	}
	for _, container := range containers {
		if e.skipped(container.Name, container.ID) {
			e.progress = append(e.progress, models.MigrationServiceProgress{Service: container.Name, ContainerID: container.ID, Skipped: true})
			continue
		}
		entry, ok := byID[container.ID]
		if !ok {
			e.progress = append(e.progress, models.MigrationServiceProgress{Service: container.Name, ContainerID: container.ID, SourceState: container.State})
			entry = &e.progress[len(e.progress)-1]
		}
		if entry.Adopted {
			continue
		}
		if err := e.checkCancel(ctx); err != nil {
			return err
		}
		serviceID, err := e.service.adoptContainer(ctx, e.run.OrganizationID, e.run.ProjectID, e.environment(ctx), container, e.selection.ImportEnv)
		if err != nil {
			entry.Error = err.Error()
			if err := e.saveProgress(ctx); err != nil {
				return err
			}
			return fmt.Errorf("adopt %s: %w", container.Name, err)
		}
		entry.ServiceID = serviceID
		entry.Adopted = true
		entry.Error = ""
		e.service.appendLog(ctx, e.run.ID, "ADOPTED "+container.Name)
		if err := e.saveProgress(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (e *pipelineExecutor) environment(ctx context.Context) string {
	envs, err := e.service.envs.ListByProject(ctx, e.run.ProjectID)
	if err != nil || len(envs) == 0 {
		return ""
	}
	return envs[0].ID
}

func (e *pipelineExecutor) skipped(name, id string) bool {
	for _, skip := range e.selection.Skips {
		if skip == name || skip == id {
			return true
		}
	}
	return false
}

func (e *pipelineExecutor) transferPhase(ctx context.Context) error {
	e.phase = "TRANSFER"
	e.service.setStatus(ctx, e.run.ID, models.MigrationStatusTransferring, e.phase, "")
	for index := range e.progress {
		entry := &e.progress[index]
		if entry.Skipped || entry.Transferred {
			continue
		}
		if e.skipped(entry.Service, entry.ContainerID) {
			entry.Skipped = true
			continue
		}
		if err := e.checkCancel(ctx); err != nil {
			return err
		}
		if err := e.transferService(ctx, entry); err != nil {
			entry.Error = err.Error()
			_ = e.saveProgress(ctx)
			return err
		}
		entry.Transferred = true
		entry.Error = ""
		e.service.appendLog(ctx, e.run.ID, "TRANSFERRED "+entry.Service)
		if err := e.saveProgress(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (e *pipelineExecutor) transferService(ctx context.Context, entry *models.MigrationServiceProgress) error {
	image, volumes, err := e.serviceImage(ctx, entry)
	if err != nil {
		return err
	}
	if image != "" {
		if err := dockerprobe.TransferImage(ctx, e.source, e.target, image, func(line string) {
			e.service.appendLog(ctx, e.run.ID, line)
		}); err != nil {
			return err
		}
	}
	for _, volume := range volumes {
		targetName := volume
		if override, ok := e.selection.Overrides["volume:"+volume]; ok && override != "" {
			targetName = override
		}
		resolved, err := e.resolveVolumeTarget(ctx, volume, targetName)
		if err != nil {
			return err
		}
		if resolved == "" {
			continue
		}
		if err := dockerprobe.TransferVolume(ctx, e.source, e.target, volume, resolved, func(line string) {
			e.service.appendLog(ctx, e.run.ID, line)
		}); err != nil {
			return err
		}
		entry.CopiedVolumes = append(entry.CopiedVolumes, resolved)
	}
	return nil
}

func (e *pipelineExecutor) serviceImage(ctx context.Context, entry *models.MigrationServiceProgress) (string, []string, error) {
	if e.run.SourceKind == "project" {
		service, err := e.service.appRepo.GetByID(ctx, entry.ServiceID)
		if err != nil {
			return "", nil, err
		}
		records, err := e.service.volumes.ListByService(ctx, entry.ServiceID)
		if err != nil {
			return "", nil, err
		}
		var volumes []string
		for _, volume := range records {
			if volume.HostPath == "" || strings.Contains(volume.HostPath, "/") {
				continue
			}
			volumes = append(volumes, volume.HostPath)
		}
		return service.ImageRef, volumes, nil
	}
	container, err := dockerprobe.InspectContainer(ctx, e.source, entry.ContainerID)
	if err != nil {
		return "", nil, fmt.Errorf("re-inspect %s: %w", entry.Service, err)
	}
	var volumes []string
	for _, volume := range container.Volumes {
		if volume.Type == "volume" && volume.Name != "" {
			volumes = append(volumes, volume.Name)
		}
	}
	return container.Image, volumes, nil
}

func (e *pipelineExecutor) resolveVolumeTarget(ctx context.Context, sourceName, targetName string) (string, error) {
	if decision, ok := e.selection.Decisions["volume:"+sourceName]; ok {
		switch {
		case decision == "use-existing":
			return "", nil
		case decision == "overwrite":
			return targetName, nil
		case len(decision) > 7 && decision[:7] == "rename:":
			return decision[7:], nil
		}
	}
	existing, err := existingVolumes(ctx, e.target)
	if err != nil {
		return "", err
	}
	if !existing[targetName] {
		return targetName, nil
	}
	hasData, err := dockerprobe.VolumeHasData(ctx, e.target, targetName)
	if err != nil {
		return "", err
	}
	if !hasData {
		return targetName, nil
	}
	answer, err := e.askVolume(ctx, sourceName, targetName)
	if err != nil {
		return "", err
	}
	switch answer.optionID {
	case "use-existing":
		e.selection.Decisions["volume:"+sourceName] = "use-existing"
		return "", nil
	case "overwrite":
		e.selection.Decisions["volume:"+sourceName] = "overwrite"
		return targetName, nil
	case "rename":
		e.selection.Decisions["volume:"+sourceName] = "rename:" + answer.value
		return answer.value, nil
	default:
		return "", fmt.Errorf("unknown volume decision")
	}
}

func (e *pipelineExecutor) askVolume(ctx context.Context, sourceName, targetName string) (promptAnswer, error) {
	prompt := newPrompt(models.MigrationPromptVolumeConflict, targetName,
		fmt.Sprintf("target already holds data in volume %s for source volume %s", targetName, sourceName),
		[]models.MigrationPromptOption{
			{ID: "overwrite", Label: "Overwrite", Description: "Replace target data with the source copy"},
			{ID: "use-existing", Label: "Keep target", Description: "Skip this volume and keep target data"},
			{ID: "rename", Label: "Rename", Description: "Copy into a new volume name"},
		})
	return e.ask(ctx, prompt)
}

func (e *Service) logWriter(runID string) *runLogWriter {
	return &runLogWriter{service: e, runID: runID}
}

type runLogWriter struct {
	service *Service
	runID   string
	buffer  []byte
}

func (w *runLogWriter) Write(data []byte) (int, error) {
	w.buffer = append(w.buffer, data...)
	for {
		index := bytes.IndexByte(w.buffer, '\n')
		if index < 0 {
			break
		}
		line := string(w.buffer[:index])
		w.buffer = w.buffer[index+1:]
		w.service.appendLog(context.Background(), w.runID, line)
	}
	return len(data), nil
}
