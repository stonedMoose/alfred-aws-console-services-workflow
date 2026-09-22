package searchers

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

const (
	aliasPrefix           = "alias/"
	awsManagedAliasPrefix = "alias/aws/"
)

type KMSKeySearcher struct{}

func (s KMSKeySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "kms_keys", s.fetch, s.addToWorkflow)
}

// fetch lists aliases rather than keys, because a bare key id is not something
// anyone searches for. AWS managed keys are dropped so that this matches the
// console's customer managed keys page.
func (KMSKeySearcher) fetch(cfg aws.Config) ([]types.AliasListEntry, error) {
	client := kms.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.AliasListEntry, *string, error) {
		resp, err := client.ListAliases(context.TODO(), &kms.ListAliasesInput{
			Limit:  aws.Int32(100), // max allowed by this API
			Marker: awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return customerManagedAliases(resp.Aliases), awspaging.NextPageTokenIf(resp.Truncated, resp.NextMarker), nil
	})
}

func customerManagedAliases(aliases []types.AliasListEntry) []types.AliasListEntry {
	var customerManaged []types.AliasListEntry
	for _, alias := range aliases {
		if isCustomerManagedAlias(alias) {
			customerManaged = append(customerManaged, alias)
		}
	}
	return customerManaged
}

func isCustomerManagedAlias(alias types.AliasListEntry) bool {
	return alias.AliasName != nil &&
		!strings.HasPrefix(*alias.AliasName, awsManagedAliasPrefix) &&
		alias.TargetKeyId != nil
}

func (KMSKeySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, alias types.AliasListEntry) {
	keyID := aws.ToString(alias.TargetKeyId)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       strings.TrimPrefix(aws.ToString(alias.AliasName), aliasPrefix),
		ConsolePath: "/kms/home#/kms/keys/" + keyID,
		ServiceID:   "kms",
		ID:          aws.ToString(alias.AliasArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(keyID)
}
