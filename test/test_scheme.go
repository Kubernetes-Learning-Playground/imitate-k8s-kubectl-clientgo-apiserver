package main

import (
	"fmt"
	"practice_ctl/pkg/apimachinery/runtime/schema"
	v1 "practice_ctl/pkg/apis/core/v1"
	"practice_ctl/pkg/apiserver/controllers"
)

func main() {
	// 取得全局scheme表
	sh := controllers.GetGlobalScheme()
	var gvk = schema.GroupVersionKind{Group: "food", Version: "v1", Kind: "Food"}
	res := sh.GetObjectKind(gvk)
	fmt.Println(res)

	var gvk1 = schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Apple"}
	res1 := sh.GetObjectKind(gvk1)
	fmt.Println(res1)
	a := res1.(*v1.Apple)
	fmt.Println(a.Kind, a.ApiVersion)

	fmt.Println(res1.GroupVersionKind())

	var gvk2 = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Car"}
	res2 := sh.GetObjectKind(gvk2).GroupVersionKind()
	fmt.Println(res2)
}
