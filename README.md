# go-do

a mix many tools framework of microservice
 

# httpServer

main.go
```
server.Start("your config path", "your service global prefix", func(r *register.Register) {
		controller.DeviceController(r) // register your controller
	})
  ```
  
  
  DeviceController.go
```
func DeviceController(r *register.Register) {
	r.RegisterController("your controller prefix")

	r.RegisterHttpHandle("service path", enum.GET(request method), func(ctx context.Context, req interface{}) (interface{}, error) {
		tp, ok := req.(*params.TestParams)
		if ok {
			return deviceService.GetDeviceInfo(tp) // your logic code
		}
		return nil, errors.New(string(common.PARAMS_ERROR))
	}(handle), &params.TestParams{}(input params))
}
  ```


# api gateway


main.go
```
server.StartGateway("config path", func(chain *chain.Chain) {
		chain.Apply("JwtFilter", middleware.JwtFilter) // 自定义中间件 (中间件名称，中间件方法)
		chain.Next()
	})
```

JwtFilter.go
```
func JwtFilter(c *chain.Chain) {
	token := c.Token
	if token == "" {
		for _, resource := range common.FRONTEND_RESOURCES {
			if match, _ := regexp.MatchString(resource, c.Url); match {
				c.Next()
				return
			}
		}

		c.Denied()
		return
	}
	c.Next()
}
```

# config file

```
package conf

type config struct {
	Server server `yaml:"server"`

	DataSource dataSource `yaml:"dataSource"`

	Nacos nacos `yaml:"nacos"`

	Wx wx `yaml:"wx"`

	Jwt jwt `yaml:"jwt"`

	GateWay gateway `yaml:"gateway"`
}

type gateway struct {
	Routers []routers `yaml:"routers"`
}

type routers struct {
	ServerName string `yaml:"serverName"`

	Path string `yaml:"path"`

	Filter []string `yaml:"filter"`
}

type jwt struct {
	TokenHeaderName string `yaml:"tokenHeaderName"`
	SecretKey       string `yaml:"secretKey"`
}

type wx struct {
	Miniapp miniapp `yaml:"miniapp"`
}

type miniapp struct {
	Configs []miniConfig `yaml:"configs"`
}

type miniConfig struct {
	AppId  string `yaml:"appid"`
	Secret string `yaml:"secret"`
	Token  string `yaml:"token"`
}

type nacos struct {
	Ip             string         `yaml:"ip"`
	Port           uint64         `yaml:"port"`
	GrpcPort       uint64         `yaml:"grpcPort"`
	ClientConfig   clientConfig   `yaml:"clientConfig"`
	InstanceConfig instanceConfig `yaml:"instanceConfig"`
	ConfigCenter   configCenter   `yaml:"configCenter"`
}

type configCenter struct {
	DataId string `yaml:"dataId"`
}

type clientConfig struct {
	TimeoutMs           uint64 `yaml:"timeoutMs"`
	ListenInterval      uint64 `yaml:"listenInterval"`
	NotLoadCacheAtStart bool   `yaml:"notLoadCacheAtStart"`
	LogDir              string `yaml:"logDir"`
}

type instanceConfig struct {
	Weight    float64 `yaml:"weight"`
	Enable    bool    `yaml:"enable"`
	Healthy   bool    `yaml:"healthy"`
	Ephemeral bool    `yaml:"ephemeral"`
}

type server struct {
	GatewayPort       string `yaml:"gatewayPort"`
	Port              string `yaml:"port"`
	GrpcPort          string `yaml:"grpcPort"`
	Name              string `yaml:"name"`
	GrpcName          string `yaml:"grpcName"`
	RpcPort           string `yaml:"rpcPort"`
	RpcName           string `yaml:"rpcName"`
	HystrixStreamPort string `yaml:"hystrixStreamPort"`
}

type dataSource struct {
	Mongo mongo `yaml:"mongo"`
	Mysql mysql `yaml:"mysql"`
	Redis redis `yaml:"redis"`
}

type mongo struct {
	Uri string `yaml:"uri"`
	Db  string `yaml:"db"`
}

type mysql struct {
	Uri                       string `yaml:"uri"`
	DefaultStringSize         uint   `yaml:"defaultStringSize"`
	DisableDatetimePrecision  bool   `yaml:"disableDatetimePrecision"`
	DontSupportRenameIndex    bool   `yaml:"dontSupportRenameIndex"`
	DontSupportRenameColumn   bool   `yaml:"dontSupportRenameColumn"`
	SkipInitializeWithVersion bool   `yaml:"skipInitializeWithVersion"`
	MaxIdleConns              int    `yaml:"maxIdleConns"`
	MaxOpenConns              int    `yaml:"maxOpenConns"`
	ConnMaxLifetime           int64  `yaml:"connMaxLifetime"`
}

type redis struct {
	Url    string `yaml:"url"`
	Passwd string `yaml:"passwd"`
	Db     int    `yaml:"db"`
}

```
