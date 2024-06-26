package inmem

import (
	"context"
	"sync"
	"testing"

	"route256/cart/internal/cart/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestStorage_AddItem(t *testing.T) {
	testItemAlpha := models.Item{
		SkuID: 1000,
		Count: 5,
	}
	testItemBeta := models.Item{
		SkuID: 5000,
		Count: 2,
	}
	tests := []struct {
		name          string
		userID        int64
		amountOfItems int
		testToRun     func(t *testing.T, userID int64, storage *Storage)
	}{
		{
			name:   "AddSingleItem",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
			},
			amountOfItems: 1,
		},
		{
			name:   "AddMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemBeta.SkuID, testItemBeta.Count)
				assert.NoError(t, err)
			},
			amountOfItems: 2,
		},
		{
			name:   "AddSingleItemMultipleTimes",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				assert.Equal(t, testItemAlpha.Count*3, storage.carts[userID].items[testItemAlpha.SkuID])
			},
			amountOfItems: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := NewStorage()

			tt.testToRun(t, tt.userID, storage)
			require.NotNil(t, storage.carts[tt.userID])
			assert.Equal(t, len(storage.carts[tt.userID].items), tt.amountOfItems)
		})
	}
}

func TestStorage_DeleteItem(t *testing.T) {
	testItemAlpha := models.Item{
		SkuID: 1000,
		Count: 5,
	}
	testItemBeta := models.Item{
		SkuID: 5000,
		Count: 2,
	}
	tests := []struct {
		name            string
		userID          int64
		amountBeforeDel int
		amountAfterDel  int
		testToRun       func(t *testing.T, userID int64, storage *Storage)
		wantErr         error
	}{
		{
			name:   "DeleteItem_AddedSingleItem",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
			},
			amountBeforeDel: 1,
			amountAfterDel:  0,
		},
		{
			name:   "DeleteItem_AddedMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemBeta.SkuID, testItemBeta.Count)
				assert.NoError(t, err)
			},
			amountBeforeDel: 2,
			amountAfterDel:  1,
		},
		{
			name:   "DeleteItem_SingleItemAddedMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				assert.Equal(t, testItemAlpha.Count*3, storage.carts[userID].items[testItemAlpha.SkuID])
			},
			amountBeforeDel: 1,
			amountAfterDel:  0,
		},
		{
			name:      "DeleteItemErr",
			userID:    999,
			testToRun: func(_ *testing.T, _ int64, _ *Storage) {},
			wantErr:   models.ErrCartIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := NewStorage()

			if tt.wantErr == nil {
				tt.testToRun(t, tt.userID, storage)
				require.NotNil(t, storage.carts[tt.userID])
				assert.Equal(t, len(storage.carts[tt.userID].items), tt.amountBeforeDel)

				err := storage.DeleteItem(context.Background(), tt.userID, testItemAlpha.SkuID)
				assert.NoError(t, err)

				assert.Equal(t, len(storage.carts[tt.userID].items), tt.amountAfterDel)
			} else {
				err := storage.DeleteItem(context.Background(), tt.userID, testItemAlpha.SkuID)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteItemsByUserID(t *testing.T) {
	testItemAlpha := models.Item{
		SkuID: 1000,
		Count: 5,
	}
	testItemBeta := models.Item{
		SkuID: 5000,
		Count: 2,
	}
	tests := []struct {
		name            string
		userID          int64
		amountBeforeDel int
		testToRun       func(t *testing.T, userID int64, storage *Storage)
		wantErr         error
	}{
		{
			name:   "DeleteItem_AddedSingleItem",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
			},
			amountBeforeDel: 1,
		},
		{
			name:   "DeleteItem_AddedMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemBeta.SkuID, testItemBeta.Count)
				assert.NoError(t, err)
			},
			amountBeforeDel: 2,
		},
		{
			name:   "DeleteItem_SingleItemAddedMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				assert.Equal(t, testItemAlpha.Count*3, storage.carts[userID].items[testItemAlpha.SkuID])
			},
			amountBeforeDel: 1,
		},
		{
			name:      "DeleteItemErr",
			userID:    999,
			testToRun: func(_ *testing.T, _ int64, _ *Storage) {},
			wantErr:   models.ErrCartIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := NewStorage()

			if tt.wantErr == nil {
				tt.testToRun(t, tt.userID, storage)
				require.NotNil(t, storage.carts[tt.userID])
				assert.Equal(t, len(storage.carts[tt.userID].items), tt.amountBeforeDel)

				err := storage.DeleteItemsByUserID(context.Background(), tt.userID)
				assert.NoError(t, err)

				assert.Nil(t, storage.carts[tt.userID])
			} else {
				err := storage.DeleteItemsByUserID(context.Background(), tt.userID)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestStorage_GetItemsByUserID(t *testing.T) {
	testItemAlpha := models.Item{
		SkuID: 1000,
		Count: 5,
	}
	testItemBeta := models.Item{
		SkuID: 5000,
		Count: 2,
	}
	tests := []struct {
		name          string
		userID        int64
		amountOfItems int
		testToRun     func(t *testing.T, userID int64, storage *Storage)
		wantErr       error
		expectItems   []models.Item
	}{
		{
			name:   "GetItems_AddedSingleItem",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
			},
			amountOfItems: 1,
			expectItems:   []models.Item{testItemAlpha},
		},
		{
			name:   "GetItems_AddedMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemBeta.SkuID, testItemBeta.Count)
				assert.NoError(t, err)
			},
			amountOfItems: 2,
			expectItems:   []models.Item{testItemAlpha, testItemBeta},
		},
		{
			name:   "GetItems_SingleItemAddedMultipleItems",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				assert.Equal(t, testItemAlpha.Count*3, storage.carts[userID].items[testItemAlpha.SkuID])
			},
			amountOfItems: 1,
			expectItems:   []models.Item{{SkuID: testItemAlpha.SkuID, Count: testItemAlpha.Count * 3}},
		},
		{
			name:      "GetItemsErr_CartNeverExisted",
			userID:    999,
			testToRun: func(_ *testing.T, _ int64, _ *Storage) {},
			wantErr:   models.ErrCartIsEmpty,
		},
		{
			name:   "GetItemsErr_CartBecameEmpty",
			userID: 999,
			testToRun: func(t *testing.T, userID int64, storage *Storage) {
				t.Helper()
				err := storage.AddItem(context.Background(), userID, testItemAlpha.SkuID, testItemAlpha.Count)
				assert.NoError(t, err)
				err = storage.DeleteItem(context.Background(), userID, testItemAlpha.SkuID)
				assert.NoError(t, err)
			},
			wantErr: models.ErrCartIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := NewStorage()

			if tt.wantErr == nil {
				tt.testToRun(t, tt.userID, storage)
				require.NotNil(t, storage.carts[tt.userID])
				assert.Equal(t, len(storage.carts[tt.userID].items), tt.amountOfItems)

				items, err := storage.GetItemsByUserID(context.Background(), tt.userID)
				assert.NoError(t, err)

				assert.ElementsMatch(t, items, tt.expectItems)
			} else {
				tt.testToRun(t, tt.userID, storage)
				_, err := storage.GetItemsByUserID(context.Background(), tt.userID)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestStorage_AddItem_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	storage := NewStorage()
	wg := sync.WaitGroup{}
	userID := int64(1)
	itemCount := 100

	for i := range itemCount {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := storage.AddItem(context.Background(), userID, int64(i), 1)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()
	items, err := storage.GetItemsByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, itemCount, len(items))
}

func TestStorage_DeleteItem_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	storage := NewStorage()
	userID := int64(1)
	skuID := int64(1)
	count := uint16(1)

	err := storage.AddItem(context.Background(), userID, skuID, count)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	goroutineCount := 10

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//nolint:errcheck // check is not required
			_ = storage.DeleteItem(context.Background(), userID, skuID)
		}()
	}

	wg.Wait()
	items, err := storage.GetItemsByUserID(context.Background(), userID)
	assert.ErrorIs(t, err, models.ErrCartIsEmpty)
	assert.Nil(t, items)
}

func TestStorage_AddDeleteItem_Concurrent(t *testing.T) {
	defer goleak.VerifyNone(t)

	storage := NewStorage()
	wg := sync.WaitGroup{}
	userID := int64(1)
	skuID := int64(1)
	count := uint16(1)

	// Add and delete items concurrently
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				err := storage.AddItem(context.Background(), userID, skuID, count)
				assert.NoError(t, err)
			} else {
				//nolint:errcheck // check is not required
				_ = storage.DeleteItem(context.Background(), userID, skuID)
			}
		}(i)
	}

	wg.Wait()
	items, err := storage.GetItemsByUserID(context.Background(), userID)
	// В зависимости от тайминга может быть разный результат (основной упор теста на проверку с помощью флага -race)
	if err == nil {
		assert.NotNil(t, items)
	} else {
		assert.Equal(t, models.ErrCartIsEmpty, err)
	}
}
