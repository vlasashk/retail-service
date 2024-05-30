package removecart_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"route256/cart/internal/cart/models/constants"
	mockRemover "route256/cart/internal/cart/ports/vanilla/handlers/removecart/mocks"

	"route256/cart/internal/cart/ports/vanilla/handlers/removecart"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mocksToUse struct {
	Remover *mockRemover.CartRemover
}

func initMocks(t *testing.T) *mocksToUse {
	return &mocksToUse{
		Remover: mockRemover.NewCartRemover(t),
	}
}

func TestRemoveCartHandler(t *testing.T) {
	url := "http://example.com/user/%d/cart"

	tests := []struct {
		name       string
		mockSetUp  func(*mocksToUse, int64)
		expectCode int
		userID     int64
		expectResp string
	}{
		{
			name:       "RemoveCartHandlerSuccess",
			expectCode: http.StatusNoContent,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Remover.On("DeleteItemsByUserID", mock.Anything, userID).Return(nil).Once()
			},
			userID: 999,
		},
		{
			name:       "RemoveCartWrongUserID",
			expectCode: http.StatusBadRequest,
			mockSetUp:  func(_ *mocksToUse, _ int64) {},
			userID:     -1,
			expectResp: `{"error":"invalid user_id value"}`,
		},
		{
			name:       "RemoveCartErr",
			expectCode: http.StatusInternalServerError,
			mockSetUp: func(m *mocksToUse, userID int64) {
				m.Remover.On("DeleteItemsByUserID", mock.Anything, userID).Return(errors.New("any error")).Once()
			},
			userID:     13,
			expectResp: `{"error":"failed to remove cart"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mocks := initMocks(t)

			r := httptest.NewRequest("DELETE", fmt.Sprintf(url, tt.userID), nil)
			r.SetPathValue(constants.PathArgUserID, strconv.Itoa(int(tt.userID)))
			w := httptest.NewRecorder()
			tt.mockSetUp(mocks, tt.userID)

			handler := removecart.New(zerolog.Logger{}, mocks.Remover)
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
