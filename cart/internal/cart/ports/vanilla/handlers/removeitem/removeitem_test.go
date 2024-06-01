package removeitem_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"route256/cart/internal/cart/constants"
	"route256/cart/internal/cart/ports/vanilla/handlers/removeitem"
	mockRemover "route256/cart/internal/cart/ports/vanilla/handlers/removeitem/mocks"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocksToUse struct {
	Remover *mockRemover.ItemRemover
}

func initMocks(t *testing.T) *mocksToUse {
	return &mocksToUse{
		Remover: mockRemover.NewItemRemover(t),
	}
}

func TestRemoveCartHandler(t *testing.T) {
	url := "http://example.com/user/%d/cart/%d"

	tests := []struct {
		name       string
		mockSetUp  func(*mocksToUse, int64, int64)
		expectCode int
		userID     int64
		skuID      int64
		expectResp string
	}{
		{
			name:       "RemoveItemHandlerSuccess",
			expectCode: http.StatusNoContent,
			mockSetUp: func(m *mocksToUse, userID, skuID int64) {
				m.Remover.On("DeleteItem", mock.Anything, userID, skuID).Return(nil).Once()
			},
			userID: 999,
			skuID:  999,
		},
		{
			name:       "RemoveItemWrongUserID",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _, _ int64) {},
			userID:     -1,
			skuID:      1000,
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "RemoveItemWrongSKUid",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _, _ int64) {},
			userID:     999,
			skuID:      -1,
			expectResp: `{"error":"invalid sku_id value"}`,
		},
		{
			name:       "RemoveItemWrongUserIDAndSKUid",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _, _ int64) {},
			userID:     -1,
			skuID:      -1,
			expectResp: `{"error":"invalid user_id value\ninvalid sku_id value"}`,
		},
		{
			name:       "RemoveItemErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID, skuID int64) {
				m.Remover.On("DeleteItem", mock.Anything, userID, skuID).Return(errors.New("any error")).Once()
			},
			userID:     13,
			skuID:      999,
			expectResp: `{"error":"failed to remove item"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mocks := initMocks(t)

			r := httptest.NewRequest("DELETE", fmt.Sprintf(url, tt.userID, tt.skuID), nil)
			r.SetPathValue(constants.UserID, strconv.Itoa(int(tt.userID)))
			r.SetPathValue(constants.SKUid, strconv.Itoa(int(tt.skuID)))
			w := httptest.NewRecorder()
			tt.mockSetUp(mocks, tt.userID, tt.skuID)

			handler := removeitem.New(zerolog.Logger{}, mocks.Remover)
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
