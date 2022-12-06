package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go-do/common/conf"
	"go-do/common/utils"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

func LoadNacos() {
	clientConfig := constant.ClientConfig{
		TimeoutMs:           conf.ConfigInfo.Nacos.ClientConfig.TimeoutMs,
		ListenInterval:      conf.ConfigInfo.Nacos.ClientConfig.ListenInterval,
		NotLoadCacheAtStart: conf.ConfigInfo.Nacos.ClientConfig.NotLoadCacheAtStart,
		LogDir:              conf.ConfigInfo.Nacos.ClientConfig.LogDir,
		BeatInterval:        1000,
		NamespaceId:         "4c22a800-178e-43ea-8bd1-7372e63c5b55",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:   conf.ConfigInfo.Nacos.Ip,
			Port:     conf.ConfigInfo.Nacos.Port,
			GrpcPort: conf.ConfigInfo.Nacos.GrpcPort,
		},
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		panic(err)
	}

	loadConfig(configClient)

	port, err := strconv.Atoi(conf.ConfigInfo.Server.Port)

	if err != nil {
		panic(err)
	}

	grpcPort, err := strconv.Atoi(conf.ConfigInfo.Server.GrpcPort)
	rpcPort, err := strconv.Atoi(conf.ConfigInfo.Server.RpcPort)

	if err != nil {
		panic(err)
	}

	ip, err := utils.GetSubnetIp()

	if err != nil {
		panic(err)
	}

	// http 服务注册到nacos
	success, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip: ip,
		//Ip:          "127.0.0.1",
		Port:        uint64(port),
		ServiceName: conf.ConfigInfo.Server.Name,
		Weight:      conf.ConfigInfo.Nacos.InstanceConfig.Weight,
		//ClusterName: conf.ConfigInfo.Server.Name,
		Enable:    conf.ConfigInfo.Nacos.InstanceConfig.Enable,
		Healthy:   conf.ConfigInfo.Nacos.InstanceConfig.Healthy,
		Ephemeral: conf.ConfigInfo.Nacos.InstanceConfig.Ephemeral,
		Metadata:  map[string]string{"source": "go"},
	})

	if !success {
		panic(err)
	}

	// grpc 服务注册到nacos
	success, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip: ip,
		//Ip:          "127.0.0.1",
		Port:        uint64(grpcPort),
		ServiceName: conf.ConfigInfo.Server.GrpcName,
		Weight:      conf.ConfigInfo.Nacos.InstanceConfig.Weight,
		//ClusterName: conf.ConfigInfo.Server.Name,
		Enable:    conf.ConfigInfo.Nacos.InstanceConfig.Enable,
		Healthy:   conf.ConfigInfo.Nacos.InstanceConfig.Healthy,
		Ephemeral: conf.ConfigInfo.Nacos.InstanceConfig.Ephemeral,
		Metadata:  map[string]string{"source": "go"},
	})

	if !success {
		panic(err)
	}

	// rpc 服务注册到nacos
	success, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip: ip,
		//Ip:          "127.0.0.1",
		Port:        uint64(rpcPort),
		ServiceName: conf.ConfigInfo.Server.RpcName,
		Weight:      conf.ConfigInfo.Nacos.InstanceConfig.Weight,
		//ClusterName: conf.ConfigInfo.Server.Name,
		Enable:    conf.ConfigInfo.Nacos.InstanceConfig.Enable,
		Healthy:   conf.ConfigInfo.Nacos.InstanceConfig.Healthy,
		Ephemeral: conf.ConfigInfo.Nacos.InstanceConfig.Ephemeral,
		Metadata:  map[string]string{"source": "go"},
	})

	if !success {
		panic(err)
	}

	//time.Sleep(20 * time.Second)
	//
	//client.DeregisterInstance(vo.DeregisterInstanceParam{
	//	Ip:          ip,
	//	Port:        uint64(port),
	//	ServiceName: conf.ConfigInfo.Server.Name,
	//	Ephemeral:   conf.ConfigInfo.Nacos.InstanceConfig.Ephemeral,
	//	//Cluster:     "cluster-a", // 默认值DEFAULT
	//	//GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	//})
}

func loadConfig(configClient config_client.IConfigClient) {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "config_dev",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(content)

	b := []byte(content)
	err = yaml.Unmarshal(b, &conf.ConfigInfo)

	if err != nil {

		fmt.Printf(" config parse failed: %s", err)

		os.Exit(-1)

	}

	fmt.Println(conf.ConfigInfo)
}
