package servers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func randomTokenHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

var (
	workloadStoreMu sync.Mutex
	workloadStore   = make(map[string][]ManagedWorkload)
)

func (h *Handler) GetManagedWorkloads(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	workloadStoreMu.Lock()
	defer workloadStoreMu.Unlock()

	items, exists := workloadStore[id]
	if !exists {
		items = []ManagedWorkload{
			{
				ID:            "wk-cron-cleaner",
				Name:          "Nightly Maintenance Worker",
				State:         "running",
				RestartPolicy: "always",
				Source:        "system",
				Manageable:    true,
			},
		}
		workloadStore[id] = items
	}

	return utils.Success(c, "Workloads retrieved", ManagedWorkloadsResponse{
		Workloads: items,
		Truncated: false,
	})
}

func (h *Handler) CreateManagedWorkload(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req CreateManagedWorkloadRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	workloadStoreMu.Lock()
	defer workloadStoreMu.Unlock()

	item := ManagedWorkload{
		ID:            fmt.Sprintf("wk-%s", randomTokenHex(6)),
		Name:          req.Name,
		State:         "running",
		RestartPolicy: req.RestartPolicy,
		Source:        "manual",
		Manageable:    true,
	}

	workloadStore[id] = append(workloadStore[id], item)
	return utils.Success(c, "Workload created", item)
}

func (h *Handler) ControlManagedWorkload(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req ControlManagedWorkloadRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	workloadStoreMu.Lock()
	defer workloadStoreMu.Unlock()

	var updated *ManagedWorkload
	items := workloadStore[id]
	newItems := make([]ManagedWorkload, 0, len(items))
	for _, it := range items {
		if it.ID == req.WorkloadID {
			if req.Action == "delete" {
				continue
			}
			if req.Action == "stop" {
				it.State = "stopped"
			} else if req.Action == "start" || req.Action == "restart" {
				it.State = "running"
			}
			updated = &it
		}
		newItems = append(newItems, it)
	}
	workloadStore[id] = newItems

	return utils.Success(c, "Workload controlled", ControlManagedWorkloadResult{
		Ok:       true,
		Workload: updated,
	})
}

func (h *Handler) GetManagedWorkloadLogs(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	logs := fmt.Sprintf("[%s] process started with PID 14209\n[%s] listening for scheduled job triggers\n[%s] job execution succeeded (duration: 412ms)\n[%s] waiting for next interval...",
		time.Now().Add(-10*time.Minute).Format(time.RFC3339),
		time.Now().Add(-9*time.Minute).Format(time.RFC3339),
		time.Now().Add(-4*time.Minute).Format(time.RFC3339),
		time.Now().Add(-1*time.Minute).Format(time.RFC3339),
	)

	return utils.Success(c, "Workload logs retrieved", ManagedLogsResponse{
		Logs:      logs,
		Truncated: false,
	})
}
