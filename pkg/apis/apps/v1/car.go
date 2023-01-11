package v1

type Car struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   `json:"metadata" yaml:"metadata"`
	Spec       CarSpec   `json:"spec" yaml:"spec"`
	Status     CarStatus `json:"status" yaml:"status"`
}


type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type CarSpec struct {
	Brand string `json:"brand" yaml:"brand"`
	Price string `json:"price" yaml:"price"`
	Color string `json:"color" yaml:"color"`
	Name  string `json:"name" yaml:"name"`
}

type CarStatus struct {
	//CreateTime time.Time
	Status     string
}


type CarList struct {
	Item []*Car
}
