/*
Provide a basic communication client with confluence.

I chose not to use the client of package "github.com/virtomize/confluence-go-api",
because it does not use the buffers for body of the request and response.
And this client does not allow passing a context to limit/cancel requests.
*/
package confluence

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	goconfluence "github.com/virtomize/confluence-go-api"
)

type Client struct {
	host          string
	authorization string
	httpClient    *http.Client
}

type errorResponse struct {
	Message string `json:"message"`
}

func GenerateBasicAuthorization(email, apiToken string) string {
	return fmt.Sprintf(
		"Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, apiToken))),
	)
}

func NewClient(httpClient *http.Client, host, authorization string) *Client {
	return &Client{
		httpClient:    httpClient,
		host:          fmt.Sprintf("%s/wiki/rest/api", host),
		authorization: authorization,
	}
}

func (c *Client) CreateContent(ctx context.Context, content *goconfluence.Content) (*goconfluence.Content, error) {
	var result *goconfluence.Content
	if err := c.request(ctx, "CreateContent", http.MethodPost, "content", content, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) UpdateContent(ctx context.Context, content *goconfluence.Content) (*goconfluence.Content, error) {
	var result *goconfluence.Content
	if err := c.request(ctx, "UpdateContent", http.MethodPut, fmt.Sprintf("content/%s", content.ID), content, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) FindContent(ctx context.Context, content *goconfluence.Content, start int) (*goconfluence.ContentSearch, error) {
	req, err := c.createRequest(ctx, "FindContent", http.MethodGet, "content", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	q.Add("expand", "space,version,ancestors")
	q.Add("start", fmt.Sprint(start))

	q.Add("type", content.Type)
	q.Add("spaceKey", content.Space.Key)
	q.Add("title", content.Title)

	req.URL.RawQuery = q.Encode()

	var result *goconfluence.ContentSearch
	if _, err := c.doRequest("FindContent", req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) request(ctx context.Context, name, method, path string, body interface{}, result interface{}) error {
	req, err := c.createRequest(ctx, name, method, path, body)
	if err != nil {
		return err
	}

	defer req.Body.Close()

	res, err := c.doRequest(name, req, result)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}

func (c *Client) createRequest(ctx context.Context, name, method, path string, body interface{}) (*http.Request, error) {
	var (
		buf *bytes.Buffer
		req *http.Request
		err error
	)

	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, c.errorF("can't encode %s request body: %w", name, err)
		}

		req, err = http.NewRequestWithContext(ctx, method, c.createUrl(path), buf)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, c.createUrl(path), nil)
	}

	if err != nil {
		return nil, c.errorF("can't create %s request: %w", name, err)
	}

	// Instrui a não validar XSRF/CSRF
	req.Header.Add("X-Atlassian-Token", "no-check")
	req.Header.Add("Authorization", c.authorization)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/json")

	return req, nil
}

func (c *Client) doRequest(name string, req *http.Request, result interface{}) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, c.errorF("can't do %s request %w", name, err)
	}

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		if err := json.NewDecoder(res.Body).Decode(result); err != nil {
			return nil, c.errorF("can't decode %s response body: %w", name, err)
		}

		fallthrough
	case http.StatusNoContent, http.StatusResetContent:
		return res, nil
	case http.StatusUnauthorized:
		return nil, c.errorF("authentication failed in %s", name)
	case http.StatusServiceUnavailable:
		return nil, c.errorF("service is not available for %s: %s", name, res.Status)
	case http.StatusInternalServerError:
		return nil, c.errorF("internal server error on %s: %s", name, res.Status)
	case http.StatusConflict:
		return nil, c.errorF("conflict on %s: %s", name, res.Status)
	}

	if res.StatusCode > 399 && res.StatusCode < 500 {
		var errorResp errorResponse
		if err := json.NewDecoder(res.Body).Decode(&errorResp); err != nil {
			return nil, c.errorF("%s: can't decode error response for status code '%d': %w", name, res.StatusCode, err)
		}

		return nil, c.errorF("%s: Confluence API return the following error message: '%d' - '%s'", name, res.StatusCode, errorResp.Message)
	}

	return nil, c.errorF("%s: status code returned not mapped: '%d'", name, res.StatusCode)
}

func (c *Client) errorF(format string, a ...interface{}) error {
	return fmt.Errorf("confluence client: "+format, a...)
}

func (c *Client) createUrl(path string) string {
	return fmt.Sprintf("%s/%s", c.host, path)
}
