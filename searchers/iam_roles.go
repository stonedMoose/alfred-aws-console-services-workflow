package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type IAMRoleSearcher struct{}

func (s IAMRoleSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s IAMRoleSearcher) fetch(cfg aws.Config) ([]types.Role, error) {
	svc := iam.NewFromConfig(cfg)

	var entities []types.Role
	marker := ""
	for {
		params := &iam.ListRolesInput{
			MaxItems: aws.Int32(1000), // get as many as we can
		}
		if marker != "" {
			params.Marker = aws.String(marker)
		}
		resp, err := svc.ListRoles(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Roles...)

		if resp.IsTruncated && resp.Marker != nil {
			marker = *resp.Marker
		} else {
			break
		}
	}

	return entities, nil
}

func (s IAMRoleSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Role) {
	title := *entity.RoleName

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	subtitleArray = util.AppendString(subtitleArray, entity.Path)
	if entity.CreateDate != nil {
		subtitleArray = append(subtitleArray, "Created "+entity.CreateDate.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/iamv2/home#/roles/details/%s", url.PathEscape(title))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("iam")).
		Valid(true)

	arn := ""
	if entity.Arn != nil {
		arn = *entity.Arn
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
