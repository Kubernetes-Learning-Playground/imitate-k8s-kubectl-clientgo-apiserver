package apps

import "practice_ctl/pkg/util/stores/rest"

type AppsInterface interface {
	CarGetter
}

var _ AppsInterface = &AppsClient{}

type AppsClient struct {
	client *rest.RESTClient
}


func (c *AppsClient) Car() CarInterface {
	return newCar(c.client)
}

func New(c *rest.RESTClient) *AppsClient {
	return &AppsClient{client: c}
}

