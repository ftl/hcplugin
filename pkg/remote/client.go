package remote

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http *http.Client
	port int
}

func New(port int) *Client {
	return &Client{
		http: &http.Client{Timeout: 250 * time.Millisecond},
		port: port,
	}
}

func (c *Client) Do(actionID string) error {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", c.port),
		Path:   "/do",
	}
	q := u.Query()
	q.Set("action", actionID)
	u.RawQuery = q.Encode()
	return c.post(u.String())
}

func (c *Client) Send(text string) error {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", c.port),
		Path:   "/send",
	}
	q := u.Query()
	q.Set("text", text)
	u.RawQuery = q.Encode()
	return c.post(u.String())
}

func (c *Client) post(rawURL string) error {
	resp, err := c.http.Post(rawURL, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("hellocontest: %s", resp.Status)
	}
	return nil
}
