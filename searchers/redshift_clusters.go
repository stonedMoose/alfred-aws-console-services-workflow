package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type RedshiftClusterSearcher struct{}

func (s RedshiftClusterSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s RedshiftClusterSearcher) fetch(cfg aws.Config) ([]types.Cluster, error) {
	svc := redshift.NewFromConfig(cfg)

	var entities []types.Cluster
	marker := ""
	for {
		params := &redshift.DescribeClustersInput{
			MaxRecords: aws.Int32(100), // max allowed by this API
		}
		if marker != "" {
			params.Marker = aws.String(marker)
		}
		resp, err := svc.DescribeClusters(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Clusters...)

		if resp.Marker != nil && *resp.Marker != "" {
			marker = *resp.Marker
		} else {
			break
		}
	}

	return entities, nil
}

func (s RedshiftClusterSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Cluster) {
	title := *entity.ClusterIdentifier

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.ClusterStatus)
	subtitleArray = util.AppendString(subtitleArray, entity.NodeType)
	if entity.NumberOfNodes != nil {
		subtitleArray = append(subtitleArray, strconv.Itoa(int(*entity.NumberOfNodes))+" nodes")
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/redshiftv2/home#cluster-details?cluster=%s", url.QueryEscape(title))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("redshift")).
		Valid(true)

	arn := ""
	if entity.ClusterNamespaceArn != nil {
		arn = *entity.ClusterNamespaceArn
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
