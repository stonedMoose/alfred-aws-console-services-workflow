package searchers

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ACMCertificateSearcher struct{}

func (s ACMCertificateSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "acm_certificates", s.fetch, s.addToWorkflow)
}

func (ACMCertificateSearcher) fetch(cfg aws.Config) ([]types.CertificateSummary, error) {
	client := acm.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.CertificateSummary, *string, error) {
		resp, err := client.ListCertificates(context.TODO(), &acm.ListCertificatesInput{
			MaxItems:  aws.Int32(1000), // get as many as we can
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.CertificateSummaryList, resp.NextToken, nil
	})
}

func (ACMCertificateSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, certificate types.CertificateSummary) {
	arn := aws.ToString(certificate.CertificateArn)
	id := certificateID(arn)
	title := aws.ToString(certificate.DomainName)
	if title == "" {
		title = id
	}
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/acm/home#/certificates/" + id,
		ServiceID:   "acm",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		strings.ToLower(string(certificate.Status)),
		strings.ToLower(string(certificate.Type)),
		expiryDetail(certificate.NotAfter),
	))
}

// certificateID is what the console addresses a certificate by: the part of
// its ARN after the last slash, rather than after the last colon.
func certificateID(arn string) string {
	return arn[strings.LastIndex(arn, "/")+1:]
}

func expiryDetail(notAfter *time.Time) string {
	if notAfter == nil {
		return ""
	}
	return "expires " + notAfter.Format("2006-01-02")
}
