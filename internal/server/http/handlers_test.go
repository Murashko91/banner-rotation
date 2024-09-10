package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/otus-murashko/banners-rotation/internal/app"
	"github.com/otus-murashko/banners-rotation/internal/storage"
	"github.com/otus-murashko/banners-rotation/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestEventHTTPApi(t *testing.T) {
	db := memory.NewMemoryStorage()
	if err := db.Connect(); err != nil {
		t.Errorf("storage connect error %s", err.Error())
	}

	bannerApp := app.New(db)
	h := Handler{
		app: bannerApp,
	}

	bannerRouter := http.NewServeMux()

	bannerRouter.Handle("/banner-rotation", loggingMiddleware(http.HandlerFunc(h.bannerRotationHandler)))
	bannerRouter.Handle("/banner", loggingMiddleware(http.HandlerFunc(h.bannerHandler)))
	bannerRouter.Handle("/slot", loggingMiddleware(http.HandlerFunc(h.slotHandler)))
	bannerRouter.Handle("/group", loggingMiddleware(http.HandlerFunc(h.groupHandler)))
	bannerRouter.Handle("/stat", loggingMiddleware(http.HandlerFunc(h.statHandler)))

	server := httptest.NewServer(bannerRouter)
	urlServer, _ := url.Parse(server.URL)
	urlServer.Path = "banner"

	t.Run("events API", func(t *testing.T) {
		// insert banners

		for i := 0; i < 30; i++ {
			requestJSON := storage.Banner{
				Descr: fmt.Sprintf("test %d", i),
			}

			reqBody, err := json.Marshal(requestJSON)
			require.NoError(t, err, "not expected body marshal error")
			resp, err := doHTTPCall(urlServer.String(), http.MethodPost, bytes.NewReader(reqBody))
			require.NoError(t, err, "not expected request error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			resBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err, "not expected read body error")

			var resJSON storage.Banner
			err = json.Unmarshal(resBody, &resJSON)
			require.NoError(t, err, "response unmarshal error")
			require.Equal(t, resJSON.Descr, requestJSON.Descr, "not expected Descr")
			require.NotNil(t, resJSON.ID, "banner id not populated")
		}

		// insert slots

		// insert groups

		// insert banner to slot

		// Tests clicks

		// Test banner bandit

	})
}

func doHTTPCall(url string, method string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}

	// set content-type header to JSON
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// create HTTP client and execute request
	client := &http.Client{}

	return client.Do(req)
}
