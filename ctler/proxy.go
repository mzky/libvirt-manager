package ctler

import (
	"github.com/gin-gonic/gin"
	"libvirt-manager/utils"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	ProxyTimeout = 120 * time.Minute
)

func PublicProxy(ctx *gin.Context, ip string) {
	u := &url.URL{
		Scheme: "http",
		Host:   ip + ":" + utils.Port,
	}
	remoteProxy := httputil.NewSingleHostReverseProxy(u)
	remoteProxy.ServeHTTP(ctx.Writer, ctx.Request)

	//var simpleHostProxy = httputil.ReverseProxy{
	//	Director: func(req *http.Request) {
	//		req.URL.Scheme = "http"
	//		req.URL.Host = u.Host
	//		req.Host = u.Host
	//
	//		req.Header.Add("X-Forwarded-For", ctx.ClientIP()) // 客户端IP
	//		req.Header.Add("X-Forwarded-Host", req.Host)
	//		req.Header.Add("X-Origin-Host", req.Header.Get("Host"))
	//		// fmt.Println(req.URL.Scheme, req.Host, req.URL.Path)
	//	},
	//	Transport: &http.Transport{
	//		DialContext:       (&net.Dialer{Timeout: ProxyTimeout}).DialContext,
	//		DisableKeepAlives: true,                                  // 短连接
	//		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, // true 跳过证书验证
	//	},
	//	ModifyResponse: func(r *http.Response) error {
	//		respLocationHeader := r.Header.Get("Location")
	//		if gonet.IsRelativeForward(r.StatusCode, respLocationHeader) {
	//			// 301/302时，本地相对路径跳转时，改写Location返回头
	//			basePath := strings.TrimRight(ctx.Request.URL.Path, u.Host)
	//			r.Header.Set("Location", basePath+respLocationHeader)
	//		}
	//		return nil
	//	},
	//}
	//simpleHostProxy.ServeHTTP(ctx.Writer, ctx.Request)
}
