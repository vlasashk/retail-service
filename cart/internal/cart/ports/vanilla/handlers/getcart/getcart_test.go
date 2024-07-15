package getcart_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/models"
	"route256/cart/internal/cart/ports/vanilla/handlers/getcart"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mocksToUse struct {
	Retriever *CartRetrieverMock
}

func initMocks(t *testing.T) *mocksToUse {
	mc := minimock.NewController(t)
	return &mocksToUse{
		Retriever: NewCartRetrieverMock(mc),
	}
}

func TestGetCartHandler(t *testing.T) {
	url := "http://example.com/user/%d/cart/list"

	testCart := models.ItemsInCart{
		Items: []models.Item{
			{
				SkuID: 1000,
				Count: 5,
				Info: models.ItemDescription{
					Name:  "TEST",
					Price: 1000,
				},
			},
			{
				SkuID: 2000,
				Count: 1,
				Info: models.ItemDescription{
					Name:  "TEST",
					Price: 500,
				},
			},
		},
		TotalPrice: 5500,
	}

	tests := []struct {
		name       string
		mockSetUp  func(*mocksToUse, int64)
		expectCode int
		userID     int64
		expectResp string
	}{
		{
			name:       "GetCartHandlerSuccess",
			expectCode: http.StatusOK,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then(testCart, nil)
			},
			userID:     999,
			expectResp: `{"items":[{"sku_id":1000,"name":"TEST","count":5,"price":1000},{"sku_id":2000,"name":"TEST","count":1,"price":500}],"total_price":5500}`,
		},
		{
			name:       "GetCartWrongUserID",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     -1,
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "GetCartIsEmpty",
			expectCode: http.StatusNotFound,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then(models.ItemsInCart{}, models.ErrCartIsEmpty)
			},
			userID:     999,
			expectResp: `{"error":"cart is empty or doesn't exist"}`,
		},
		{
			name:       "GetCartRetrieverErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then(models.ItemsInCart{}, errors.New("any err"))
			},
			userID:     13,
			expectResp: `{"error":"failed to checkout cart"}`,
		},
		{
			name:       "GetCartProductDoesntExist",
			expectCode: http.StatusPreconditionFailed,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.GetItemsByUserIDMock.When(minimock.AnyContext, userID).Then(models.ItemsInCart{}, models.ErrNotFound)
			},
			userID:     42,
			expectResp: `{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mocks := initMocks(t)

			r := httptest.NewRequest("GET", fmt.Sprintf(url, tt.userID), nil)
			r.SetPathValue(constants.UserID, strconv.Itoa(int(tt.userID)))
			w := httptest.NewRecorder()
			tt.mockSetUp(mocks, tt.userID)

			handler := getcart.New(mocks.Retriever)
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectCode, w.Code)
			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()

			assert.JSONEq(t, tt.expectResp, string(body))
		})
	}
}
