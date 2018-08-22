package model

type SkuSupplierRelation struct {
	Id                  int64
	Sku                 string
	ProductSupplierCode string
	SupplierId          int64
	IsDel               int8
	LastUpdateTime      string
}
