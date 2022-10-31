package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"practice_ctl/pkg/apis/core/unverstioned"
)

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334

func NewServerVersionInfo() *unverstioned.Version {
	return &unverstioned.Version{Version: "0.1.1", GoVersion: "go1.18"}
}

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
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

//本课程来自 程序员在囧途(www.jtthink.com) 咨询群：98514334
