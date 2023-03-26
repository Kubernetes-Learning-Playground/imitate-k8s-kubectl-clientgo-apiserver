package runtime

import (
	"fmt"
	"practice_ctl/pkg/apimachinery/runtime/schema"
	"reflect"
)


type Scheme struct {
	typeToGVK map[reflect.Type][]schema.GroupVersionKind
	object    map[schema.GroupVersionKind]Object
}



func (s *Scheme) AddKnownTypes(gvk schema.GroupVersionKind, obj Object) {
	t := reflect.TypeOf(obj)
	if len(gvk.Version) == 0 {
		panic(fmt.Sprintf("version is required on all types: %s %v", gvk, t))
	}
	if t.Kind() != reflect.Ptr {
		panic("All types must be pointers to structs.")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("All types must be pointers to structs.")
	}
	s.typeToGVK[t] = append(s.typeToGVK[t], gvk)
	s.object[gvk] = obj
}

func (s *Scheme) GetObjectKind(g schema.GroupVersionKind) schema.ObjectKind {
	obj, ok := s.object[g]
	if !ok {
		return nil
	}
	return obj.GetObjectKind(g)
}

func NewScheme() *Scheme {
	return &Scheme{
		typeToGVK: map[reflect.Type][]schema.GroupVersionKind{},
		object:   map[schema.GroupVersionKind]Object{},
	}
}

type SchemeBuilder []func(s *Scheme) error

func (sb *SchemeBuilder) AddScheme(s *Scheme) error {
	for _, f := range *sb {
		if err := f(s); err != nil {
			return err
		}
	}
	return nil
}

func (sb *SchemeBuilder) Register(funcs ...func(s *Scheme) error) {
	for _, f := range funcs {
		*sb = append(*sb, f)
	}
}

func NewSchemeBuilder(funcs ...func(*Scheme) error) SchemeBuilder {
	sb := SchemeBuilder{}
	sb.Register(funcs...)
	return sb
}
