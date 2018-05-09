package bitbucket_v1

import (
	"encoding/json"
	"net/url"
)

type Repository struct {
	Slug          string   `json:"slug"`
	Id            uint64   `json:"id"`
	Name          string   `json:"name"`
	ScmId         string   `json:"scmId"`
	State         string   `json:"state"`
	StatusMessage string   `json:"statusMessage"`
	Forkable      bool     `json:"forkable"`
	Project       *Project `json:"project"`
	Public        bool     `json:"public"`
	Links         Links    `json:"links"`
}

type GetRepositoriesOpts struct {
	ProjectKey string
	Ammount    int
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

func (client *Client) GetRepositories(opts *GetRepositoriesOpts) ([]*Project, error) {
	result := []*Project{}
	urlString, err := getRepositoriesUrl(opts)
	if err != nil {
		return result, err
	}
	values, err := client.getPaged(urlString, opts.Ammount, "")
	if err != nil {
		return result, err
	}
	// Unmarshal the values
	for _, value := range values {
		var v = &Project{}
		err = json.Unmarshal(value, v)
		if err != nil {
			return result, err
		}
		result = append(result, v)
	}
	return result, nil
}

func (client *Client) GetRepository(projectKey, repositorySlug string) (*Repository, error) {
	var repository Repository
	urlString := buildUrl("/projects/%s/repos/%s", projectKey, repositorySlug)
	respBytes, err := client.execute("GET", urlString, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBytes, &repository)
	if err != nil {
		return nil, err
	}
	return &repository, nil
}
