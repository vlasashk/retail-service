package additem

type itemCountReq struct {
	Count uint16 `json:"count" validate:"required,gt=0"`
}
