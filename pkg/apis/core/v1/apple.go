package v1

import (
	metav1 "practice_ctl/pkg/apis/meta"
)

type Apple struct {
	metav1.TypeMeta   		`json:"" yaml:"",inline`
	metav1.ObjectMeta   	`json:"metadata" yaml:"metadata"`
	Spec 	   AppleSpec    `json:"spec" yaml:"spec"`
	Status     AppleStatus  `json:"status" yaml:"status"`
}

type AppleSpec struct {
	Size   	   string			`json:"size" yaml:"size"`
	Price  	   string			`json:"price" yaml:"price"`
	Place      string			`json:"place" yaml:"place"`
	Color      string			`json:"color" yaml:"color"`
}

type AppleStatus struct {
	//CreateTime time.Time
	Status     string
}

type AppleList struct {
	Item []*Apple
}
