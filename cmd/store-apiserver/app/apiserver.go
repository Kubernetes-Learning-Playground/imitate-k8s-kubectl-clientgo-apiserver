package app

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/klog/v2"
	"net/http"
	"practice_ctl/pkg/storeapiserver/controllers"
)

type APIServer struct {
	ServerCount int

	Server *http.Server

	//Config *apiserverconfig.Config
	Config *Config

	// webservice container, where all webservice defines
	container *restful.Container


}

type Config struct {

}

type ServerRunOptions struct {
	Config
	port string
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{}
}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*ServerRunOptions
}

func Complete(s *ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	s.port = "8888"


	options.ServerRunOptions = s

	return options, nil
}

func (s *ServerRunOptions) Validate() []error {
	var errs []error


	return errs
}

// Run runs the specified APIServer.  This should never exit.
func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {

	server, err := completeOptions.NewAPIServer(stopCh)
	if err != nil {
		return err
	}

	err = server.PrepareRun(stopCh)
	if err != nil {
		return err
	}

	return server.Run(context.Background())
}

func (s *ServerRunOptions) NewAPIServer(stopCh <-chan struct{}) (*APIServer, error) {
	apiServer := &APIServer{
		Config:     &s.Config,
	}


	server := &http.Server{
		Addr: fmt.Sprintf(":%v", s.port),
	}

	apiServer.Server = server

	return apiServer, nil
}


func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.container = restful.NewContainer()

	// 设定路由为CurlyRouter(快速路由)
	s.container.Router(restful.CurlyRouter{})

	for _, ws := range s.container.RegisteredWebServices() {
		klog.V(2).Infof("%s", ws.RootPath())
	}

	// container作为http server的handler
	s.Server.Handler = s.container

	// 注册服务
	s.installAIscopeAPIs()

	// handler chain
	s.buildHandlerChain(stopCh)

	return nil
}

func (s *APIServer) installAIscopeAPIs() {
	err := AddToContainer(s.container)
	err = AddKVServiceToContainer(s.container)
	if err != nil {
		panic(err)
	}
}

func AddToContainer(container *restful.Container) error {
	ws := new(restful.WebService)

	ws.Route(ws.GET("/hello").
		To(func(request *restful.Request, response *restful.Response) {


			response.WriteAsJson("hello world")
		})).
		Doc("hello world")

	container.Add(ws)

	return nil
}

func AddKVServiceToContainer(container *restful.Container) error {
	kvWs := new(restful.WebService)
	kvWs.Path("/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	kvWs.Route(kvWs.POST("/apple").To(appleCtl.CreateApple))
	kvWs.Route(kvWs.GET("/applelist").To(appleCtl.ListApple))
	kvWs.Route(kvWs.GET("/apple/watch").To(appleCtl.WatchApple))
	kvWs.Route(kvWs.GET("/apple").To(appleCtl.GetApple))
	kvWs.Route(kvWs.PUT("/apple").To(appleCtl.UpdateApple))
	kvWs.Route(kvWs.DELETE("/apple").To(appleCtl.DeleteApple))

	kvWs.Route(kvWs.POST("/car").To(carCtl.CreateCar))
	kvWs.Route(kvWs.GET("/carlist").To(carCtl.ListCar))
	kvWs.Route(kvWs.GET("/car/watch").To(carCtl.WatchCar))
	kvWs.Route(kvWs.GET("/car").To(carCtl.GetCar))
	kvWs.Route(kvWs.PUT("/car").To(carCtl.UpdateCar))
	kvWs.Route(kvWs.DELETE("/car").To(carCtl.DeleteCar))

	container.Add(kvWs)
	return nil
}


func (s *APIServer) buildHandlerChain(stopCh <-chan struct{}) {

	handler := s.Server.Handler

	s.Server.Handler = handler
}



func (s *APIServer) Run(ctx context.Context) (err error) {
	initController()
	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
	}()

	klog.V(0).Infof("Start listening on %s", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return err
}

var (
	appleCtl *controllers.AppleRestfulCtl
	carCtl   *controllers.CarRestfulCtl
)

func initController() {
	appleCtl = controllers.NewAppleRestfulCtl()
	carCtl = controllers.NewCarRestfulCtl()
}