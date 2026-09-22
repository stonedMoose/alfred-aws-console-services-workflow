package searchers

import (
	"fmt"
	"strings"
	"time"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// detailSeparator sets the details of a subtitle apart from each other.
const detailSeparator = " – "

// subtitleFrom joins the details known about a resource, skipping the empty
// ones, in the house style of the item subtitles.
func subtitleFrom(details ...string) string {
	return joinKnown(detailSeparator, details...)
}

// joinKnown joins the non-empty details with separator.
func joinKnown(separator string, details ...string) string {
	known := make([]string, 0, len(details))
	for _, detail := range details {
		if detail != "" {
			known = append(known, detail)
		}
	}
	return strings.Join(known, separator)
}

// dateDetail renders "<label> <date>", or nothing when the date is unknown.
func dateDetail(label string, date *time.Time) string {
	if date == nil {
		return ""
	}
	return label + " " + date.Format(time.UnixDate)
}

// formatIfSet renders value with format, or nothing when value is unknown.
func formatIfSet[Value any](format string, value *Value) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf(format, *value)
}

// namedOrID picks the title of a resource that may carry a name next to its
// id: the name when there is one, in which case the id becomes a detail of the
// subtitle, otherwise the id.
func namedOrID(name, id string) (title, idDetail string) {
	if name == "" {
		return id, ""
	}
	return name, id
}

// ec2NameTag returns the value of the "Name" tag, the closest thing an EC2
// resource has to a name.
func ec2NameTag(tags []ec2types.Tag) string {
	for _, tag := range tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}

// arnResourceName returns what follows the last colon of an ARN, which is the
// resource name for the services that leave it unqualified.
func arnResourceName(arn string) string {
	return arn[strings.LastIndex(arn, ":")+1:]
}
