package remote

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serverPort(t *testing.T, serverURL string) int {
	t.Helper()
	u, err := url.Parse(serverURL)
	require.NoError(t, err)
	_, portStr, err := net.SplitHostPort(u.Host)
	require.NoError(t, err)
	port, err := strconv.Atoi(portStr)
	require.NoError(t, err)
	return port
}

func TestDo(t *testing.T) {
	var capturedPath, capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(serverPort(t, srv.URL))
	err := c.Do("entry.next_esm_step", nil)
	require.NoError(t, err)
	assert.Equal(t, "/do", capturedPath)
	assert.Equal(t, "action=entry.next_esm_step", capturedQuery)
}

func TestDoURLEncoding(t *testing.T) {
	var capturedAction string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAction = r.URL.Query().Get("action")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(serverPort(t, srv.URL))
	require.NoError(t, c.Do("hello world", nil))
	assert.Equal(t, "hello world", capturedAction)
}

func TestDoWithParams(t *testing.T) {
	var capturedQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(serverPort(t, srv.URL))
	err := c.Do("bandmap.mark_with_number", map[string]string{"number": "5", "action": "ignored"})
	require.NoError(t, err)
	assert.Equal(t, "bandmap.mark_with_number", capturedQuery.Get("action"))
	assert.Equal(t, "5", capturedQuery.Get("number"))
}

func TestSend(t *testing.T) {
	var capturedPath, capturedText string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedText = r.URL.Query().Get("text")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(serverPort(t, srv.URL))
	err := c.Send("hello world")
	require.NoError(t, err)
	assert.Equal(t, "/send", capturedPath)
	assert.Equal(t, "hello world", capturedText)
}

func TestDoServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := New(serverPort(t, srv.URL))
	err := c.Do("some.action", nil)
	assert.Error(t, err)
}

func TestDoTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(serverPort(t, srv.URL))
	err := c.Do("some.action", nil)
	assert.Error(t, err)
}
