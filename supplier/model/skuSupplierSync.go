package model

type SkuSupplierSync struct {
	Id                      int64
	SupplierId              int64
	SupplierName            string
	SupplierUrl             string
	SupplierSyncTime        string
	LastUpdateTime          string
	IsDel                   int8
	SupplierTopCategoryId   int64
	SupplierTopCategoryName string
}
