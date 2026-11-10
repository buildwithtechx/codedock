package deploy

import (
	"bytes"
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
)

const (
	ResourceIssueContainerOOM     = "container-oom"
	ResourceIssueMemoryAllocation = "memory-allocation"
	ResourceIssueProcessHeap      = "process-heap"
	ResourceIssueStorageFull      = "storage-full"
	ResourceIssueStorageLow       = "storage-low"
	ResourceIssueProcessKilled    = "process-killed"
)

var (
	exitCodeRegex    = regexp.MustCompile(`(?i)(?:returned a non-zero code|exit(?:ed)? with(?: exit)? code|exit code)\s*:?\s*(\d+)`)
	storageFullRegex = regexp.MustCompile(`(?i)(no space left on device|\bENOSPC\b)`)
	processHeapRegex = regexp.MustCompile(`(?i)(JavaScript heap out of memory|Ineffective mark-compacts near heap limit|Reached heap limit Allocation failed)`)
	memoryAllocRegex = regexp.MustCompile(`(?i)(\bENOMEM\b|cannot allocate memory|Fatal process out of memory|fatal error:\s*(?:runtime:\s*)?out of memory)`)
)

type ResourceFailureError struct {
	Issue   string
	Message string
}

func (e *ResourceFailureError) Error() string {
	return e.Message
}

func NewResourceFailureError(issue, message string) *ResourceFailureError {
	return &ResourceFailureError{Issue: issue, Message: message}
}

func ExtractBuildExitCode(output string) int {
	matches := exitCodeRegex.FindAllStringSubmatch(output, -1)
	for _, m := range matches {
		if len(m) > 1 {
			if code, err := strconv.Atoi(m[1]); err == nil && code != 0 {
				return code
			}
		}
	}
	return 0
}

func ClassifyResourceIssue(output string) string {
	if storageFullRegex.MatchString(output) {
		return ResourceIssueStorageFull
	}
	if processHeapRegex.MatchString(output) {
		return ResourceIssueProcessHeap
	}
	if memoryAllocRegex.MatchString(output) {
		return ResourceIssueMemoryAllocation
	}
	if ExtractBuildExitCode(output) == 137 {
		return ResourceIssueProcessKilled
	}
	return ""
}

type ContainerStabilityVerdict struct {
	OOMKilled     bool
	ExitCode      int
	ResourceIssue string
	LogTail       string
	Reason        string
}

func DiagnoseContainerFailure(ctx context.Context, cm *ContainerManager, containerIDOrName string) ContainerStabilityVerdict {
	if cm == nil {
		return ContainerStabilityVerdict{}
	}
	inspect, err := cm.Inspect(ctx, containerIDOrName)
	if err != nil {
		return ContainerStabilityVerdict{}
	}

	tail := readContainerLogTail(ctx, cm, containerIDOrName, 30)
	combined := tail
	if inspect.State != nil && inspect.State.Error != "" {
		combined = combined + "\n" + inspect.State.Error
	}

	oomKilled := inspect.State != nil && inspect.State.OOMKilled
	exitCode := 0
	if inspect.State != nil {
		exitCode = inspect.State.ExitCode
	}

	if oomKilled {
		return ContainerStabilityVerdict{
			OOMKilled:     true,
			ExitCode:      exitCode,
			ResourceIssue: ResourceIssueContainerOOM,
			LogTail:       tail,
			Reason:        "Docker confirmed container was OOM-killed (State.OOMKilled)",
		}
	}

	if issue := ClassifyResourceIssue(combined); issue != "" {
		return ContainerStabilityVerdict{
			OOMKilled:     false,
			ExitCode:      exitCode,
			ResourceIssue: issue,
			LogTail:       tail,
			Reason:        issue,
		}
	}

	if exitCode == 137 {
		return ContainerStabilityVerdict{
			OOMKilled:     false,
			ExitCode:      137,
			ResourceIssue: ResourceIssueProcessKilled,
			LogTail:       tail,
			Reason:        "process terminated with SIGKILL (exit code 137)",
		}
	}

	return ContainerStabilityVerdict{
		OOMKilled: false,
		ExitCode:  exitCode,
		LogTail:   tail,
	}
}

func readContainerLogTail(ctx context.Context, cm *ContainerManager, containerIDOrName string, lines int) string {
	reader, err := cm.dockerClient.ContainerLogs(ctx, containerIDOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       strconv.Itoa(lines),
	})
	if err != nil {
		return ""
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(reader)
	raw := buf.String()
	clean := stripDockerLogHeaders(raw)
	return strings.TrimSpace(clean)
}

func stripDockerLogHeaders(s string) string {
	var b strings.Builder
	for len(s) > 0 {
		if len(s) >= 8 && (s[0] == 1 || s[0] == 2) && s[1] == 0 && s[2] == 0 && s[3] == 0 {
			s = s[8:]
			continue
		}
		idx := strings.IndexByte(s, '\n')
		if idx == -1 {
			b.WriteString(s)
			break
		}
		b.WriteString(s[:idx+1])
		s = s[idx+1:]
	}
	return b.String()
}
