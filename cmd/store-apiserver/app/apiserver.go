package app

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"k8s.io/klog/v2"
	"net/http"
	"practice_ctl/pkg/storeapiserver/auth"
	"practice_ctl/pkg/storeapiserver/controllers"
	"practice_ctl/pkg/storeapiserver/filters"
	"practice_ctl/pkg/util/helpers"
)

type APIServer struct {

	// server实例
	Server *http.Server
	// Aggregater server
	AggregaterServer *AggregationApiServer

	// Config 配置文件
	Config *Config

	// webservice container, where all webservice defines
	container *restful.Container


}

type Config struct {

}

type ServerRunOptions struct {
	Config
	Port 		 string
	EtcdEndpoint string
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{}
}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*ServerRunOptions
}

// Complete 完成server options配置
func Complete(s *ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	if s.Port == "" {
		s.Port = "8888"
	}
	if s.EtcdEndpoint == "" {
		s.EtcdEndpoint = "127.0.0.1:2379"
	}

	options.ServerRunOptions = s

	return options, nil
}

// TODO: 处理启动的配置校验逻辑
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
		AggregaterServer: NewAggregationApiServer(),
	}


	server := &http.Server{
		Addr: fmt.Sprintf(":%v", s.Port),
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
	s.installAllAPIs()

	// handler chain
	s.buildHandlerChain(s.Server.Handler)

	return nil
}

// installAllAPIs 注册api
func (s *APIServer) installAllAPIs() {
	helpers.Must(s.AddCommonApiToContainer(s.container))
	helpers.Must(s.AddServiceV1ApiToContainer(s.container))
}

func (s *APIServer) AddCommonApiToContainer(container *restful.Container) error {
	ws := new(restful.WebService)

	// 测试接口
	ws.Route(ws.GET("/test").
		To(func(request *restful.Request, response *restful.Response) {
			response.WriteAsJson("hello world")
		})).
		Doc("hello world")
	// ping接口
	ws.Route(ws.GET("/ping").To(func(request *restful.Request, response *restful.Response) {
		response.WriteAsJson("pong")
	})).Doc("keepalive ping")

	// 测试panic回复接口
	ws.Route(ws.GET("/try_panic").To(func(request *restful.Request, response *restful.Response) {
		panic("panic")
	})).Doc("try panic")

	// 测试请求超时接口
	ws.Route(ws.GET("/try_timeout").To(controllers.TimedHandler)).Doc("try timeout")

	// 登入接口
	ws.Route(ws.POST("/login").To(auth.LoginHandler))

	// 注册接口
	ws.Route(ws.POST("/register").To(func(request *restful.Request, response *restful.Response) {
		req := struct {
			Path string `json:"path"`
			Host string `json:"host"`
		}{}
		if err := request.ReadEntity(&req); err != nil {
			fmt.Println("bind json err!")
			errResp := struct {
				Code int    `json:"code"`
				Err  string `json:"err"`
			}{Code: http.StatusBadRequest, Err: err.Error()}
			response.WriteEntity(errResp)
		}

		s.AggregaterServer.AggregationMap[req.Path] = req.Host
		resp := struct {
			Code int    `json:"code"`
			Res  interface{} `json:"res"`
		}{Code: http.StatusOK, Res: req}
		response.WriteAsJson(resp)

	}))

	container.Add(ws)

	return nil
}

func (s *APIServer) AddServiceV1ApiToContainer(container *restful.Container) error {
	serviceWs := new(restful.WebService)
	serviceWs.Path("/v1").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	serviceWs.Route(serviceWs.POST("/apple").To(appleCtl.CreateApple))
	serviceWs.Route(serviceWs.PATCH("/apple").To(appleCtl.PatchApple))
	serviceWs.Route(serviceWs.GET("/applelist").To(appleCtl.ListApple))
	serviceWs.Route(serviceWs.GET("/apple/watch").To(appleCtl.WatchApple))
	serviceWs.Route(serviceWs.GET("/apple").To(appleCtl.GetApple))
	serviceWs.Route(serviceWs.PUT("/apple").To(appleCtl.UpdateApple))
	serviceWs.Route(serviceWs.DELETE("/apple").To(appleCtl.DeleteApple))

	serviceWs.Route(serviceWs.POST("/car").To(carCtl.CreateCar))
	serviceWs.Route(serviceWs.GET("/carlist").To(carCtl.ListCar))
	serviceWs.Route(serviceWs.GET("/car/watch").To(carCtl.WatchCar))
	serviceWs.Route(serviceWs.GET("/car").To(carCtl.GetCar))
	serviceWs.Route(serviceWs.PUT("/car").To(carCtl.UpdateCar))
	serviceWs.Route(serviceWs.PATCH("/car").To(carCtl.PatchCar))
	serviceWs.Route(serviceWs.DELETE("/car").To(carCtl.DeleteCar))

	container.Add(serviceWs)
	return nil
}

// TODO: 注册v1alpha路由

// buildHandlerChain 中间件
func (s *APIServer) buildHandlerChain(apiHandler http.Handler) {
	// TODO: 增加其他中间件，认证 鉴权 准入
	handler := apiHandler

	handler = s.AggregaterServer.SearchHandler(handler)		// 聚合中间件
	handler = filters.AuthorizeMiddleware(handler)          // 授权中间件
	handler = filters.AuthenticateMiddleware(handler)		// 认证中间件
	handler = filters.RequestTimeoutMiddleware(handler)		// 请求超时中间件
	handler = filters.IpLimiterMiddleware(handler)			// ip限流中间件
	handler = filters.LoggerMiddleware(handler)				// 日志中间件
	handler = filters.RecoveryMiddleware(handler)			// panic中间件

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