package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"practice_ctl/pkg/apis/core/v1"
)

func NewServerVersionInfo() *v1.Version {
	return &v1.Version{
		Version: "0.1.1",
		GoVersion: "go1.18",
	}
}

type VersionCtl struct {
}

func (v *VersionCtl) Version(c *gin.Context) goft.Json {
	return NewServerVersionInfo()
}

func (v *VersionCtl) Name() string {
	return "VersionCtl"
}

func (v *VersionCtl) Build(goft *goft.Goft) {
	// GET  http://localhost:8080/version
	goft.Handle("GET", "/version", v.Version)

}

