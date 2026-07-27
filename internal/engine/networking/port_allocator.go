package networking

import (
	"fmt"
	"net"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/utils"
)

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
