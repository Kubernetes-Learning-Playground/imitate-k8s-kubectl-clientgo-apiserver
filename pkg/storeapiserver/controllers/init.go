package controllers

import (
	corev1 "practice_ctl/pkg/apis/core/v1"
	appv1 "practice_ctl/pkg/apis/apps/v1"
	"practice_ctl/pkg/apimachinery/runtime"
)

var (
	//schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &GlobalSchemeBuilder
	AddToScheme        = localSchemeBuilder.AddScheme
)

var GlobalScheme = runtime.NewScheme()

// GlobalSchemeBuilder 初始化 需要注册
var GlobalSchemeBuilder = runtime.SchemeBuilder{
	corev1.AddToScheme,
	appv1.AddToScheme,
}

func GetGlobalScheme() *runtime.Scheme {
	return GlobalScheme
}

func init() {
	err := AddToScheme(GlobalScheme)
	if err != nil {
		panic("scheme error!!")
	}
}