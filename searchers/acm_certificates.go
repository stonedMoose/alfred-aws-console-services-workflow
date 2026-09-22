package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type ACMCertificateSearcher struct{}

func (s ACMCertificateSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s ACMCertificateSearcher) fetch(cfg aws.Config) ([]types.CertificateSummary, error) {
	svc := acm.NewFromConfig(cfg)

	var entities []types.CertificateSummary
	nextToken := ""
	for {
		params := &acm.ListCertificatesInput{
			MaxItems: aws.Int32(1000), // get as many as we can
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListCertificates(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.CertificateSummaryList...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s ACMCertificateSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.CertificateSummary) {
	title := ""
	if entity.DomainName != nil {
		title = *entity.DomainName
	}

	subtitleArray := []string{}
	if entity.Status != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.Status)))
	}
	if entity.Type != "" {
		subtitleArray = append(subtitleArray, strings.ToLower(string(entity.Type)))
	}
	if entity.NotAfter != nil {
		subtitleArray = append(subtitleArray, "expires "+entity.NotAfter.Format("2006-01-02"))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	arn := ""
	if entity.CertificateArn != nil {
		arn = *entity.CertificateArn
	}
	if title == "" {
		title = util.GetEndOfArn(arn)
	}

	path := fmt.Sprintf("/acm/home#/certificates/%s", util.GetEndOfArn(arn))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("acm")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", arn, title)
}
