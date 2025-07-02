package lang_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
)

func TestHttpHandler(t *testing.T) {
	parser := lang.Parser{}
	loop := &painter.Loop{}

	handler := lang.HttpHandler(loop, &parser)

	t.Run("POST valid script", func(t *testing.T) {
		body := "white\nupdate\n"
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("GET valid cmd param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost/?cmd=green", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	t.Run("POST invalid command", func(t *testing.T) {
		body := "invalidcmd\n"
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
		}
	})
}
