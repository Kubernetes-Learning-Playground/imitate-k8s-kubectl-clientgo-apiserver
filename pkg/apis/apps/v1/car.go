package v1

import metav1 "practice_ctl/pkg/apis/meta"

type Car struct {
	metav1.TypeMeta   `json:""`
	metav1.ObjectMeta   `json:"metadata" yaml:"metadata"`
	Spec       CarSpec   `json:"spec" yaml:"spec"`
	Status     CarStatus `json:"status" yaml:"status"`
}


type CarSpec struct {
	Brand string `json:"brand" yaml:"brand"`
	Price string `json:"price" yaml:"price"`
	Color string `json:"color" yaml:"color"`
}

type CarStatus struct {
	//CreateTime time.Time
	Status     string
}


type CarList struct {
	Item []*Car
}
