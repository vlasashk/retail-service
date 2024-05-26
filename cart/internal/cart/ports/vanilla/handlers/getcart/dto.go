package getcart

import "route256/cart/internal/cart/models"

type userCartResp struct {
	Items      []itemResp `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type itemResp struct {
	SkuId int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

func cartToDTO(cart models.Cart) userCartResp {
	resp := userCartResp{
		TotalPrice: cart.TotalPrice,
	}
	resp.Items = make([]itemResp, 0, len(cart.Items))
	for _, item := range cart.Items {
		resp.Items = append(resp.Items, itemResp{
			SkuId: item.SkuId,
			Name:  item.Info.Name,
			Count: item.Count,
			Price: item.Info.Price,
		})
	}
	return resp
}
