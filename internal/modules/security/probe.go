package security

import (
	"fmt"
	"regexp"
	"strings"
)

var runtimeCommands = map[string]string{
	"node":   "node --version",
	"bun":    "bun --version",
	"deno":   "deno --version",
	"python": "python3 --version",
	"php":    "php --version",
	"ruby":   "ruby --version",
	"java":   "java -version",
	"dotnet": "dotnet --list-runtimes",
	"go":     "go version",
}

var versionExtractRegex = regexp.MustCompile(`(\d+\.\d+\.\d+)`)

func GenerateRuntimeProbeScript(runtimes []string) string {
	if len(runtimes) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("unset NODE_OPTIONS BUN_OPTIONS RUBYOPT JAVA_TOOL_OPTIONS JDK_JAVA_OPTIONS _JAVA_OPTIONS PYTHONSTARTUP PYTHONPATH; ")

	for _, rt := range runtimes {
		rt = strings.ToLower(strings.TrimSpace(rt))
		cmd, ok := runtimeCommands[rt]
		if !ok {
			continue
		}

		// check binary existence first, then echo sentinel and execute version command capped at 1024 bytes
		binary := strings.Fields(cmd)[0]
		sb.WriteString(fmt.Sprintf(
			"if command -v %s >/dev/null 2>&1; then printf '\\nCODEDOCK_RUNTIME_%s\\n'; %s 2>&1 | head -c 1024; fi; ",
			binary, rt, cmd,
		))
	}

	return strings.TrimSpace(sb.String())
}

func ParseRuntimeProbeOutput(raw string) map[string]string {
	results := make(map[string]string)
	lines := strings.Split(raw, "\n")

	var currentRuntime string
	var currentOutput strings.Builder

	flush := func() {
		if currentRuntime != "" {
			text := currentOutput.String()
			matches := versionExtractRegex.FindStringSubmatch(text)
			if len(matches) > 1 {
				results[currentRuntime] = matches[1]
			}
			currentOutput.Reset()
			currentRuntime = ""
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "CODEDOCK_RUNTIME_") || strings.HasPrefix(trimmed, "OPENSHIP_RUNTIME_") {
			flush()
			if strings.HasPrefix(trimmed, "CODEDOCK_RUNTIME_") {
				currentRuntime = strings.TrimPrefix(trimmed, "CODEDOCK_RUNTIME_")
			} else {
				currentRuntime = strings.TrimPrefix(trimmed, "OPENSHIP_RUNTIME_")
			}
			continue
		}

		if currentRuntime != "" {
			currentOutput.WriteString(line)
			currentOutput.WriteString(" ")
		}
	}

	flush()
	return results
}
