package core

import "practice_ctl/pkg/util/stores/rest"

type CoreInterface interface {
	VersionGetter
	AppleGetter
}

var _ CoreInterface = &CoreClient{}

type CoreClient struct {
	client *rest.RESTClient
}

func (c *CoreClient) Version() VersionInterface {
	return newVersion(c.client)
}

func (c *CoreClient) Apple() AppleInterface {
	return newApple(c.client)
}

func New(c *rest.RESTClient) *CoreClient {
	return &CoreClient{client: c}
}
