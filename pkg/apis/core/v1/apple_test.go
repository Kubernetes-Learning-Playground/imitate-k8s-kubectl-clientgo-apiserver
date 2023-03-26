package v1

import (
	"fmt"
	"practice_ctl/pkg/apimachinery/runtime"
	"practice_ctl/pkg/apimachinery/runtime/schema"
	"testing"
)

var testScheme = runtime.NewScheme()


func TestApple(t *testing.T) {
	err := AddToScheme(testScheme)
	if err != nil {
		return
	}

	var gvk1 = schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Apple"}
	res1, err := testScheme.GetObjectKind(gvk1)
	fmt.Println(res1, err)

	aa, _ := res1.GetObjectKind(gvk1)
	fmt.Println(aa.GroupVersionKind())
}
