package v1

type Apple struct {
	Size   string	`json:"size"`
	Price  string	`json:"price"`
	Place  string	`json:"place"`
	Color  string	`json:"color"`
	Name   string	`json:"name"`
}

type AppleList struct {
	Item []*Apple
}
