package v1

import (
	"practice_ctl/pkg/apimachinery/runtime"
	"practice_ctl/pkg/apimachinery/runtime/schema"
	"strings"
)

func (f *Car) SetGroupVersionKind(kind schema.GroupVersionKind) {
	f.Kind = kind.Kind
	if kind.Group == "" {
		f.ApiVersion = kind.Version
	} else {
		f.ApiVersion = kind.Group + "/" + kind.Version
	}
}

func (f *Car) GroupVersionKind() schema.GroupVersionKind {
	res := strings.Split(f.ApiVersion, "/")
	var s schema.GroupVersionKind
	if len(res) < 2 {
		s = schema.GroupVersionKind{
			Group: "",
			Version: res[0],
			Kind: f.Kind,
		}

	} else {
		s = schema.GroupVersionKind{
			Group: res[0],
			Version: res[1],
			Kind: f.Kind,
		}
	}
	return s
}



func (f *Car) GetObjectKind(g schema.GroupVersionKind) schema.ObjectKind {
	f.SetGroupVersionKind(g)
	return f
}

var SchemeGroupVersion = schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Car"}

var (
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &schemeBuilder
	AddToScheme        = localSchemeBuilder.AddScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	f := &Car{
		ApiVersion: "apps/v1",
		Kind: "Car",
	}
	scheme.AddKnownTypes(SchemeGroupVersion, f)
	return nil
}

func init() {
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &schemeBuilder
	AddToScheme        = localSchemeBuilder.AddScheme
}

var _ schema.ObjectKind = &Car{}
var _ runtime.Object = &Car{}


