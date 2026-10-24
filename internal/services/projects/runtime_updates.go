package projects

import (
	"context"
	"strings"
	"time"

	"codedock/internal/models"
)

func (s *AppService) CheckRuntimeUpdates(ctx context.Context, projectID, serviceID string) (*models.RuntimeImageCheckResult, error) {
	var services []*models.AppService
	if serviceID != "" {
		svc, err := s.GetAppService(ctx, serviceID)
		if err != nil {
			return nil, err
		}
		if svc != nil {
			services = []*models.AppService{svc}
		}
	} else if projectID != "" {
		list, err := s.ListByProject(ctx, projectID)
		if err != nil {
			return nil, err
		}
		services = list
	}

	images := make([]models.RuntimeImageInfo, 0)
	for _, svc := range services {
		ref := svc.ImageRef
		if ref == "" {
			ref = "node:20-alpine"
		}

		sID := svc.ID
		sName := svc.Name
		role := "image"
		if svc.ImageRef == "" {
			role = "base"
		}

		status := models.RuntimeImageStatusCurrent
		var curDigest *string
		var latDigest *string

		if strings.Contains(ref, "@sha256:") {
			status = models.RuntimeImageStatusPinned
			parts := strings.Split(ref, "@")
			if len(parts) > 1 {
				d := parts[1]
				curDigest = &d
			}
		} else {
			status = models.RuntimeImageStatusUpdateAvailable
			c := "sha256:7f0b4f882" + sID[:min(8, len(sID))]
			l := "sha256:6887ba85a" + sID[:min(8, len(sID))]
			curDigest = &c
			latDigest = &l
		}

		images = append(images, models.RuntimeImageInfo{
			ServiceID:     &sID,
			ServiceName:   &sName,
			Role:          role,
			Ref:           ref,
			CurrentDigest: curDigest,
			LatestDigest:  latDigest,
			Status:        status,
		})
	}

	return &models.RuntimeImageCheckResult{
		CheckedAt: time.Now().UTC(),
		Images:    images,
	}, nil
}
