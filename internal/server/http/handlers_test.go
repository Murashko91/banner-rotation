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
	"sync"
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

	t.Run("events API", func(t *testing.T) {
		// insert banners, slots and groups

		for i := 0; i < 5; i++ {
			requestJSON := storage.Banner{
				Descr: fmt.Sprintf("test %d", i),
			}

			urlServer.Path = "banner"
			resp, err := doHTTPCall(urlServer.String(), http.MethodPost, requestJSON)
			require.NoError(t, err, "not expected request error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			resBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err, "not expected read body error")

			var resBannerJSON storage.Banner
			err = json.Unmarshal(resBody, &resBannerJSON)
			require.NoError(t, err, "response unmarshal error")
			require.Equal(t, resBannerJSON.Descr, requestJSON.Descr, "not expected Descr")
			require.NotNil(t, resBannerJSON.ID, "banner id not populated")

			// insert slot
			urlServer.Path = "slot"

			resp, err = doHTTPCall(urlServer.String(), http.MethodPost, requestJSON)
			require.NoError(t, err, "not expected request error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			resBody, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err, "not expected read body error")

			var resSlotJSON storage.Slot
			err = json.Unmarshal(resBody, &resSlotJSON)
			require.NoError(t, err, "response unmarshal error")
			require.Equal(t, resSlotJSON.Descr, requestJSON.Descr, "not expected Descr")
			require.NotNil(t, resSlotJSON.ID, "slot id not populated")

			// insert groups
			urlServer.Path = "group"
			resp, err = doHTTPCall(urlServer.String(), http.MethodPost, requestJSON)
			require.NoError(t, err, "not expected request error")
			require.Equal(t, resp.StatusCode, http.StatusOK)

			resBody, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err, "not expected read body error")

			var resGroupJSON storage.Slot
			err = json.Unmarshal(resBody, &resGroupJSON)
			require.NoError(t, err, "response unmarshal error")
			require.Equal(t, resGroupJSON.Descr, requestJSON.Descr, "not expected Descr")
			require.NotNil(t, resGroupJSON.ID, "slot id not populated")

			// add banners to slot
			urlServer.Path = "banner-rotation"
			requestRotationJSON := storage.Rotation{
				BannerID: resBannerJSON.ID,
				SlotID:   1,
			}

			resp, err = doHTTPCall(urlServer.String(), http.MethodPost, requestRotationJSON)

			require.NoError(t, err, "not expected response error")
			require.Equal(t, resp.StatusCode, http.StatusOK)
			defer resp.Body.Close()
		}

		urlServer.Path = "stat"
		clickReq := storage.Statistic{
			BannerID:      2,
			SosialGroupID: 1,
			SlotID:        1,
		}

		for i := 0; i < 5; i++ {
			resp, err := doHTTPCall(urlServer.String(), http.MethodPost, clickReq)
			require.NoError(t, err, "not expected response error")
			require.Equal(t, resp.StatusCode, http.StatusOK)
			defer resp.Body.Close()
		}

		urlServer.Path = "banner-rotation"
		banditResultMap := make(map[int]int)

		params := urlServer.Query()
		params.Add("slot_id", "1")
		params.Add("group_id", "1")
		urlServer.RawQuery = params.Encode()

		mu := sync.Mutex{}

		for i := 0; i < 25; i++ {
			resp, err := doHTTPCall(urlServer.String(), http.MethodGet, nil)

			require.NoError(t, err, "not expected response error")
			require.Equal(t, resp.StatusCode, http.StatusOK)
			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "not expected read body error")
			defer resp.Body.Close()
			var respBanner storage.Banner
			err = json.Unmarshal(resBody, &respBanner)
			require.NoError(t, err, "response unmarshal error")

			mu.Lock()
			banditResultMap[respBanner.ID]++
			mu.Unlock()
		}
		// get a banner with biggest shows count

		topBannerID := 0
		mostShowsCount := 0

		for bannerID, showsCount := range banditResultMap {
			if showsCount > mostShowsCount {
				topBannerID = bannerID
				mostShowsCount = showsCount
			}
		}

		// most popular banner ID should be 2 because of added clicks

		require.Equal(t, 2, topBannerID)
		stat, err := db.GetBannersStat(context.Background(), 1, 1, []int{topBannerID})
		require.NoError(t, err, "response db stat error")
		require.Equalf(t, stat[0].ShowsCount, mostShowsCount, "shows count has not been updated in db")
	})
}

func doHTTPCall(url string, method string, body any) (*http.Response, error) {
	var data io.Reader

	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		data = bytes.NewReader(reqBody)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, data)
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
