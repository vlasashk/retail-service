package additem

type itemCountReq struct {
	UserID int64  `json:"-"`
	SKUid  int64  `json:"-"`
	Count  uint16 `json:"count" validate:"required,gt=0"`
}
