package dns

import (
	"context"
	"fmt"
	"log/slog"
)

type Service struct {
	repo        Repository
	providerSvc *ProviderService
}

func NewService(repo Repository, providerSvc *ProviderService) *Service {
	return &Service{repo: repo, providerSvc: providerSvc}
}

func (s *Service) CreateRecord(ctx context.Context, req *CreateDNSRecordRequest) (*DNSRecord, error) {
	if req.DomainName == "" || req.RecordType == "" || req.RecordName == "" || req.RecordValue == "" {
		return nil, ErrFieldsRequired
	}

	record := &DNSRecord{
		DomainName:  req.DomainName,
		RecordType:  req.RecordType,
		RecordName:  req.RecordName,
		RecordValue: req.RecordValue,
		TTL:         req.TTL,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return nil, err
	}

	if record.RecordName != "" && s.providerSvc != nil {
		domain := record.DomainName
		if record.RecordName != "@" {
			domain = fmt.Sprintf("%s.%s", record.RecordName, record.DomainName)
		}
		if _, err := s.providerSvc.ProvisionRecord(ctx, domain, record.RecordType, record.RecordValue); err != nil {
			if deleteErr := s.repo.Delete(ctx, record.ID); deleteErr != nil {
				slog.Error("failed to roll back DNS record after provisioning failure", "record_id", record.ID, "error", deleteErr)
			}
			return nil, fmt.Errorf("provision DNS record: %w", err)
		}
	}

	return record, nil
}

func (s *Service) ListByDomain(ctx context.Context, domainName string) ([]*DNSRecord, error) {
	return s.repo.ListByDomain(ctx, domainName)
}

func (s *Service) UpdateRecord(ctx context.Context, id string, req *UpdateDNSRecordRequest) (*DNSRecord, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.RecordType != "" {
		record.RecordType = req.RecordType
	}
	if req.RecordName != "" {
		record.RecordName = req.RecordName
	}
	if req.RecordValue != "" {
		record.RecordValue = req.RecordValue
	}
	if req.TTL > 0 {
		record.TTL = req.TTL
	}

	if err := s.repo.Update(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *Service) DeleteRecord(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

