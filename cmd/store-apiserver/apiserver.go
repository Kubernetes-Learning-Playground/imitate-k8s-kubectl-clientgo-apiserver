package main

import (
	"github.com/gin-gonic/gin"
	"practice_ctl/pkg/storeapiserver/controllers"
)

func main() {
	router := gin.New()

	defer func() {
		_ = router.Run(":8888")
	}()

	register(router)
}

var (
	appleCtl *controllers.AppleCtl
	versionCtl *controllers.VersionCtl
	carCtl     *controllers.CarCtl
)

func initController() {
	appleCtl = controllers.NewAppleCtl()
	versionCtl = controllers.NewVersionCtl()
	carCtl = controllers.NewCarCtl()
}

func register(router *gin.Engine) {
	initController()
	r := router.Group("/v1")


	{
		// version
		r.GET("/version", versionCtl.Version)

		// apple
		r.GET("/apple", appleCtl.GetApple)
		r.GET("/apple/watch", appleCtl.WatchApple)
		r.GET("/applelist", appleCtl.ListApple)
		r.POST("/apple", appleCtl.CreateApple)
		r.PUT("/apple", appleCtl.UpdateApple)
		r.DELETE("/apple", appleCtl.DeleteApple)

		// car
		r.GET("/car", carCtl.GetCar)
		r.GET("/car/watch", carCtl.WatchCar)
		r.GET("/carlist", carCtl.ListCar)
		r.POST("/car", carCtl.CreateCar)
		r.PUT("/car", carCtl.UpdateCar)
		r.DELETE("/car", carCtl.DeleteCar)


	}
}

