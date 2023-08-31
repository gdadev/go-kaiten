package kaiten

import (
	"net/http"
	"net/url"
)

type Client struct {
	baseUrl *url.URL

	//
	Spaces *SpacesService
}

func (c *Client) NewRequest(method, path string, opt interface{}) (*http.Request, error) {
	return nil, nil
}

type Response struct {
	*http.Response
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	return nil, nil
}
