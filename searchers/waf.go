package searchers

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

// wafScope restricts the WAF searchers to regional resources.
// TODO support the CLOUDFRONT scope somehow
const wafScope = types.ScopeRegional

// wafConsolePath builds the path of a WAF resource page. WAF counts as a
// global service, yet its resources are regional, so the region has to be
// written into the path rather than derived from the service.
func wafConsolePath(format string, name, id, region string) string {
	return fmt.Sprintf(format, name, id, region)
}
