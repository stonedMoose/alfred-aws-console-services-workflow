package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// apiGatewayAPI unifies the REST APIs of API Gateway v1 and the HTTP and
// WebSocket APIs of v2: they come from two different AWS APIs, with different
// types, but sit side by side in the console.
type apiGatewayAPI struct {
	Id           string
	Name         string
	Description  string
	ProtocolType string
}

const restProtocolType = "REST"

type APIGatewayAPISearcher struct{}

func (s APIGatewayAPISearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "apigateway_apis", s.fetch, s.addToWorkflow)
}

func (APIGatewayAPISearcher) fetch(cfg aws.Config) ([]apiGatewayAPI, error) {
	restAPIs, err := fetchRestAPIs(cfg)
	if err != nil {
		return nil, err
	}
	v2APIs, err := fetchV2APIs(cfg)
	if err != nil {
		return nil, err
	}
	return append(restAPIs, v2APIs...), nil
}

func fetchRestAPIs(cfg aws.Config) ([]apiGatewayAPI, error) {
	client := apigateway.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(position string) ([]apiGatewayAPI, *string, error) {
		resp, err := client.GetRestApis(context.TODO(), &apigateway.GetRestApisInput{
			Limit:    aws.Int32(500), // max allowed by this API
			Position: awspaging.TokenOrNil(position),
		})
		if err != nil {
			return nil, nil, err
		}
		apis := make([]apiGatewayAPI, 0, len(resp.Items))
		for _, item := range resp.Items {
			apis = append(apis, apiGatewayAPI{
				Id:           aws.ToString(item.Id),
				Name:         aws.ToString(item.Name),
				Description:  aws.ToString(item.Description),
				ProtocolType: restProtocolType,
			})
		}
		return apis, resp.Position, nil
	})
}

func fetchV2APIs(cfg aws.Config) ([]apiGatewayAPI, error) {
	client := apigatewayv2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]apiGatewayAPI, *string, error) {
		resp, err := client.GetApis(context.TODO(), &apigatewayv2.GetApisInput{
			MaxResults: aws.String("500"), // this API takes its page size as a string
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		apis := make([]apiGatewayAPI, 0, len(resp.Items))
		for _, item := range resp.Items {
			apis = append(apis, apiGatewayAPI{
				Id:           aws.ToString(item.ApiId),
				Name:         aws.ToString(item.Name),
				Description:  aws.ToString(item.Description),
				ProtocolType: string(item.ProtocolType),
			})
		}
		return apis, resp.NextToken, nil
	})
}

func (APIGatewayAPISearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, api apiGatewayAPI) {
	title := api.Name
	if title == "" {
		title = api.Id
	}
	region := searchArgs.GetRegion()
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title: title,
		// the region goes in the path because the console reads it from the
		// query string alongside the api id
		ConsolePath: fmt.Sprintf("/apigateway/main/apis/%s/%s?api=%s&region=%s", api.Id, api.consoleTab(), api.Id, region),
		ServiceID:   "apigateway",
	}).Subtitle(subtitleFrom(api.Id, api.ProtocolType, api.Description))
}

// consoleTab is the tab the console opens an API on: REST APIs are made of
// resources, the others of routes.
func (api apiGatewayAPI) consoleTab() string {
	if api.ProtocolType == restProtocolType {
		return "resources"
	}
	return "routes"
}
