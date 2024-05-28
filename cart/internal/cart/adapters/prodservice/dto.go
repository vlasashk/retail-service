package prodservice

type getProductReq struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type getProductResp struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
