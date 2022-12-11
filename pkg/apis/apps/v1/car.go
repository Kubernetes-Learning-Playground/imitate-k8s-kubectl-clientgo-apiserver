package v1

type Car struct {
	Brand   string	`json:"brand"`
	Price  string	`json:"price"`
	Color  string	`json:"color"`
	Name   string	`json:"name"`
}

type CarList struct {
	Item []*Car
}
