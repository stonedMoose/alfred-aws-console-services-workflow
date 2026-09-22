package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type APIGatewayAPISearcher struct{}

func (s APIGatewayAPISearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// REST APIs and HTTP/WebSocket APIs come from two different AWS APIs but sit
// side by side in the console, so both are fetched here
func (s APIGatewayAPISearcher) fetch(cfg aws.Config) ([]awsworkflow.APIGatewayAPI, error) {
	entities, err := fetchRestAPIs(cfg)
	if err != nil {
		return nil, err
	}

	v2Entities, err := fetchV2APIs(cfg)
	if err != nil {
		return nil, err
	}

	return append(entities, v2Entities...), nil
}

func fetchRestAPIs(cfg aws.Config) ([]awsworkflow.APIGatewayAPI, error) {
	svc := apigateway.NewFromConfig(cfg)

	var entities []awsworkflow.APIGatewayAPI
	position := ""
	for {
		params := &apigateway.GetRestApisInput{
			Limit: aws.Int32(500), // max allowed by this API
		}
		if position != "" {
			params.Position = aws.String(position)
		}
		resp, err := svc.GetRestApis(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			entity := awsworkflow.APIGatewayAPI{ProtocolType: "REST"}
			if item.Id != nil {
				entity.Id = *item.Id
			}
			if item.Name != nil {
				entity.Name = *item.Name
			}
			if item.Description != nil {
				entity.Description = *item.Description
			}
			entities = append(entities, entity)
		}

		if resp.Position != nil && *resp.Position != "" {
			position = *resp.Position
		} else {
			break
		}
	}

	return entities, nil
}

func fetchV2APIs(cfg aws.Config) ([]awsworkflow.APIGatewayAPI, error) {
	svc := apigatewayv2.NewFromConfig(cfg)

	var entities []awsworkflow.APIGatewayAPI
	nextToken := ""
	for {
		params := &apigatewayv2.GetApisInput{
			MaxResults: aws.String("500"), // this API takes its page size as a string
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.GetApis(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			entity := awsworkflow.APIGatewayAPI{ProtocolType: string(item.ProtocolType)}
			if item.ApiId != nil {
				entity.Id = *item.ApiId
			}
			if item.Name != nil {
				entity.Name = *item.Name
			}
			if item.Description != nil {
				entity.Description = *item.Description
			}
			entities = append(entities, entity)
		}

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s APIGatewayAPISearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity awsworkflow.APIGatewayAPI) {
	title := entity.Name
	if title == "" {
		title = entity.Id
	}

	subtitleArray := []string{entity.Id}
	if entity.ProtocolType != "" {
		subtitleArray = append(subtitleArray, entity.ProtocolType)
	}
	if entity.Description != "" {
		subtitleArray = append(subtitleArray, entity.Description)
	}
	subtitle := strings.Join(subtitleArray, " – ")

	// a REST API opens on its resources, everything else on its routes
	tab := "routes"
	if entity.ProtocolType == "REST" {
		tab = "resources"
	}
	// the region goes in the path because the console reads it from the query
	// string alongside the api id
	path := fmt.Sprintf("/apigateway/main/apis/%s/%s?api=%s&region=%s", entity.Id, tab, entity.Id, searchArgs.GetRegion())
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("apigateway")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
