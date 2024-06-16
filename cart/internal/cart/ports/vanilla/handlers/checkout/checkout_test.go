package checkout_test

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
	"route256/cart/internal/cart/ports/vanilla/handlers/checkout"

	"github.com/gojuno/minimock/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mocksToUse struct {
	Retriever *CartCheckoutMock
}

func initMocks(t *testing.T) *mocksToUse {
	mc := minimock.NewController(t)
	return &mocksToUse{
		Retriever: NewCartCheckoutMock(mc),
	}
}

func TestCheckoutCartHandler(t *testing.T) {
	url := "http://example.com/user/%d/cart/checkout"

	orderIDSample := int64(111)

	tests := []struct {
		name       string
		mockSetUp  func(*mocksToUse, int64)
		expectCode int
		userID     int64
		expectResp string
	}{
		{
			name:       "CartCheckoutHandlerSuccess",
			expectCode: http.StatusOK,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.CartCheckoutMock.When(minimock.AnyContext, userID).Then(orderIDSample, nil)
			},
			userID:     999,
			expectResp: `{"order_id":111}`,
		},
		{
			name:       "CartCheckoutWrongUserID",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     -1,
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "CartCheckoutIsEmpty",
			expectCode: http.StatusNotFound,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.CartCheckoutMock.When(minimock.AnyContext, userID).Then(0, models.ErrCartIsEmpty)
			},
			userID:     999,
			expectResp: `{"error":"cart is empty or doesn't exist"}`,
		},
		{
			name:       "CartCheckoutRetrieverErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.CartCheckoutMock.When(minimock.AnyContext, userID).Then(0, errors.New("any err"))
			},
			userID:     13,
			expectResp: `{"error":"failed to checkout cart"}`,
		},
		{
			name:       "CartCheckoutProductDoesntExist",
			expectCode: http.StatusPreconditionFailed,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.CartCheckoutMock.When(minimock.AnyContext, userID).Then(0, models.ErrNotFound)
			},
			userID:     42,
			expectResp: `{"error":"not found"}`,
		},
		{
			name:       "CartCheckoutInsufficientStock",
			expectCode: http.StatusPreconditionFailed,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Retriever.CartCheckoutMock.When(minimock.AnyContext, userID).Then(0, models.ErrInsufficientStock)
			},
			userID:     42,
			expectResp: `{"error":"insufficient stock"}`,
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

			handler := checkout.New(zerolog.Logger{}, mocks.Retriever)
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
