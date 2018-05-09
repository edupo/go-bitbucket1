package bitbucket_v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	DEFAULT_PAGE_LENGHT = 10
	DEC_RADIX           = 10
)

var (
	apiBaseURL           = "https://bitbucket.org/api/1.0"
	ErrNotPagedResponse  = errors.New("Is not a paged response")
	ErrNoFieldInResponse = errors.New("The field does not exist in the response")
	ErrBadResponse       = errors.New("Response is malformed")
)

func GetApiBaseURL() string {
	return apiBaseURL
}

func SetApiBaseURL(urlStr string) {
	apiBaseURL = urlStr
}

type Client struct {
	Auth       *auth
	PageLenght uint64
}

type auth struct {
	user, password string
}

func NewClient(user, password string) *Client {
	return &Client{
		PageLenght: DEFAULT_PAGE_LENGHT,
		Auth: &auth{
			user:     user,
			password: password,
		}}
}

func (client *Client) execute(method, urlStr, text string) (interface{}, error) {

	var result interface{}
	resBodyBytes, err := client.executeString(method, urlStr, text)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resBodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (client *Client) executeString(method, urlStr, text string) ([]byte, error) {
	body := strings.NewReader(text)
	req, err := http.NewRequest(method, urlStr, body)

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
	IsLastPage    bool              `json:"isLatPage"`
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
		q.Set("limit", strconv.FormatUint(limit, DEC_RADIX))
		q.Set("start", strconv.FormatInt(int64(resp.NextPageStart), DEC_RADIX))
		urlObject.RawQuery = q.Encode()
		urlString = urlObject.String()

		// Perform the GET
		stream, err := client.executeString("GET", urlString, text)
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

func (c *Client) requestUrl(template string, args ...interface{}) string {

	if len(args) == 1 && args[0] == "" {
		return GetApiBaseURL() + template
	}
	return GetApiBaseURL() + fmt.Sprintf(template, args...)
}
