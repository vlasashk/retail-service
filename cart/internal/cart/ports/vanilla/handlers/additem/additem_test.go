package additem_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/additem"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mocksToUse struct {
	Adder *CartAdderMock
}

type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func initMocks(t *testing.T) *mocksToUse {
	mc := minimock.NewController(t)
	return &mocksToUse{
		Adder: NewCartAdderMock(mc),
	}
}

func TestAddItemHandler(t *testing.T) {
	url := "http://example.com/user/%d/cart/%d"

	testItem := models.Item{
		SkuID: 1000,
		Count: 5,
		Info: models.ItemDescription{
			Name:  "TEST",
			Price: 1000,
		},
	}

	tests := []struct {
		name       string
		mockSetUp  func(*mocksToUse, int64)
		expectCode int
		userID     int64
		skuID      int64
		body       io.Reader
		expectResp string
	}{
		{
			name:       "AddItemHandlerSuccess",
			expectCode: http.StatusOK,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(nil)
			},
			userID: 999,
			skuID:  1000,
			body:   bytes.NewBuffer([]byte(`{"count":5}`)),
		},
		{
			name:       "AddItemWrongUserID",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     -1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "AddItemWrongSKUid",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     999,
			skuID:      -1,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"invalid sku_id value"}`,
		},
		{
			name:       "AddItemWrongUserIDAndSKUid",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     -1,
			skuID:      -1,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"invalid user_id value\ninvalid sku_id value"}`,
		},
		{
			name:       "AddItemWrongCount",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":0}`)),
			expectResp: `{"error":"invalid amount of products"}`,
		},
		{
			name:       "AddItemBadCountBody",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":-10}`)),
			expectResp: `{"error":"failed to process request body"}`,
		},
		{
			name:       "AddItemProductDoesntExist",
			expectCode: http.StatusPreconditionFailed,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(models.ErrNotFound)
			},
			userID:     42,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"not found"}`,
		},
		{
			name:       "AddItemAdderErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(errors.New("any error"))
			},
			userID:     13,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"failed to add item"}`,
		},
		{
			name:       "AddItemProductProviderErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(models.ErrItemProvider)
			},
			userID:     1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"failed to request item info"}`,
		},
		{
			name:       "AddItemStockProviderErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(models.ErrStockProvider)
			},
			userID:     1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"failed to request stock info"}`,
		},
		{
			name:       "AddItemInsufficientStock",
			expectCode: http.StatusPreconditionFailed,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Adder.AddItemMock.When(minimock.AnyContext, userID, testItem.SkuID, testItem.Count).Then(models.ErrInsufficientStock)
			},
			userID:     1,
			skuID:      1000,
			body:       bytes.NewBuffer([]byte(`{"count":5}`)),
			expectResp: `{"error":"insufficient stock"}`,
		},
		{
			name:       "AddItemReaderErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     1,
			skuID:      1000,
			body:       &errorReader{},
			expectResp: `{"error":"failed to read body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mocks := initMocks(t)

			r := httptest.NewRequest("POST", fmt.Sprintf(url, tt.userID, tt.skuID), tt.body)
			r.SetPathValue(constants.UserID, strconv.Itoa(int(tt.userID)))
			r.SetPathValue(constants.SKUid, strconv.Itoa(int(tt.skuID)))
			w := httptest.NewRecorder()
			tt.mockSetUp(mocks, tt.userID)

			handler := additem.New(zerolog.Logger{}, mocks.Adder)
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectCode, w.Code)
			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			if len(body) > 0 || len(tt.expectResp) > 0 {
				assert.JSONEq(t, tt.expectResp, string(body))
			}
		})
	}
}
