package kaiten

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	baseUrlTemplate = `https://%v.kaiten.ru/api/v1/`
)

type Client struct {
	client *retryablehttp.Client

	domain string

	baseUrl *url.URL

	token string

	//
	Spaces *SpacesService
}

func NewClient(domain, token string) (*Client, error) {
	c := &Client{
		domain: domain,
		token:  token,
	}

	//nolint:errcheck
	c.setBaseUrl(baseUrlTemplate, domain)

	c.client = retryablehttp.NewClient()
	c.client.RetryMax = 0
	c.client.Logger = nil

	c.Spaces = &SpacesService{client: c}

	return c, nil
}

func (c *Client) setBaseUrl(baseUrlTemplate, domain string) error {
	if !strings.HasSuffix(baseUrlTemplate, "/") {
		baseUrlTemplate += "/"
	}

	baseUrlRaw := fmt.Sprintf(baseUrlTemplate, domain)

	baseUrl, err := url.Parse(baseUrlRaw)
	if err != nil {
		return err
	}

	c.baseUrl = baseUrl
	return nil
}

func (c *Client) NewRequest(method, path string, opt interface{}) (*retryablehttp.Request, error) {
	u := *c.baseUrl
	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseUrl.Path + path
	u.Path = c.baseUrl.Path + unescaped

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	// if c.UserAgent != "" {
	// 	reqHeaders.Set("User-Agent", c.UserAgent)
	// }

	var body interface{}
	switch {
	case method == http.MethodPatch || method == http.MethodPost || method == http.MethodPut:
		reqHeaders.Set("Content-Type", "application/json")

		if opt != nil {
			body, err = json.Marshal(opt)
			if err != nil {
				return nil, err
			}

		}

	case opt != nil:
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// for _, fn := range append(c.defaultRequestOptions, options...) {
	// 	if fn == nil {
	// 		continue
	// 	}
	// 	if err := fn(req); err != nil {
	// 		return nil, err
	// 	}
	// }

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

type Response struct {
	*http.Response
}

func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := &Response{Response: resp}

	if err = checkResponse(resp); err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

type ErrorResponse struct {
	Body     []byte
	Response *http.Response
	Message  string
}

func (e *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(e.Response.Request.URL.Path)
	u := fmt.Sprintf("%s://%s%s", e.Response.Request.URL.Scheme, e.Response.Request.URL.Host, path)
	return fmt.Sprintf("%s %s: %d %s", e.Response.Request.Method, u, e.Response.StatusCode, e.Message)
}

type errorResponseBody struct {
	Message string `json:"message"`
}

func checkResponse(r *http.Response) error {
	if r.StatusCode == 200 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Body = data

		var raw errorResponseBody
		if err := json.Unmarshal(data, &raw); err != nil {
			errorResponse.Message = fmt.Sprintf("failed to parse unknown error format: %s", data)
		} else {
			errorResponse.Message = raw.Message
		}
	}

	return errorResponse
}
