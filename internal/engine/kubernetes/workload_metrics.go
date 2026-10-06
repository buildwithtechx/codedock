package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func applyPodMetrics(observation *models.WorkloadObservation, raw string) error {
	var response struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Containers []struct {
				Usage map[string]string `json:"usage"`
			} `json:"containers"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return fmt.Errorf("decode pod metrics: %w", err)
	}
	for i := range observation.Pods {
		found := false
		for _, pod := range response.Items {
			if pod.Metadata.Name != observation.Pods[i].Name {
				continue
			}
			found = len(pod.Containers) > 0
			for _, container := range pod.Containers {
				cpu, err := metricQuantity(container.Usage["cpu"])
				if err != nil {
					return err
				}
				memory, err := metricQuantity(container.Usage["memory"])
				if err != nil {
					return err
				}
				observation.Pods[i].CPU += cpu
				observation.Pods[i].MemoryBytes += int64(memory)
			}
		}
		if !found && observation.Pods[i].Ready {
			return fmt.Errorf("metrics missing for ready pod %s", observation.Pods[i].Name)
		}
	}
	return nil
}
func metricQuantity(raw string) (float64, error) {
	units := []struct {
		Suffix     string
		Multiplier float64
	}{{"Ei", 1 << 60}, {"Pi", 1 << 50}, {"Ti", 1 << 40}, {"Gi", 1 << 30}, {"Mi", 1 << 20}, {"Ki", 1 << 10}, {"E", 1e18}, {"P", 1e15}, {"T", 1e12}, {"G", 1e9}, {"M", 1e6}, {"k", 1e3}, {"m", 1e-3}, {"u", 1e-6}, {"n", 1e-9}}
	multiplier := float64(1)
	for _, unit := range units {
		if strings.HasSuffix(raw, unit.Suffix) {
			raw = strings.TrimSuffix(raw, unit.Suffix)
			multiplier = unit.Multiplier
			break
		}
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil || value < 0 || math.IsNaN(value) || math.IsInf(value, 0) || math.IsInf(value*multiplier, 0) {
		return 0, fmt.Errorf("invalid metric quantity")
	}
	return value * multiplier, nil
}
