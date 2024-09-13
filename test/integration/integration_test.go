//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/otus-murashko/banners-rotation/internal/config"
	internalhttp "github.com/otus-murashko/banners-rotation/internal/server/http"
	"github.com/otus-murashko/banners-rotation/internal/storage"
	"github.com/stretchr/testify/suite"
)

var configFile string

const intgTest = "Integration test"

func init() {
	flag.StringVar(&configFile, "conf", "./conf/config.yaml", "Path to configuration file")
}

type BannerSuite struct {
	suite.Suite
	ctx    context.Context
	db     *sqlx.DB
	client *http.Client
	path   string
}

func (s *BannerSuite) SetupSuite() {
	flag.Parse()

	config := config.GetBannersConfig(configFile)

	s.client = &http.Client{}

	s.ctx = context.Background()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBName)
	db, err := sqlx.Open("pgx", psqlInfo)
	s.Require().NoError(err)
	s.db = db
	s.path = fmt.Sprintf("http://%s:%d", config.Server.Host, config.Server.Port)
}

func (s *BannerSuite) SetupTest() {
}

// execute after each test.
func (s *BannerSuite) TearDownTest() {
	query := ` DELETE FROM  statistic s USING banner b
		WHERE b.id = s.banner AND b.Descr like '%s%s%s';`
	query = fmt.Sprintf(query, "%", intgTest, "%")

	_, err := s.db.Exec(query)
	s.Require().NoError(err)

	query = ` DELETE FROM  rotation r USING banner b
		WHERE b.id = r.banner AND b.Descr like '%s%s%s'`
	query = fmt.Sprintf(query, "%", intgTest, "%")
	_, err = s.db.Exec(query)
	s.Require().NoError(err)

	query = `DELETE FROM social_group WHERE Descr like '%s%s%s'`
	query = fmt.Sprintf(query, "%", intgTest, "%")
	_, err = s.db.Exec(query)
	s.Require().NoError(err)

	query = `DELETE FROM slot  WHERE Descr like '%s%s%s'`
	query = fmt.Sprintf(query, "%", intgTest, "%")
	_, err = s.db.Exec(query)
	s.Require().NoError(err)

	query = `DELETE FROM banner  WHERE Descr like '%s%s%s'`
	query = fmt.Sprintf(query, "%", intgTest, "%")
	_, err = s.db.Exec(query)
	s.Require().NoError(err)

}

// will run after all the tests in the suite have been run.
func (s *BannerSuite) TearDownSuite() {
	defer s.db.Close()
}

func TestBannerPost(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}

func (s *BannerSuite) Test_CreateBanner() {
	firstSlotID := 0
	firstGroupID := 0
	bannerIDs := make([]int, 0, 5)
	u, err := url.Parse(s.path)
	for i := 0; i < 5; i++ {

		// Test Banner Creation
		requestJSON := internalhttp.Item{
			Descr: fmt.Sprintf("%s test %d", intgTest, i),
		}

		s.Require().NoError(err)

		u.Path = "banner"
		resp, err := doHTTPCall(u.String(), http.MethodPost, requestJSON, s.client)
		s.Require().NoError(err)
		s.Equalf(resp.StatusCode, http.StatusOK, "not expected response status")

		resBody, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		s.Require().NoError(err, "not expected read body error")

		var resBannerJSON internalhttp.ItemResult
		err = json.Unmarshal(resBody, &resBannerJSON)
		s.Require().NoError(err, "response unmarshal error")
		s.Require().Equal(resBannerJSON.Descr, requestJSON.Descr, "not expected Descr")
		s.Require().NotNil(resBannerJSON.ID, "banner id not populated")

		bannerIDs = append(bannerIDs, resBannerJSON.ID)

		// Test if banner exists in DB
		dbBanner, err := getBannerFromStorage(s.db, resBannerJSON.ID)
		s.Require().NoError(err, "db error")
		s.Require().Equal(resBannerJSON.Descr, dbBanner.Descr, "not expected Descr")

		// Test Slot Creation

		u.Path = "slot"
		resp, err = doHTTPCall(u.String(), http.MethodPost, requestJSON, s.client)
		s.Require().NoError(err)
		s.Equalf(resp.StatusCode, http.StatusOK, "not expected response status")

		resBody, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		s.Require().NoError(err, "not expected read body error")

		var resSlotJSON internalhttp.ItemResult
		err = json.Unmarshal(resBody, &resSlotJSON)
		s.Require().NoError(err, "response unmarshal error")
		s.Require().Equal(resSlotJSON.Descr, requestJSON.Descr, "not expected Descr")
		s.Require().NotNil(resSlotJSON.ID, "slot id not populated")

		if firstSlotID == 0 {
			firstSlotID = resSlotJSON.ID
		}

		// Test if slot exists in DB
		dbSlot, err := getSlotFromStorage(s.db, resSlotJSON.ID)
		s.Require().NoError(err, "db error")
		s.Require().Equal(resSlotJSON.Descr, dbSlot.Descr, "not expected Descr")

		// Test Group Creation

		u.Path = "group"
		resp, err = doHTTPCall(u.String(), http.MethodPost, requestJSON, s.client)
		s.Require().NoError(err)
		s.Equalf(resp.StatusCode, http.StatusOK, "not expected response status")

		resBody, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		s.Require().NoError(err, "not expected read body error")

		var resGroupJSON internalhttp.ItemResult
		err = json.Unmarshal(resBody, &resGroupJSON)
		s.Require().NoError(err, "response unmarshal error")
		s.Require().Equal(resGroupJSON.Descr, requestJSON.Descr, "not expected Descr")
		s.Require().NotNil(resGroupJSON.ID, "group id not populated")
		if firstGroupID == 0 {
			firstGroupID = resGroupJSON.ID
		}

		// Test if group exists in DB
		group, err := getGroupFromStorage(s.db, resGroupJSON.ID)
		s.Require().NoError(err, "db error")
		s.Require().Equal(resGroupJSON.Descr, group.Descr, "not expected Descr")

		// Test Add Banner to Slot

		u.Path = "banner-rotation"
		reqRotationJSON := storage.Rotation{
			BannerID: resBannerJSON.ID,
			SlotID:   firstSlotID,
		}
		resp, err = doHTTPCall(u.String(), http.MethodPost, reqRotationJSON, s.client)
		s.Require().NoError(err)
		s.Equalf(resp.StatusCode, http.StatusOK, "not expected response status")

		u.Path = "stat"
		// Add clicks to second banner
		if len(bannerIDs) > 1 {

			bannerStatToClick := internalhttp.Statistic{
				BannerID:      bannerIDs[1],
				SlotID:        firstSlotID,
				SocialGroupID: firstGroupID,
			}
			resp, err = doHTTPCall(u.String(), http.MethodPost, bannerStatToClick, s.client)
			s.Require().NoError(err)
			s.Equalf(resp.StatusCode, http.StatusOK, "not expected response status")
		}

	}

	// get best banner

	u.Path = "banner-rotation"
	banditResultMap := make(map[int]int)

	params := u.Query()
	params.Add("slot_id", strconv.Itoa(firstSlotID))
	params.Add("group_id", strconv.Itoa(firstGroupID))
	u.RawQuery = params.Encode()

	mu := sync.Mutex{}

	for i := 0; i < 25; i++ {
		resp, err := doHTTPCall(u.String(), http.MethodGet, nil, s.client)

		s.Require().NoError(err)
		s.Equalf(resp.StatusCode, http.StatusOK, "not expected response status")
		resBody, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		defer resp.Body.Close()
		var respBanner storage.Banner
		err = json.Unmarshal(resBody, &respBanner)
		s.Require().NoError(err)

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
	s.Require().Equal(bannerIDs[1], topBannerID)

}

func doHTTPCall(url string, method string, body any, client *http.Client) (*http.Response, error) {
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

	return client.Do(req)
}

func getBannerFromStorage(db *sqlx.DB, id int) (storage.Banner, error) {
	sql := `SELECT id,  descr
		FROM Banner
		WHERE id = $1`

	row := db.QueryRowx(sql, id)
	if row.Err() != nil {
		return storage.Banner{}, row.Err()
	}

	var qBanner storage.Banner

	err := row.StructScan(&qBanner)
	if err != nil {
		return storage.Banner{}, err
	}

	return qBanner, nil
}

func getSlotFromStorage(db *sqlx.DB, id int) (storage.Slot, error) {
	sql := `SELECT id,  descr
		FROM slot
		WHERE id = $1`

	row := db.QueryRowx(sql, id)
	if row.Err() != nil {
		return storage.Slot{}, row.Err()
	}

	var qSlot storage.Slot

	err := row.StructScan(&qSlot)
	if err != nil {
		return storage.Slot{}, err
	}

	return qSlot, nil
}

func getGroupFromStorage(db *sqlx.DB, id int) (storage.SocialGroup, error) {
	sql := `SELECT id,  descr
		FROM social_group
		WHERE id = $1`

	row := db.QueryRowx(sql, id)
	if row.Err() != nil {
		return storage.SocialGroup{}, row.Err()
	}

	var qGroup storage.SocialGroup

	err := row.StructScan(&qGroup)
	if err != nil {
		return storage.SocialGroup{}, err
	}

	return qGroup, nil
}
