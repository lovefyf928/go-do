package chain

import (
	"github.com/gin-gonic/gin"
	"go-do/common/authorization"
	"go-do/common/utils"
	"go-do/middleware"
)

type MiddlewareFn func(ctx *Chain)

type Chain struct {
	fnChain []map[string]MiddlewareFn

	index int

	permitted chan int

	denied chan int

	MiddlewareNames []string

	Token string

	Url string
}

func (c *Chain) Next() {
	if len(c.fnChain)-1 == c.index {
		close(c.permitted)
		return
	}
	c.index++
	fnMap := c.fnChain[c.index]
	for s := range fnMap {
		if num := utils.StringsContains(c.MiddlewareNames, s); c.MiddlewareNames == nil && num == -1 {
			c.Next()
			return
		}
		fnMap[s](c)
	}
}

func (c *Chain) Denied() {
	close(c.denied)
}

func (c *Chain) Apply(middlewareName string, fn MiddlewareFn) {
	fnMap := map[string]MiddlewareFn{}
	fnMap[middlewareName] = fn
	c.fnChain = append(c.fnChain, fnMap)
}

func NewChain(middlewareNames []string, token string, url string) *Chain {
	return &Chain{fnChain: []map[string]MiddlewareFn{}, index: -1, permitted: make(chan int), denied: make(chan int), MiddlewareNames: middlewareNames, Token: token, Url: url}
}

func ChainMiddleware(fn func(chain *Chain)) func(*gin.Context) {
	return func(c *gin.Context) {
		chain := NewChain(middleware.GetModuleMiddlewareFilterName(c.Request.URL), c.Request.Header.Get(authorization.TOKEN_HEADER_NAME), c.Request.URL.Path)
		fn(chain)
		select {
		case <-chain.permitted:
			c.Next()
		case <-chain.denied:
			c.Abort()
		}
	}
}
