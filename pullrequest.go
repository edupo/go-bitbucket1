package bitbucket_v1

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type PullRequest struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	FromRef     Ref    `json:"fromRef"`
	ToRef       Ref    `json:"toRef"`
}

type GetPullRequestsOpts struct {
	Project, Repo  string
	Ammount        int
	Direction      string
	At             string
	State          string
	Order          string
	WithAttributes bool
	WithProperties bool
}

func getPullRequestUrl(opts *GetPullRequestsOpts) (string, error) {
	urlString := buildUrl("/projects/%s/repos/%s/pull-requests", opts.Project, opts.Repo)
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}
	query := url.Values{}
	query.Add("direction", opts.Direction)
	query.Add("at", opts.At)
	query.Add("state", opts.State)
	query.Add("order", opts.Order)
	query.Add("withAttributes", strconv.FormatBool(opts.WithAttributes))
	query.Add("withProperties", strconv.FormatBool(opts.WithProperties))
	urlObject.RawQuery = query.Encode()
	return urlObject.String(), nil
}

func (client *Client) GetPullRequests(opts *GetPullRequestsOpts) ([]*PullRequest, error) {
	result := []*PullRequest{}
	urlString, err := getPullRequestUrl(opts)
	if err != nil {
		return result, err
	}
	values, err := client.getPaged(urlString, opts.Ammount, "")
	if err != nil {
		return result, err
	}
	// Unmarshal the values
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
