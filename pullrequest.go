package bitbucket_v1

import "encoding/json"

type PullRequest struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
}

type Ref struct {
	Id           string `json:"id"`
	DisplayId    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
}

type User struct {
	Id           uint64 `json:"id"`
	Name         string `json:"name"`
	emailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

func (client *Client) GetPullrequests(project, repo string, ammount int) ([]*PullRequest, error) {
	var result = []*PullRequest{}
	url := client.requestUrl("/projects/%s/repos/%s/pull-requests?state=ALL", project, repo)
	values, err := client.getPaged(url, ammount, "")
	if err != nil {
		return nil, err
	}
	for _, value := range values {
		var v = &PullRequest{}
		err = json.Unmarshal(value, v)
		if err != nil {
			return result, err
		}
		result = append(result, v)
	}
	return result, nil
}
