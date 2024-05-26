package models

type ItemDescription struct {
	Name  string
	Price uint32
}

type Item struct {
	SkuId int64
	Count uint16
	Info  ItemDescription
}
