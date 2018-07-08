

package main

import (
	"fmt"
	"encoding/json"
)

type OrderInfo struct{

	OrderId string `json:"order_id"`
	OrderGoodsId string `json:"order_goods_id"`
	Total int `json:"total"`
	userId string `json:"user_id"`
}


var gOrderInfo []*OrderInfo 
func main() {
	
	orderOne := &OrderInfo{
		OrderId : "12345",
		OrderGoodsId : "45633",
		Total:11,
		userId:"uid-111",
	}

	gOrderInfo = append(gOrderInfo,orderOne)

	fmt.Printf("%v\n",gOrderInfo)
	fmt.Printf("%T\n",gOrderInfo)

	info,err:= json.Marshal(gOrderInfo)

	if err != nil{
		panic(fmt.Sprintf("json error %s",err))
	}
	fmt.Println(string(info))

	json.Unmarshal((info),gOrderInfo)

	for k,v := range gOrderInfo{
		fmt.Println(k,v)
	}
	// fmt.Println(gOrderInfo)

}




