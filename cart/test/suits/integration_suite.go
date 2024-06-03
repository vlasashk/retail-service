//go:build integration

package suits

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"route256/cart/config"
	"route256/cart/internal/cart"
	"route256/cart/internal/cart/constants"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

const (
	defaultUserID      = 111111
	defaultAlphaItemID = 2618151
	defaultBetaItemID  = 2956315
	additionalItemID   = 1076963
)

type IntegrationSuite struct {
	suite.Suite
	client         http.Client
	serviceAddress string
	cancel         context.CancelFunc
}

func (s *IntegrationSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	cfg, err := config.New()
	if err != nil {
		s.T().Fatal(err)
	}

	s.serviceAddress = "http://" + cfg.Address
	s.client = http.Client{
		Timeout: time.Second,
	}

	go func() {
		if err = cart.Run(ctx, cfg); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	if err = s.healthCheck(10); err != nil {
		s.T().Fatal(err)
	}
}

func (s *IntegrationSuite) TearDownSuite() {
	s.cancel()
}

func (s *IntegrationSuite) healthCheck(attempts int) error {
	var err error

	healthURL := s.serviceAddress + "/healthz"

	for attempts > 0 {
		if _, err = s.client.Get(healthURL); err != nil {
			log.Debug().Int("attempts left", attempts).Str("URL", healthURL).Msg("Service is not available for integration rests")
			time.Sleep(time.Second)
			attempts--
			continue
		}
		return nil
	}

	return err
}

func (s *IntegrationSuite) SetupTest() {
	s.addItemHelper(defaultUserID, defaultAlphaItemID, 5)
	s.addItemHelper(defaultUserID, defaultBetaItemID, 1)
}

func (s *IntegrationSuite) TearDownTest() {
	cleanUP, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%d/cart", s.serviceAddress, defaultUserID), nil)
	s.Require().NoError(err)
	cleanUP.SetPathValue(constants.UserID, strconv.Itoa(defaultUserID))

	resp, err := s.client.Do(cleanUP)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *IntegrationSuite) addItemHelper(userID, skuID, count int) {
	addItem, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%d/cart/%d", s.serviceAddress, userID, skuID),
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"count":%d}`, count))))
	s.Require().NoError(err)
	addItem.SetPathValue(constants.UserID, strconv.Itoa(userID))
	addItem.SetPathValue(constants.SKUid, strconv.Itoa(skuID))

	resp, err := s.client.Do(addItem)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *IntegrationSuite) delItemHelper(userID, skuID int) {
	delItem, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%d/cart/%d", s.serviceAddress, userID, skuID), nil)
	s.Require().NoError(err)
	delItem.SetPathValue(constants.UserID, strconv.Itoa(userID))
	delItem.SetPathValue(constants.SKUid, strconv.Itoa(skuID))

	resp, err := s.client.Do(delItem)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *IntegrationSuite) TestAddItem() {
	tests := []struct {
		name       string
		expectCode int
		userID     int64
		skuID      int64
		body       io.Reader
		expectResp string
	}{
		{
			name:       "AddItemHandlerSuccess",
			expectCode: http.StatusOK,
			userID:     defaultUserID,
			skuID:      1076963,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
		},
		{
			name:       "AddItemWrongUserID",
			expectCode: http.StatusBadRequest,
			userID:     -1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "AddItemWrongSKUid",
			expectCode: http.StatusBadRequest,
			userID:     999,
			skuID:      -1,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"invalid sku_id value"}`,
		},
		{
			name:       "AddItemWrongUserIDAndSKUid",
			expectCode: http.StatusBadRequest,
			userID:     -1,
			skuID:      -1,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"invalid user_id value\ninvalid sku_id value"}`,
		},
		{
			name:       "AddItemWrongCount",
			expectCode: http.StatusBadRequest,
			userID:     1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":0}`)),
			expectResp: `{"error":"invalid amount of products"}`,
		},
		{
			name:       "AddItemBadCountBody",
			expectCode: http.StatusBadRequest,
			userID:     defaultUserID,
			skuID:      defaultAlphaItemID,
			body:       bytes.NewBuffer([]byte(`{"count":-10}`)),
			expectResp: `{"error":"failed to process request body"}`,
		},
		{
			name:       "AddItemProductDoesntExist",
			expectCode: http.StatusPreconditionFailed,
			userID:     defaultUserID,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"item not found"}`,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			r, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%d/cart/%d", s.serviceAddress, tt.userID, tt.skuID), tt.body)
			s.Require().NoError(err)
			r.SetPathValue(constants.UserID, strconv.Itoa(int(tt.userID)))
			r.SetPathValue(constants.SKUid, strconv.Itoa(int(tt.skuID)))

			resp, err := s.client.Do(r)
			s.Require().NoError(err)
			s.Equal(tt.expectCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)
			_ = resp.Body.Close()
			if len(body) > 0 || len(tt.expectResp) > 0 {
				s.JSONEq(tt.expectResp, string(body))
			}
		})
	}
}

func (s *IntegrationSuite) TestGetItems() {
	tests := []struct {
		name       string
		setUP      func()
		tearDown   func()
		expectCode int
		userID     int64
		expectResp string
	}{
		{
			name:       "GetCartHandlerSuccess",
			expectCode: http.StatusOK,
			userID:     defaultUserID,
			expectResp: `{"items":[{"sku_id":2618151,"name":"Пора снимать бикини","count":5,"price":452},{"sku_id":2956315,"name":"Eloy. Time To Turn","count":1,"price":3130}],"total_price":5390}`,
			setUP:      func() {},
			tearDown:   func() {},
		},
		{
			name:       "GetCartHandlerSuccess_AfterNewAdd",
			expectCode: http.StatusOK,
			userID:     defaultUserID,
			expectResp: `{"items":[{"sku_id":1076963,"name":"Теория нравственных чувств | Смит Адам","count":10,"price":3379},{"sku_id":2618151,"name":"Пора снимать бикини","count":5,"price":452},{"sku_id":2956315,"name":"Eloy. Time To Turn","count":1,"price":3130}],"total_price":39180}`,
			setUP: func() {
				s.addItemHelper(defaultUserID, additionalItemID, 10)
			},
			tearDown: func() {
				s.delItemHelper(defaultUserID, additionalItemID)
			},
		},
		{
			name:       "GetCartHandlerSuccess_AfterAddIdentical",
			expectCode: http.StatusOK,
			userID:     defaultUserID,
			expectResp: `{"items":[{"sku_id":2618151,"name":"Пора снимать бикини","count":5,"price":452},{"sku_id":2956315,"name":"Eloy. Time To Turn","count":2,"price":3130}],"total_price":8520}`,
			setUP: func() {
				s.addItemHelper(defaultUserID, defaultBetaItemID, 1)
			},
			tearDown: func() {
				s.delItemHelper(defaultUserID, defaultBetaItemID)
				// to restore initial test state
				s.addItemHelper(defaultUserID, defaultBetaItemID, 1)
			},
		},
		{
			name:       "GetCartWrongUserID",
			expectCode: http.StatusBadRequest,
			userID:     -1,
			expectResp: `{"error":"invalid user_id value"}`,
			setUP:      func() {},
			tearDown:   func() {},
		},
		{
			name:       "GetCartDoesntExist",
			expectCode: http.StatusNotFound,
			userID:     999,
			expectResp: `{"error":"cart is empty or doesn't exist"}`,
			setUP:      func() {},
			tearDown:   func() {},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setUP()
			r, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%d/cart/list", s.serviceAddress, tt.userID), nil)
			s.Require().NoError(err)
			r.SetPathValue(constants.UserID, strconv.Itoa(int(tt.userID)))

			resp, err := s.client.Do(r)
			s.Require().NoError(err)
			s.Equal(tt.expectCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			s.Require().NoError(err)
			_ = resp.Body.Close()
			if len(body) > 0 || len(tt.expectResp) > 0 {
				s.JSONEq(tt.expectResp, string(body))
			}
			tt.tearDown()
		})
	}
}
