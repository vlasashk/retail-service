package getcart_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"route256/cart/internal/cart/models"
	mockProvider "route256/cart/internal/cart/ports/vanilla/handlers/common/mocks"
	"route256/cart/internal/cart/ports/vanilla/handlers/getcart"
	mockRetriever "route256/cart/internal/cart/ports/vanilla/handlers/getcart/mocks"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCartHandler(t *testing.T) {
	url := "http://example.com/user/%d/cart/list"

	retriever := mockRetriever.NewCartRetriever(t)
	provider := mockProvider.NewProductProvider(t)

	testItem := models.Item{
		SkuId: 1000,
		Count: 5,
		Info: models.ItemDescription{
			Name:  "TEST",
			Price: 1000,
		},
	}

	tests := []struct {
		name       string
		mockSetUp  func(int64)
		expectCode int
		userID     int64
		wantBody   assert.ComparisonAssertionFunc
		expectResp string
	}{
		{
			name:       "GetCartHandlerSuccess",
			expectCode: http.StatusOK,
			mockSetUp: func(userID int64) {
				retriever.On("GetItemsByUserID", mock.Anything, userID).Return([]models.Item{testItem, testItem}, nil).Once()
				provider.On("GetProduct", mock.Anything, testItem.SkuId).Return(testItem.Info, nil).Twice()
			},
			userID:     999,
			expectResp: `{"items":[{"sku_id":1000,"name":"TEST","count":5,"price":1000},{"sku_id":1000,"name":"TEST","count":5,"price":1000}],"total_price":10000}`,
		},
		{
			name:       "GetCartWrongUserID",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ int64) {},
			userID:     -1,
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "GetCartIsEmpty",
			expectCode: http.StatusNotFound,
			mockSetUp: func(userID int64) {
				retriever.On("GetItemsByUserID", mock.Anything, userID).Return(nil, models.ErrCartIsEmpty).Once()
			},
			userID:     999,
			expectResp: `{"error":"cart is empty or doesn't exist"}`,
		},
		{
			name:       "GetCartRetrieverErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(userID int64) {
				retriever.On("GetItemsByUserID", mock.Anything, userID).Return(nil, errors.New("any err")).Once()
			},
			userID:     13,
			expectResp: `{"error":"failed to get items"}`,
		},
		{
			name:       "GetCartProductDoesntExist",
			expectCode: http.StatusNotFound,
			mockSetUp: func(userID int64) {
				retriever.On("GetItemsByUserID", mock.Anything, userID).Return([]models.Item{testItem}, nil).Once()
				provider.On("GetProduct", mock.Anything, testItem.SkuId).Return(models.ItemDescription{}, models.ErrCartIsEmpty).Once()
			},
			userID:     42,
			expectResp: `{"error":"cart is empty or doesn't exist"}`,
		},
		{
			name:       "GetCartProviderErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(userID int64) {
				retriever.On("GetItemsByUserID", mock.Anything, userID).Return([]models.Item{testItem}, nil).Once()
				provider.On("GetProduct", mock.Anything, testItem.SkuId).Return(models.ItemDescription{}, errors.New("any error")).Once()
			},
			userID:     1,
			expectResp: `{"error":"failed to checkout cart"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", fmt.Sprintf(url, tt.userID), nil)
			r.SetPathValue("user_id", strconv.Itoa(int(tt.userID)))
			w := httptest.NewRecorder()
			tt.mockSetUp(tt.userID)

			handler := getcart.New(zerolog.Logger{}, retriever, provider)
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectCode, w.Code)
			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()

			bodyActual := strings.TrimSpace(string(body))
			assert.Equal(t, tt.expectResp, bodyActual)
		})
	}
}
