package meta

import (
	"practice_ctl/pkg/apis/meta/types"
	"time"
)

type TypeMeta struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
}

type ObjectMeta struct {
	Name            string            `json:"name" yaml:"name"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	UID             types.UID         `json:"uid,omitempty" yaml:"UID,omitempty"`
	ResourceVersion string            `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
	CreateTimestamp time.Time         `json:"createTimestamp,omitempty" yaml:"createTimestamp,omitempty"`
	UpdateTimestamp time.Time         `json:"updateTimestamp,omitempty" yaml:"updateTimestamp,omitempty"`
	// TODO: ResourceVersion UID CreateTimestamp UpdateTimestamp Labels 增加

}
