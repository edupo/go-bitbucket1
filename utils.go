package bitbucket_v1

import "fmt"

func buildUrl(template string, args ...interface{}) string {
	if len(args) == 1 && args[0] == "" {
		return template
	}
	return fmt.Sprintf(template, args...)
}
