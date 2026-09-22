package awsworkflow

import "github.com/aws/aws-sdk-go-v2/aws"

// AwsService is a console service, or one of its sub-services, as described
// in the services catalogue.
type AwsService struct {
	Id               string       `yaml:"id"`
	Name             string       `yaml:"name"`
	ShortName        string       `yaml:"short_name"`
	Description      string       `yaml:"description"`
	Url              string       `yaml:"url"`
	HomeID           string       `yaml:"home_id"`
	ExtraSearchTerms []string     `yaml:"extra_search_terms"`
	SubServices      []AwsService `yaml:"sub_services"`
	HasGlobalRegion  bool         `yaml:"has_global_region"`
}

// GetRegion returns the region the service's console lives in: the
// configured one, or none for the services with a single global console.
func (s *AwsService) GetRegion(cfg aws.Config) string {
	if s.HasGlobalRegion {
		return ""
	}
	return cfg.Region
}

func (s *AwsService) HasSubServices() bool {
	return len(s.SubServices) > 0
}

// FindServiceById returns the service of that id within services, or nil.
func FindServiceById(services []AwsService, id string) *AwsService {
	for i := range services {
		if services[i].Id == id {
			return &services[i]
		}
	}
	return nil
}
