package awsworkflow

// APIGatewayAPI flattens a REST API and an HTTP/WebSocket API into one shape.
// The two are served by different AWS APIs but share a single entry in both
// the console and this workflow's service list.
type APIGatewayAPI struct {
	Id           string
	Name         string
	Description  string
	ProtocolType string
}
