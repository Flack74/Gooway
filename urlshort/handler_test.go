package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapHandler_Redirect(t *testing.T) {
	paths := map[string]string{
		"/foo": "https://example.com/foo",
	}
	fallbackCalled := false
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fallbackCalled = true
	})

	handler := MapHandler(paths, fallback)

	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("expected status %d, got %d", http.StatusMovedPermanently, resp.StatusCode)
	}
	if fallbackCalled {
		t.Errorf("fallback should not have been called for known path")
	}

	// Test fallback
	req2 := httptest.NewRequest("GET", "/unknown", nil)
	w2 := httptest.NewRecorder()
	handler(w2, req2)

	resp2 := w2.Result()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected fallback to respond with 200, got %d", resp2.StatusCode)
	}
	if !fallbackCalled {
		t.Errorf("fallback should have been called for unknown path")
	}
}
