package fireworks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewFromEnvDefaultsToIntegrationWithoutKey(t *testing.T) {
	t.Setenv("FIREWORKS_API_BASE", "")
	t.Setenv("FIREWORKS_BASE_URL", "")
	t.Setenv("FIREWORKS_API_KEY", "")
	c, err := NewFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if c.baseURL != IntegrationBase || c.apiKey != "" {
		t.Fatalf("got base %q, key set %v", c.baseURL, c.apiKey != "")
	}
}

func TestNewFromEnvIntegrationDropsKey(t *testing.T) {
	t.Setenv("FIREWORKS_API_BASE", "")
	t.Setenv("FIREWORKS_API_KEY", "should-not-be-sent")
	c, err := NewFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if c.apiKey != "" {
		t.Fatal("integration route must not carry a key")
	}
}

func TestNewFromEnvDirectNeedsKey(t *testing.T) {
	t.Setenv("FIREWORKS_API_BASE", DirectBase)
	t.Setenv("FIREWORKS_API_KEY", "")
	if _, err := NewFromEnv(); err == nil || !strings.Contains(err.Error(), "FIREWORKS_API_KEY") {
		t.Fatalf("want missing-key error, got %v", err)
	}
	t.Setenv("FIREWORKS_API_KEY", "k")
	if _, err := NewFromEnv(); err != nil {
		t.Fatal(err)
	}
}

func TestAuthorizationOnlyWithKey(t *testing.T) {
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Header.Get("Authorization"))
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"ok"}}]}`))
	}))
	defer srv.Close()
	for _, key := range []string{"", "k"} {
		c := &Client{baseURL: srv.URL, apiKey: key, http: srv.Client()}
		if _, err := c.Chat(context.Background(), "m", "s", "u", 1); err != nil {
			t.Fatal(err)
		}
	}
	if got[0] != "" || got[1] != "Bearer k" {
		t.Fatalf("authorization headers %q", got)
	}
}

func TestIsIntegration(t *testing.T) {
	if !IsIntegration(IntegrationBase) || IsIntegration(DirectBase) || IsIntegration("http://127.0.0.1:9/v1") {
		t.Fatal("IsIntegration misclassified a base")
	}
}
