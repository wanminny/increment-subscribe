package model

// sku supplierId 组合
type SkuSupplierId struct {
	Id interface{} `json:"id"`
	Sku string `json:"sku"`
	OriginSupplierId interface{} `json:"origin_supplier_id,omitempty"` //omitempty 为空的时候忽略
	SupplierId interface{} `json:"supplier_id"`
	CurrentTime string `json:"current_time"`
}

// 事件类型
type BinLogEventType struct {
	EventType string `json:"event_type"`
}

//一条binlog日志记录
//to mq data structure
type Record struct {
	//BinLogEventType `json:"bin_log_event_type"`
	EventType string `json:"event_type"`
	Rows []SkuSupplierId `json:"rows"`
}