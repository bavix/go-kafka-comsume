package retailrocket

import (
	"bytes"
	"errors"
	"github.com/bavix/go-kafka-consume/configure"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ApiKey    = "apiKey"
	PartnerId = "partnerId"
)

type Client struct {
	baseURL    string
	apiKey     string
	partnerId  string
	httpClient *http.Client
}

func NewClient(config *configure.RetailRocketConfig) *Client {
	return &Client{
		baseURL:   config.URL,
		apiKey:    config.ApiKey,
		partnerId: config.PartnerID,
		httpClient: &http.Client{
			Timeout: time.Second,
		},
	}
}

func (c *Client) getUrl(method string) (string, error) {
	baseUrl := strings.Trim(c.baseURL, "/")
	method = strings.Trim(method, "/")
	urlObject, err := url.Parse(baseUrl + "/" + method + "/")
	if err != nil {
		return "", err
	}

	query := urlObject.Query()
	query.Add(ApiKey, c.apiKey)
	query.Add(PartnerId, c.partnerId)

	urlObject.RawQuery = query.Encode()

	return urlObject.String(), nil
}

func (c *Client) Post(postMessage *Message) error {
	urlString, err := c.getUrl(postMessage.Method)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer([]byte(postMessage.Body))
	request, err := http.NewRequest("POST", urlString, buffer)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return errors.New("error sending data")
	}

	return nil
}
