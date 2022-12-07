package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func Proxy(c *gin.Context) {

	err := setProxyPath(c.Request.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("填写的地址有误: %s", err.Error()))
		c.Abort()
		return
	}

	req, err := http.NewRequestWithContext(c, c.Request.Method, c.Request.URL.String(), c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	defer req.Body.Close()
	req.Header = c.Request.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	for k := range resp.Header {
		for j := range resp.Header[k] {
			c.Header(k, resp.Header[k][j])
		}
	}
	extraHeaders := make(map[string]string)
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, extraHeaders)
	c.Abort()
}

func setProxyPath(rawUrl *url.URL) error {

	proxyUrl, err := getModuleProxyPath(rawUrl)

	if err != nil {
		return err
	}

	if proxyUrl == "" {
		return errors.New("服务不可用")
	}

	u, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}

	rawUrl.Scheme = u.Scheme
	rawUrl.Host = u.Host
	ruq := rawUrl.Query()
	rawUrl.RawQuery = ruq.Encode()
	return nil
}
