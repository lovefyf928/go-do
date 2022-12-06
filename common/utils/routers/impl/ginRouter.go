package impl

import (
	"github.com/gin-gonic/gin"
	"go-do/common/enum"
	"go-do/common/utils/routers"
	"strings"
)

type GinRouter struct {
	methodMapping []map[enum.HttpMethod]map[string]func(*gin.Context)
	engine        *gin.Engine
	prefix        string
	groupList     map[int]string
	count         int
	urlList       map[string]bool
}

func (r *GinRouter) NewRouter(engine *gin.Engine, prefix string) routers.IRouter {
	ginRouter := &GinRouter{engine: engine, methodMapping: []map[enum.HttpMethod]map[string]func(*gin.Context){}, prefix: prefix, groupList: map[int]string{}, count: 0, urlList: map[string]bool{}}
	ginRouter.methodMapping = append(ginRouter.methodMapping, nil)
	return ginRouter
}

func (r *GinRouter) AddRouter(method enum.HttpMethod, path string, fn func(*gin.Context)) {
	if r.urlList[strings.Split(path, ":")[0]] {
		panic("url重复")
	}
	var methodMap map[enum.HttpMethod]map[string]func(*gin.Context)
	methodMap = r.methodMapping[len(r.methodMapping)-1]

	if methodMap == nil {
		methodMap = map[enum.HttpMethod]map[string]func(*gin.Context){}
	}
	if methodMap[method] == nil {
		methodMap[method] = map[string]func(*gin.Context){}
	}
	methodMap[method][path] = fn
	r.methodMapping[len(r.methodMapping)-1] = methodMap
	r.urlList[strings.Split(path, ":")[0]] = true
	r.count++
}

func (r *GinRouter) AddGroup(groupPrefix string) {
	r.groupList[r.count+1] = groupPrefix
	r.urlList = map[string]bool{}
	r.methodMapping = append(r.methodMapping, nil)
}

func (r *GinRouter) Register() {
	count := 0
	var group *gin.RouterGroup
	for i := range r.methodMapping {
		for httpMethod := range r.methodMapping[i] {
			for path := range r.methodMapping[i][httpMethod] {
				count++
				if groupPrefix := r.groupList[count]; groupPrefix != "" {
					group = r.engine.Group(r.prefix + groupPrefix)
				}
				if group == nil {
					r.engine.Handle(string(httpMethod), r.prefix+path, r.methodMapping[i][httpMethod][path])
				} else {
					group.Handle(string(httpMethod), path, r.methodMapping[i][httpMethod][path])
				}
			}
		}
	}
}
