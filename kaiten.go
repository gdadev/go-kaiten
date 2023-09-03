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

	// if resp.StatusCode == http.StatusUnauthorized && c.authType == BasicAuth {
	// 	resp.Body.Close()
	// 	// The token most likely expired, so we need to request a new one and try again.
	// 	if _, err := c.requestOAuthToken(req.Context(), basicAuthToken); err != nil {
	// 		return nil, err
	// 	}
	// 	return c.Do(req, v)
	// }
	// defer resp.Body.Close()

	// // If not yet configured, try to configure the rate limiter
	// // using the response headers we just received. Fail silently
	// // so the limiter will remain disabled in case of an error.
	// c.configureLimiterOnce.Do(func() { c.configureLimiter(req.Context(), resp.Header) })

	// response := newResponse(resp)

	// err = CheckResponse(resp)
	// if err != nil {
	// 	// Even though there was an error, we still return the response
	// 	// in case the caller wants to inspect it further.
	// 	return response, err
	// }

	response := &Response{Response: resp}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}
