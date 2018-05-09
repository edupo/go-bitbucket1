package bitbucket_v1

import (
	"net/url"
	"strconv"
)

type GetProjectsOpts struct {
	Name, Permission string
	Ammount          int
}

type GetRepositoriesOpts struct {
	ProjectKey string
	Ammount    int
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

func getProjectsUrl(opts *GetProjectsOpts) (string, error) {
	urlString := buildUrl("/projects")
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}
	query := url.Values{}
	query.Add("name", opts.Name)
	query.Add("permission", opts.Permission)
	urlObject.RawQuery = query.Encode()
	return urlObject.String(), nil
}

func getRepositoriesUrl(opts *GetRepositoriesOpts) (string, error) {
	urlString := buildUrl("/projects")
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}
	query := url.Values{}
	query.Add("projectKey", opts.ProjectKey)
	urlObject.RawQuery = query.Encode()
	return urlObject.String(), nil
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
