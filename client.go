package bitbucket_v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultPageLenght = 10
	decRadix          = 10
)

const (
	defaultApiBaseURL = "https://bitbucket.org/api/1.0"
)

type Auth struct {
	user, password string
}

func NewSimpleAuth(user, password string) Auth {
	return Auth{
		user,
		password,
	}
}

type Client struct {
	BaseUrl    *url.URL
	Auth       *Auth
	PageLenght uint64
}

func NewClient(urlString string, auth *Auth) (*Client, error) {
	var client = Client{
		Auth:       auth,
		PageLenght: defaultPageLenght,
	}
	var err error

	if urlString == "" {
		urlString = defaultApiBaseURL
	}
	client.BaseUrl, err = url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	client.Auth = auth

	return &client, nil
}

func (client *Client) execute(method, urlString, text string) ([]byte, error) {
	body := strings.NewReader(text)
	req, err := http.NewRequest(method, client.BaseUrl.String()+urlString, body)

	if err != nil {
		return nil, err
	}
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	req.SetBasicAuth(client.Auth.user, client.Auth.password)

	httpClient := new(http.Client)
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) {
		return nil, fmt.Errorf(resp.Status)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	return ioutil.ReadAll(resp.Body)
}

type PagedResponse struct {
	Size          int               `json:"size"`
	Limit         int               `json:"limit"`
	IsLastPage    bool              `json:"isLastPage"`
	Values        []json.RawMessage `json:"values"`
	Start         int               `json:"start"`
	NextPageStart int               `json:"nextPageStart"`
}

func (client *Client) getPaged(urlString string, ammount int, text string) ([]json.RawMessage, error) {

	remaining := ammount
	var resp = &PagedResponse{}
	var limit uint64 = 0
	var values []json.RawMessage

	for remaining > 0 {

		// Check remaining values to get
		if uint64(remaining) > client.PageLenght || ammount < 0 {
			limit = client.PageLenght
		} else {
			limit = uint64(remaining)
		}

		// Setting page limit
		urlObject, err := url.Parse(urlString)
		if err != nil {
			return nil, err
		}
		q := urlObject.Query()
		q.Set("limit", strconv.FormatUint(limit, decRadix))
		q.Set("start", strconv.FormatInt(int64(resp.NextPageStart), decRadix))
		urlObject.RawQuery = q.Encode()
		urlString = urlObject.String()

		// Perform the GET
		stream, err := client.execute("GET", urlString, text)
		if err != nil {
			return values, err
		}

		var resp = &PagedResponse{}
		err = json.Unmarshal(stream, resp)
		if err != nil {
			return values, err
		}

		values = append(values, resp.Values...)
		remaining -= resp.Size
		if resp.IsLastPage || remaining == 0 {
			break
		}
	}

	return values, nil
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
