package middleware

import (
	"errors"
	"go-do/common/conf"
	"go-do/nacos"
	"net/url"
	"strings"
)

func getPathByRawUrl(rawUrl *url.URL) string {
	path := strings.Split(rawUrl.Path, "/")[1]

	path = "/" + path

	return path
}

func getModuleProxyPath(rawUrl *url.URL) (string, error) {

	path := getPathByRawUrl(rawUrl)

	for _, router := range conf.ConfigInfo.GateWay.Routers {
		if path == router.Path {
			host, err := nacos.GetHttpLbHost(router.ServerName)
			if err != nil {
				return "", err
			}
			return "http://" + host, nil
		}
	}
	return "", errors.New("未找到路径")
}

func GetModuleMiddlewareFilterName(rawUrl *url.URL) []string {

	path := getPathByRawUrl(rawUrl)

	for _, router := range conf.ConfigInfo.GateWay.Routers {
		if path == router.Path {
			return router.Filter
		}
	}
	return nil
}
