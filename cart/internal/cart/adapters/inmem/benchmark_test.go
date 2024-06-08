package inmem_test

import (
	"context"
	"testing"

	"route256/cart/internal/cart/adapters/inmem"
)

func BenchmarkAddItem(b *testing.B) {
	ctx := context.Background()
	count := uint16(5)

	b.Run("AddItem_UniqueUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})

	b.Run("AddItem_SingleUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(5)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})
}

func BenchmarkDeleteItem(b *testing.B) {
	ctx := context.Background()
	count := uint16(5)

	b.Run("DeleteItem_UniqueUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			skuID := int64(1000 + i)
			err := storage.DeleteItem(ctx, userID, skuID)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})

	b.Run("DeleteItem_SingleUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		for i := 0; i < b.N; i++ {
			userID := int64(5)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(5)
			skuID := int64(1000 + i)
			err := storage.DeleteItem(ctx, userID, skuID)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})
}

func BenchmarkDeleteItemsByUserID(b *testing.B) {
	ctx := context.Background()
	count := uint16(5)

	b.Run("DeleteItemsByUserID_UniqueUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			err := storage.DeleteItemsByUserID(ctx, userID)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})
}

func BenchmarkGetItems(b *testing.B) {
	ctx := context.Background()
	count := uint16(5)

	b.Run("GetItems_UniqueUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(i)
			_, err := storage.GetItemsByUserID(ctx, userID)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})

	b.Run("GetItems_SingleUser_UniqueItem", func(b *testing.B) {
		storage := inmem.NewStorage()
		for i := 0; i < b.N; i++ {
			userID := int64(5)
			skuID := int64(1000 + i)
			err := storage.AddItem(ctx, userID, skuID, count)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userID := int64(5)
			_, err := storage.GetItemsByUserID(ctx, userID)
			if err != nil {
				b.Fatalf("AddItem failed: %v", err)
			}
		}
	})
}
