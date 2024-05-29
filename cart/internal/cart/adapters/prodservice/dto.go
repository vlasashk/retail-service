package prodservice

type getProductReq struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type getProductResp struct {
	Name  string `json:"name" validate:"required"`
	Price uint32 `json:"price" validate:"required,gt=0"`
}
