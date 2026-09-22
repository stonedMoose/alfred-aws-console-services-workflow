package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

var HeadersToIgnore = []string{
	"Authorization",
	"X-Amz-Date",
	"X-Amz-Content-Sha256",
	"User-Agent",
	"Amz-Sdk-Request",
	"Amz-Sdk-Invocation-Id",
	"X-Amzn-Requestid",
	"Date",
	"X-Amz-Id-2",
	"X-Amz-Request-Id",
	"Content-Length",
}

func deepEqualContents(x, y any) bool {
	if reflect.ValueOf(x).IsNil() {
		if reflect.ValueOf(y).IsNil() {
			return true
		} else {
			return reflect.ValueOf(y).Len() == 0
		}
	} else {
		if reflect.ValueOf(y).IsNil() {
			return reflect.ValueOf(x).Len() == 0
		} else {
			return reflect.DeepEqual(x, y)
		}
	}
}

// TODO fix rest of matchers

func bodyMatches(r *http.Request, i cassette.Request) bool {
	if r.Body != nil {
		var buffer bytes.Buffer
		if _, err := buffer.ReadFrom(r.Body); err != nil {
			return false
		}

		r.Body = io.NopCloser(bytes.NewBuffer(buffer.Bytes()))
		if buffer.String() != i.Body {
			return false
		}
	} else {
		if len(i.Body) != 0 {
			return false
		}
	}

	return true
}

// modified version from default matcher
var CustomMatcher = func(r *http.Request, i cassette.Request) bool {
	if r.Method != i.Method {
		return false
	}

	if r.URL.String() != i.URL {
		return false
	}

	if r.Proto != i.Proto {
		return false
	}

	if r.ProtoMajor != i.ProtoMajor {
		return false
	}

	if r.ProtoMinor != i.ProtoMinor {
		return false
	}

	requestHeader := r.Header.Clone()
	cassetteRequestHeaders := i.Headers.Clone()

	for _, header := range HeadersToIgnore {
		delete(requestHeader, header)
		delete(cassetteRequestHeaders, header)
	}

	if !deepEqualContents(requestHeader, cassetteRequestHeaders) {
		return false
	}

	if !bodyMatches(r, i) {
		return false
	}

	if !deepEqualContents(r.TransferEncoding, i.TransferEncoding) {
		return false
	}

	if r.Host != i.Host {
		return false
	}

	// Only ParseForm for non-GET requests since that would use query params
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		err := r.ParseForm()
		if err != nil {
			return false
		}
	}
	if !deepEqualContents(r.Form, i.Form) {
		return false
	}

	if !deepEqualContents(r.Trailer, i.Trailer) {
		return false
	}

	if r.RemoteAddr != i.RemoteAddr {
		return false
	}

	if r.RequestURI != i.RequestURI {
		return false
	}

	return true
}

func prettyPrint(raw string) string {
	rawBytes := []byte(raw)

	if json.Valid(rawBytes) {
		var v interface{}
		if err := json.Unmarshal(rawBytes, &v); err == nil {
			if pretty, err := json.MarshalIndent(v, "", "  "); err == nil {
				return string(pretty)
			}
		}
	}

	return ""
}

func prettyFormatInteraction(i *cassette.Interaction) error {
	// comment out because this messages with matcher
	// if i.Request.Body != "" {
	// 	// Read and reformat
	// 	reqRaw := []byte(i.Request.Body)
	// 	i.Request.Body = prettyPrint(reqRaw)
	// }

	if i.Response.Body != "" {
		prettyResponseBody := prettyPrint(i.Response.Body)
		if prettyResponseBody != "" {
			i.Response.Body = prettyResponseBody
		}
	}

	return nil
}

var environmentIdRegex *regexp.Regexp = regexp.MustCompile(`e-[a-zA-Z0-9]{8,}`)
var instanceIdRegex *regexp.Regexp = regexp.MustCompile(`i-[a-zA-Z0-9]{8,}`)
var dbIdRegex *regexp.Regexp = regexp.MustCompile(`db-[a-zA-Z0-9]{8,}`)
var amiIdRegex *regexp.Regexp = regexp.MustCompile(`ami-[a-zA-Z0-9]{8,}`)
var vpcIdRegex *regexp.Regexp = regexp.MustCompile(`vpc-[a-zA-Z0-9]{8,}`)
var subnetIdRegex *regexp.Regexp = regexp.MustCompile(`subnet-[a-zA-Z0-9]{8,}`)
var namespaceRegex *regexp.Regexp = regexp.MustCompile(`ns-[a-zA-Z0-9]{8,}`)
var securityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`sg-[a-zA-Z0-9]{8,}`)
var expandedSecurityGroupIdRegex *regexp.Regexp = regexp.MustCompile(`securitygroup-[a-zA-Z0-9]{8,}`)
var volumeIdRegex *regexp.Regexp = regexp.MustCompile(`vol-[a-zA-Z0-9]{8,}`)
var attachmentIdRegex *regexp.Regexp = regexp.MustCompile(`eni-attach-[a-zA-Z0-9]{8,}`)
var reservationIdRegex *regexp.Regexp = regexp.MustCompile(`r-[a-zA-Z0-9]{8,}`)
var snapshotIdRegex *regexp.Regexp = regexp.MustCompile(`snap-[a-zA-Z0-9]{8,}`)
var fileSystemIdRegex *regexp.Regexp = regexp.MustCompile(`fs-[a-zA-Z0-9]{8,}`)
var keyPairIdRegex *regexp.Regexp = regexp.MustCompile(`key-[a-zA-Z0-9]{8,}`)

var accountIdInArn *regexp.Regexp = regexp.MustCompile(`:[0-9]{10,}:`)
var longNumberInXmlTag *regexp.Regexp = regexp.MustCompile(`>[0-9]{8,}<`) // we're going to assume that any numeric xml values are identifications of some sort, so just sanitize it
var uuidv2Regex *regexp.Regexp = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
var iso8601Regex *regexp.Regexp = regexp.MustCompile(`\\d{4}-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\d(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?`)

// Updated idTagRegex for better content capture, case-insensitive, dot-all
var idTagRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(id|DbiResourceId|HostedZoneId)>(.*?)</(?:id|DbiResourceId|HostedZoneId)>`)
var keyNameTagRegex *regexp.Regexp = regexp.MustCompile(`(?i)<(keyName)>.+</(keyName)>`)
var masterUsernameTagRegex *regexp.Regexp = regexp.MustCompile(`(?i)<(MasterUsername)>.+</(MasterUsername)>`)
var nextTokenTagRegex *regexp.Regexp = regexp.MustCompile(`(?i)<(NextToken)>.+</(NextToken)>`)

// Regex to detect path-like IDs and capture the prefix
var idPathPrefixRegex = regexp.MustCompile(`^(.*\/)[^/]+$`)

const sanitizedIDInPath = "000000"
const sanitizedIDDefault = "000000000000"

var ipv4Regex *regexp.Regexp = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
var macAddressRegex *regexp.Regexp = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)

var beanstalkSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkLoadBalancerSecurityGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBLoadBalancerSecurityGroup-[0-9A-Z]{10,}`)
var beanstalkAutoScalingGroupNameRegex *regexp.Regexp = regexp.MustCompile(`AWSEBAutoScalingGroup-[0-9A-Z]{10,}`)

var amazonawsUrlRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+\.amazonaws\.com`)
var beanstalkUrlSubdomainRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.elasticbeanstalk\.com`)
var internalUrlRegex *regexp.Regexp = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+\.internal`)

// Regex to find <Versions> blocks (case-insensitive, dot-all)
var versionsBlockRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(Versions)>(.*?)</Versions>`)

// Regex to find <member> tags within a block (case-insensitive, dot-all)
var memberTagRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(member)>(.*?)</member>`)

func sanitizeBody(body string) string {
	body = uuidv2Regex.ReplaceAllString(body, "00000000-0000-0000-0000-000000000000")
	body = environmentIdRegex.ReplaceAllString(body, "e-aaaaaaaaaa")
	body = instanceIdRegex.ReplaceAllString(body, "i-aaaaaaaaaa")
	body = dbIdRegex.ReplaceAllString(body, "db-AAAAAAAAAA")
	body = amiIdRegex.ReplaceAllString(body, "ami-aaaaaaaaaa")
	body = vpcIdRegex.ReplaceAllString(body, "vpc-aaaaaaaaaa")
	body = subnetIdRegex.ReplaceAllString(body, "subnet-aaaaaaaaaa")
	body = namespaceRegex.ReplaceAllString(body, "ns-aaaaaaaaaa")
	body = securityGroupIdRegex.ReplaceAllString(body, "sg-aaaaaaaaaa")
	body = expandedSecurityGroupIdRegex.ReplaceAllString(body, "securitygroup-aaaaaaaaaa")
	body = volumeIdRegex.ReplaceAllString(body, "vol-aaaaaaaaaa")
	body = attachmentIdRegex.ReplaceAllString(body, "eni-attach-aaaaaaaaaa")
	body = reservationIdRegex.ReplaceAllString(body, "r-aaaaaaaaaa")
	body = snapshotIdRegex.ReplaceAllString(body, "snap-aaaaaaaaaa")
	body = fileSystemIdRegex.ReplaceAllString(body, "fs-aaaaaaaaaa")
	body = keyPairIdRegex.ReplaceAllString(body, "key-aaaaaaaaaa")

	body = accountIdInArn.ReplaceAllString(body, ":0000000000:")
	body = longNumberInXmlTag.ReplaceAllString(body, ">00000000<")
	body = iso8601Regex.ReplaceAllString(body, "2020-01-01T00:00:00.000Z")

	// Conditional sanitization for ID-like tags
	body = idTagRegex.ReplaceAllStringFunc(body, func(match string) string {
		submatches := idTagRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match // Should not happen, but safeguard
		}
		openingTagName := submatches[1] // e.g., "Id", "HostedZoneId"
		originalContent := submatches[2]

		var sanitizedValue string
		pathSubmatches := idPathPrefixRegex.FindStringSubmatch(originalContent)
		if pathSubmatches != nil {
			// Content is path-like, preserve prefix
			pathPrefix := pathSubmatches[1]
			sanitizedValue = pathPrefix + sanitizedIDInPath
		} else {
			// Content is not path-like, use default sanitization
			sanitizedValue = sanitizedIDDefault
		}
		return fmt.Sprintf("<%s>%s</%s>", openingTagName, sanitizedValue, openingTagName)
	})

	body = masterUsernameTagRegex.ReplaceAllString(body, "<$1>aaaaaaaaaaaa</$2>")
	body = keyNameTagRegex.ReplaceAllString(body, "<$1>aaaaaaaaaa</$2>")
	body = nextTokenTagRegex.ReplaceAllString(body, "<$1>BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB</$2>")

	body = ipv4Regex.ReplaceAllString(body, "0.0.0.0")
	body = macAddressRegex.ReplaceAllString(body, "00:00:00:00:00:00")

	body = beanstalkSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkLoadBalancerSecurityGroupNameRegex.ReplaceAllString(body, "AWSEBLoadBalancerSecurityGroup-AAAAAAAAAAAA")
	body = beanstalkAutoScalingGroupNameRegex.ReplaceAllString(body, "AWSEBAutoScalingGroup-AAAAAAAAAAAA")

	body = amazonawsUrlRegex.ReplaceAllString(body, "subdomain.us-west-2.service.amazonaws.com")
	body = beanstalkUrlSubdomainRegex.ReplaceAllString(body, "subdomain.us-west-2.elasticbeanstalk.com")
	body = internalUrlRegex.ReplaceAllString(body, "subdomain.us-west-2.service.internal")

	// --- New Versions/Member Sanitization ---
	body = versionsBlockRegex.ReplaceAllStringFunc(body, func(versionsBlock string) string {
		// Extract the content between <Versions> and </Versions>
		matches := versionsBlockRegex.FindStringSubmatch(versionsBlock)
		if len(matches) < 3 {
			return versionsBlock // Should not happen with the regex, but safe guard
		}
		openingTag := matches[1] // e.g., "Versions" or "versions"
		content := matches[2]

		memberCounter := 0
		sanitizedContent := memberTagRegex.ReplaceAllStringFunc(content, func(memberBlock string) string {
			memberMatches := memberTagRegex.FindStringSubmatch(memberBlock)
			if len(memberMatches) < 3 {
				return memberBlock // Safeguard
			}
			memberOpeningTag := memberMatches[1] // e.g., "member" or "Member"

			// Generate the incremental value (repeat digit 8 times)
			digit := memberCounter % 10
			sanitizedValue := strings.Repeat(fmt.Sprintf("%d", digit), 8)
			memberCounter++

			// Reconstruct the member tag preserving case
			return fmt.Sprintf("<%s>%s</%s>", memberOpeningTag, sanitizedValue, memberOpeningTag)
		})

		// Reconstruct the Versions block preserving case
		return fmt.Sprintf("<%s>%s</%s>", openingTag, sanitizedContent, openingTag)
	})

	return body
}

// Interactions get recorded against whatever region the person doing the
// recording has configured, but they are always replayed with the region in
// tests/test_aws_config_file. Rewriting the region as the cassette is saved
// lets a fixture be recorded from any region and still match on replay.
const replayRegion = "us-west-2"

var awsRegionRegex *regexp.Regexp = regexp.MustCompile(`\b(af|ap|ca|eu|il|me|mx|sa|us)-(gov-)?(north|south|east|west|central|northeast|northwest|southeast|southwest)-[0-9]\b`)

func normalizeRegion(s string) string {
	return awsRegionRegex.ReplaceAllString(s, replayRegion)
}

// The JSON protocol services carry resource names in plain fields that the
// XML-oriented sanitizers above never see. Rather than enumerate every field
// of every service, any field whose name ends in one of these is treated as
// identifying.
var sensitiveFieldSuffixes = []string{
	"arn",
	"arns",
	"description",
	"name",
	"names",
	"token",
	"uri",
	"uris",
	"url",
	"urls",
}

// fields that carry identifying values but are not named after them
var extraSensitiveFields = map[string]bool{
	"apiendpoint":             true,
	"apiid":                   true,
	"bucket":                  true,
	"clusters":                true,
	"comment":                 true,
	"exclusivestarttablename": true,
	"families":                true,
	"id":                      true,
	"keyfingerprint":          true,
	"originpath":              true,
	"prefix":                  true,
	"keyid":                   true,
	"registryid":              true,
	"targetkeyid":             true,
	"value":                   true,
}

func isSensitiveField(key string) bool {
	lower := strings.ToLower(key)
	if extraSensitiveFields[lower] {
		return true
	}
	for _, suffix := range sensitiveFieldSuffixes {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

var resourceNameWordRegex *regexp.Regexp = regexp.MustCompile(`[A-Za-z0-9]+`)

// sanitizeResourceName replaces the identifying words of a resource name with
// filler while keeping its shape — slashes, dashes, arn and url prefixes — so
// that fixtures still exercise the url escaping each searcher does.
func sanitizeResourceName(name string) string {
	prefix := ""
	rest := name
	if strings.HasPrefix(name, "arn:") {
		if parts := strings.SplitN(name, ":", 6); len(parts) == 6 {
			prefix = strings.Join(parts[:5], ":") + ":"
			rest = parts[5]
		}
	} else if strings.HasPrefix(name, "https://") {
		if i := strings.Index(name[len("https://"):], "/"); i >= 0 {
			prefix = name[:len("https://")+i+1]
			rest = name[len("https://")+i+1:]
		}
	}

	return prefix + sanitizeWordsOutsideEntities(rest)
}

// XML character entities have to survive intact — turning `&quot;` into
// `&aaaa;` leaves a body that no longer parses
var xmlEntityRegex *regexp.Regexp = regexp.MustCompile(`&(?:[a-zA-Z]+|#[0-9]+);`)

func sanitizeWordsOutsideEntities(s string) string {
	entities := xmlEntityRegex.FindAllString(s, -1)
	parts := xmlEntityRegex.Split(s, -1)

	var builder strings.Builder
	for i, part := range parts {
		builder.WriteString(resourceNameWordRegex.ReplaceAllStringFunc(part, func(word string) string {
			return strings.Repeat("a", len(word))
		}))
		if i < len(entities) {
			builder.WriteString(entities[i])
		}
	}

	return builder.String()
}

// array members inherit the key of the array they belong to, so that
// "TableNames": ["orders"] is sanitized the same way a "TableName" would be
func sanitizeJSONValue(key string, value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		for k, v := range typed {
			typed[k] = sanitizeJSONValue(k, v)
		}
		return typed
	case []interface{}:
		for i, v := range typed {
			typed[i] = sanitizeJSONValue(key, v)
		}
		return typed
	case string:
		if isSensitiveField(key) {
			return sanitizeResourceName(typed)
		}
		return typed
	}
	return value
}

// The XML protocol services carry the same values in element text. A leaf
// element's text cannot contain "<", so the closing tag needs no backreference
// (which RE2 could not express anyway).
var xmlLeafElementRegex *regexp.Regexp = regexp.MustCompile(`<([A-Za-z][A-Za-z0-9]*)>([^<]*)</[A-Za-z][A-Za-z0-9]*>`)

// a paginated request echoes the token it was given, so the token in a
// response body has to stay byte for byte what the next request carries
var xmlPaginationTags = map[string]bool{
	"marker":     true,
	"nextmarker": true,
	"nexttoken":  true,
}

// policy documents are url encoded json that can name external accounts and
// identity providers, and nothing in this workflow reads them
var xmlTagsToBlank = map[string]string{
	"assumerolepolicydocument": "%7B%7D",
	"policydocument":           "%7B%7D",
}

func sanitizeXMLTags(body string) string {
	return xmlLeafElementRegex.ReplaceAllStringFunc(body, func(match string) string {
		submatches := xmlLeafElementRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		tag, text := submatches[1], submatches[2]

		if xmlPaginationTags[strings.ToLower(tag)] {
			// these already get consistent treatment above and in the request
			// query params, and rewriting them here would break replay
			return match
		}
		if blanked, ok := xmlTagsToBlank[strings.ToLower(tag)]; ok {
			return fmt.Sprintf("<%s>%s</%s>", tag, blanked, tag)
		}
		if !isSensitiveField(tag) {
			return match
		}
		return fmt.Sprintf("<%s>%s</%s>", tag, sanitizeResourceName(text), tag)
	})
}

// EC2 hangs resource names off a tag set rather than a dedicated element, so
// the values inside one are sanitized while tags elsewhere are left alone
var ec2TagSetRegex *regexp.Regexp = regexp.MustCompile(`(?is)<(tagSet|tagSpecificationSet)>(.*?)</(?:tagSet|tagSpecificationSet)>`)
var ec2TagValueRegex *regexp.Regexp = regexp.MustCompile(`(?is)<value>(.*?)</value>`)

func sanitizeEC2TagSets(body string) string {
	return ec2TagSetRegex.ReplaceAllStringFunc(body, func(tagSet string) string {
		return ec2TagValueRegex.ReplaceAllStringFunc(tagSet, func(match string) string {
			submatches := ec2TagValueRegex.FindStringSubmatch(match)
			if len(submatches) < 2 {
				return match
			}
			return fmt.Sprintf("<value>%s</value>", sanitizeResourceName(submatches[1]))
		})
	})
}

func sanitizeJSONBody(body string) string {
	if !json.Valid([]byte(body)) {
		return body
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return body
	}

	sanitized, err := json.Marshal(sanitizeJSONValue("", parsed))
	if err != nil {
		return body
	}

	return string(sanitized)
}

func sanitizeAndFormatBodyHook(i *cassette.Interaction) error {
	if i.WasReplayed() {
		// don't reformat again if we are playing back
		return nil
	}

	for _, header := range HeadersToIgnore {
		delete(i.Request.Headers, header)
		delete(i.Response.Headers, header)
	}
	i.Request.ContentLength = 0

	i.Request.URL = normalizeRegion(i.Request.URL)
	i.Request.Host = normalizeRegion(i.Request.Host)
	i.Request.RequestURI = normalizeRegion(i.Request.RequestURI)

	i.Request.Body = sanitizeJSONBody(sanitizeEC2TagSets(sanitizeXMLTags(sanitizeBody(normalizeRegion(i.Request.Body)))))
	i.Response.Body = sanitizeJSONBody(sanitizeEC2TagSets(sanitizeXMLTags(sanitizeBody(normalizeRegion(i.Response.Body)))))

	if err := prettyFormatInteraction(i); err != nil {
		return err
	}

	if i.Request.Body != "" {
		parsedQuery, parseErr := url.ParseQuery(i.Request.Body)
		if parseErr == nil {
			changed := false
			for _, qpSanitizer := range requestQueryParamSanitizers {
				if value, ok := parsedQuery[qpSanitizer.ParamName]; ok {
					if len(value) > 0 && value[0] != "" {
						parsedQuery.Set(qpSanitizer.ParamName, qpSanitizer.SanitizedValue)
						changed = true
					}
				}
			}

			if changed {
				// Re-encode the modified query string and update Request.Body
				i.Request.Body = parsedQuery.Encode()

				// Update Request.Form as well for consistency
				if i.Request.Form == nil {
					i.Request.Form = make(url.Values)
				}
				for key, values := range parsedQuery {
					i.Request.Form[key] = values
				}
			}
		} else {
			fmt.Printf("Warning: Could not parse request body as query params: %v\n", parseErr)
		}
	}

	return nil
}
