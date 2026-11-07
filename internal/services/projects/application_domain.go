package projects

import (
	"context"
	"fmt"
)

func (s *EnvironmentService) ValidateApplicationDomain(ctx context.Context, serviceID, domain string) (string, error) {
	canonical, err := canonicalDomainName(domain)
	if err != nil {
		return "", err
	}
	domains, err := s.domainRepo.ListAll(ctx)
	if err != nil {
		return "", fmt.Errorf("check domain ownership: %w", err)
	}
	for _, current := range domains {
		if current.DomainName == canonical && current.ServiceID != serviceID {
			return "", fmt.Errorf("domain is already assigned to another service")
		}
	}
	return canonical, nil
}
