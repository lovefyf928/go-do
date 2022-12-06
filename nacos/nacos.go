package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go-do/common/conf"
	"go-do/common/utils"
	"go-do/nacos/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

var (
	RpcPool *pool.Pool

	GrpcPool *pool.Pool

	nacosClient naming_client.INamingClient
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

	RpcPool = pool.NewRpcPool(getRpcInstance, 300, time.Hour*60)

	GrpcPool = pool.NewGrpcPool(getInstance, 300, time.Hour*60)

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

func GetHttpLbHost(serverName string) (string, error) {
	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serverName,
		GroupName:   "DEFAULT_GROUP",
		Clusters:    []string{"DEFAULT"},
	})
	//update(serverName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", instance.Ip, instance.Port), nil
}

func getInstance(serverName string) (conn *grpc.ClientConn, err error) {

	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: string(serverName),
		GroupName:   "DEFAULT_GROUP",
		Clusters:    []string{"DEFAULT"},
	})
	if err != nil {
		return nil, err
	}

	//_ = fmt.Sprintf("获取到的实例IP:%s;端口:%d", instance., instance.Port)
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getRpcInstance(serverName string) (conn *rpc.Client, err error) {

	instance, err := nacosClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: string(serverName),
		GroupName:   "DEFAULT_GROUP",
		Clusters:    []string{"DEFAULT"},
	})
	if err != nil {
		return
	}

	//_ = fmt.Sprintf("获取到的实例IP:%s;端口:%d", instance., instance.Port)
	//conn, err = grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err = rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", instance.Ip, instance.Port))

	if err != nil {
		return
	}

	return
}

func loadConfig(configClient config_client.IConfigClient) {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: conf.ConfigInfo.Nacos.ConfigCenter.DataId,
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
