package models

type ItemDescription struct {
	Name  string
	Price uint32
}

type Item struct {
	SkuID int64
	Count uint16
	Info  ItemDescription
}
