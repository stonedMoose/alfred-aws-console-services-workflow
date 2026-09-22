package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type RDSDatabaseSearcher struct{}

func (s RDSDatabaseSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "rds_databases", s.fetch, s.addToWorkflow)
}

func (RDSDatabaseSearcher) fetch(cfg aws.Config) ([]types.DBInstance, error) {
	client := rds.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.DBInstance, *string, error) {
		resp, err := client.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{
			MaxRecords: aws.Int32(100),
			Marker:     awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DBInstances, resp.Marker, nil
	})
}

func (RDSDatabaseSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, instance types.DBInstance) {
	id := aws.ToString(instance.DBInstanceIdentifier)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       id,
		ConsolePath: "/rds/home#database:id=" + id + ";is-cluster=false",
		ServiceID:   "rds",
		ID:          aws.ToString(instance.DBInstanceArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		joinKnown(" ", aws.ToString(instance.Engine), aws.ToString(instance.EngineVersion)),
		aws.ToString(instance.DBInstanceClass),
		databaseNameDetail(instance, id),
	))
}

// databaseNameDetail mentions the database name only when it adds to the
// instance identifier.
func databaseNameDetail(instance types.DBInstance, instanceID string) string {
	name := aws.ToString(instance.DBName)
	if name == instanceID {
		return ""
	}
	return name
}
