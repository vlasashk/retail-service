package getcart

import "route256/cart/internal/cart/models"

type userCartResp struct {
	Items      []itemResp `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type itemResp struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

func itemsToDTO(items models.ItemsInCart) userCartResp {
	resp := userCartResp{
		TotalPrice: items.TotalPrice,
	}
	resp.Items = make([]itemResp, 0, len(items.Items))
	for _, item := range items.Items {
		resp.Items = append(resp.Items, itemResp{
			SkuID: item.SkuID,
			Name:  item.Info.Name,
			Count: item.Count,
			Price: item.Info.Price,
		})
	}
	return resp
}
