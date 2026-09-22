// Package awspaging walks the token based pagination that the AWS list APIs
// share: every call hands back one page of items and, while more remain, a
// token for the next call.
package awspaging

import "github.com/aws/aws-sdk-go-v2/aws"

// PageFetcher performs one call of a paginated list API. It receives the token
// of the page to fetch, empty for the first page, and returns the items of that
// page together with the token of the following page. A nil or empty token
// means the last page has been read.
type PageFetcher[Item any] func(pageToken string) (items []Item, nextPageToken *string, err error)

// FetchAllPages calls fetchPage until the API stops handing back a token and
// returns every item in the order the pages arrived.
func FetchAllPages[Item any](fetchPage PageFetcher[Item]) ([]Item, error) {
	var items []Item
	pageToken := ""
	for {
		pageItems, nextPageToken, err := fetchPage(pageToken)
		if err != nil {
			return nil, err
		}
		items = append(items, pageItems...)
		if aws.ToString(nextPageToken) == "" {
			return items, nil
		}
		pageToken = *nextPageToken
	}
}

// TokenOrNil turns the empty first-page token into the nil pointer by which
// the SDK inputs mean "start from the first page".
func TokenOrNil(pageToken string) *string {
	if pageToken == "" {
		return nil
	}
	return aws.String(pageToken)
}

// NextPageTokenIf keeps the token only while the API reports more pages, for
// the APIs that flag truncation separately from the token.
func NextPageTokenIf(hasMorePages bool, nextPageToken *string) *string {
	if !hasMorePages {
		return nil
	}
	return nextPageToken
}
