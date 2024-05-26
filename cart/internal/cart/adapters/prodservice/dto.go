package prodservice

type getProductReq struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type getProductResp struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

//type skuListReq struct {
//	Token string `json:"token"`
//	Start uint32 `json:"startAfterSku"`
//	Count uint32 `json:"count"`
//}
//
//type skuListResp struct {
//	SKUs []uint32 `json:"skus"`
//}
