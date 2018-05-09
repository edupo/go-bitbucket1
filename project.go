package bitbucket_v1

import (
	"encoding/json"
	"net/url"
)

type Project struct {
	Key         string `json:"key"`
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       Links  `json:"links"`
}

type GetProjectsOpts struct {
	Name, Permission string
	Ammount          int
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

func (client *Client) GetProjects(opts *GetProjectsOpts) ([]*Project, error) {
	result := []*Project{}
	urlString, err := getProjectsUrl(opts)
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

func (client *Client) GetProject(projectKey string) (*Project, error) {
	var project Project
	urlString := buildUrl("/projects/%s", projectKey)
	respBytes, err := client.execute("GET", urlString, "")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBytes, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
