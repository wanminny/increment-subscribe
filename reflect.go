


package main

import (
	"fmt"
	"reflect"
)


func main() {
	u:= &User{"张三",20}
	t:=reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	fmt.Println(t)	
	fmt.Println(v)


	fmt.Println()

  	fmt.Printf("%T\n",u)
	fmt.Printf("%v\n",u)

	//将原来的对象还原
	if u,ok := v.Interface().(User); ok{
		fmt.Println(u,ok)
	}



	//reflect.Type
	// t0 := t.Type()
	v0 := t.Kind()
	fmt.Println(v0)


	// reflect.Value
	t1 := v.Type()
	v1 := v.Kind()
	fmt.Println(t1,v1,reflect.ValueOf("前缀"))



	method := v.MethodByName("Test")

	method.Call(nil)
	// for i := 0;i < t.Elem().NumField(); i++ {
	// 	fmt.Println("====>",t.Field(i).Name)
	// }

	for j := 0;j < t.NumMethod();j++ {
		fmt.Println("++++",t.Method(j).Name)
	}

//ok{
//	fmt.Println(u)
	// fmt.Printf("%T",u)
//}

//fmt.Println(v,ok)

}
type User struct{
	Name string
	Age int
}


func (*User)Test(){

	fmt.Println("method invoke !")
}

func (u *User)Test2(){

}

func (u User)Test3() {
	
}

