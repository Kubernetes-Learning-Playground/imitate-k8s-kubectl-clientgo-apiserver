package runtime

import "practice_ctl/pkg/apimachinery/runtime/schema"

type Object interface {
	GetObjectKind(g schema.GroupVersionKind) schema.ObjectKind
}
