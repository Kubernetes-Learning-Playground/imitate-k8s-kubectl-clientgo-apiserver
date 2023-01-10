package v1

type Apple struct {
	ApiVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind 	   string  		`json:"kind" yaml:"kind"`
	Name       string		`json:"name" yaml:"name"`
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
