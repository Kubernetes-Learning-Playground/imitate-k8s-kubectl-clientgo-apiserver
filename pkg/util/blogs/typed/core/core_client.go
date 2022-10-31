package core

import "practice_ctl/pkg/util/blogs/rest"

type CoreInterface interface {
	VersionGetter
}

var _ CoreInterface = &CoreClient{}

type CoreClient struct {
	client *rest.RESTClient
}

func (c *CoreClient) Version() VersionInterface {
	return newVersion(c.client)
}
func New(c *rest.RESTClient) *CoreClient {
	return &CoreClient{client: c}
}
