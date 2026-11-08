package networking

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"codedock/internal/config"
	"codedock/internal/utils"
)

// DefaultRangeStart and DefaultRangeEnd define the default host port pool.
const (
	DefaultRangeStart = 20000
	DefaultRangeEnd   = 29999
)

// PortRange represents an inclusive range of ports [Start, End].
type PortRange struct {
	Start int
	End   int
}

// DefaultLinuxEphemeralRange is Linux's default outbound ephemeral port range.
var DefaultLinuxEphemeralRange = PortRange{Start: 32768, End: 60999}

// AllocateHostPortOptions controls host port selection and safety constraints.
type AllocateHostPortOptions struct {
	// Preferred port to reuse if free (e.g. previously persisted hostPort).
	Preferred int
	// ReuseOccupiedPreferred treats an occupied preferred port as reusable
	// when the caller has verified target ownership.
	ReuseOccupiedPreferred bool
	// Avoid specifies host ports reserved by other workloads.
	Avoid []int
	// RangeStart overrides the search range lower bound.
	RangeStart int
	// RangeEnd overrides the search range upper bound.
	RangeEnd int
	// EphemeralRange prevents allocating or reusing ports in the kernel's
	// ephemeral outbound range.
	EphemeralRange *PortRange
}

// ParseEphemeralPortRange parses the content of /proc/sys/net/ipv4/ip_local_port_range.
func ParseEphemeralPortRange(raw string) PortRange {
	fields := strings.Fields(raw)
	if len(fields) >= 2 {
		start, err1 := strconv.Atoi(fields[0])
		end, err2 := strconv.Atoi(fields[1])
		if err1 == nil && err2 == nil && start >= 1 && start <= end && end <= 65535 {
			return PortRange{Start: start, End: end}
		}
	}
	return DefaultLinuxEphemeralRange
}

// ReadEphemeralPortRange reads the host's ip_local_port_range, falling back to Linux defaults.
func ReadEphemeralPortRange() PortRange {
	data, err := os.ReadFile("/proc/sys/net/ipv4/ip_local_port_range")
	if err != nil {
		return DefaultLinuxEphemeralRange
	}
	return ParseEphemeralPortRange(string(data))
}

// PickHostPort picks the first safe port in [RangeStart, RangeEnd] outside the ephemeral
// range, reusing Preferred when safe and available.
func PickHostPort(occupied map[int]bool, opts AllocateHostPortOptions) (int, error) {
	start := opts.RangeStart
	if start <= 0 {
		start = DefaultRangeStart
	}
	end := opts.RangeEnd
	if end <= 0 {
		end = DefaultRangeEnd
	}
	if start > end {
		start, end = end, start
	}

	avoidSet := make(map[int]bool, len(opts.Avoid))
	for _, p := range opts.Avoid {
		avoidSet[p] = true
	}

	ephemeral := opts.EphemeralRange
	if ephemeral == nil {
		r := ReadEphemeralPortRange()
		ephemeral = &r
	}

	inEphemeral := func(p int) bool {
		return ephemeral != nil && p >= ephemeral.Start && p <= ephemeral.End
	}

	isFree := func(p int) bool {
		return !inEphemeral(p) && !occupied[p] && !avoidSet[p]
	}

	// Reuse preferred if valid, non-ephemeral, and available/reusable
	if opts.Preferred > 0 &&
		!inEphemeral(opts.Preferred) &&
		!avoidSet[opts.Preferred] &&
		(opts.ReuseOccupiedPreferred || !occupied[opts.Preferred]) {
		return opts.Preferred, nil
	}

	for p := start; p <= end; p++ {
		if isFree(p) {
			return p, nil
		}
	}

	excludedInfo := ""
	if ephemeral != nil {
		excludedInfo = fmt.Sprintf(" outside ephemeral range %d-%d", ephemeral.Start, ephemeral.End)
	}
	return 0, fmt.Errorf("no free host port available in %d-%d%s", start, end, excludedInfo)
}

// SuggestAlternativeHostPort suggests a viable replacement host port, excluding specific ports.
func SuggestAlternativeHostPort(occupied map[int]bool, excluded []int, opts AllocateHostPortOptions) (int, error) {
	optsCopy := opts
	optsCopy.Avoid = append(append([]int(nil), opts.Avoid...), excluded...)
	return PickHostPort(occupied, optsCopy)
}

// GetAvailablePort provides legacy allocation within configured Docker port range.
func GetAvailablePort() (int, error) {
	cfg := config.Get()
	start := cfg.Docker.PortStart
	end := cfg.Docker.PortEnd

	if start > end {
		start, end = end, start
	}

	for port := start; port <= end; port++ {
		addr := fmt.Sprintf("0.0.0.0:%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			l.Close()
			return port, nil
		}
	}

	return 0, &utils.NonReportableError{Message: fmt.Sprintf("no available ports found between %d and %d", start, end)}
}
