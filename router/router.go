package router

import (
	"bytes"
	"context"
	"errors"
	"github.com/bingoohuang/golog/pkg/ginlogrus"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"libvirt-manager/ctler"
	"libvirt-manager/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func PostPreProcess(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := c.GetRawData()
		if body == nil {
			handler(c)
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 必须写回

		// 通过url地址判断方式
		//if !strings.Contains(c.Request.Host, data.IP) {
		//	ctler.PublicProxy(c, data.IP)
		//	return
		//}

		ip, _ := jsonparser.GetString(body, "ip")
		// 查找本地IP方式
		ipList, _ := utils.GetLocalIPList()
		if !utils.Contains(ipList, ip) {
			ctler.PublicProxy(c, ip)
			return
		}

		// 执行本地函数
		handler(c)
		c.Next()
	}
}

func StartServer() {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(ginlogrus.Logger(nil, true), gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, nil)
	})

	/***********  路由配置  ************/
	r.POST("/test", ctler.Test) // 测试用可删除
	r.POST("/createvm", PostPreProcess(ctler.CreateVm))
	r.POST("/getvmlist", PostPreProcess(ctler.GetVmList))
	r.POST("/getvmstate", PostPreProcess(ctler.GetVMState))
	r.POST("/vmstart", PostPreProcess(ctler.VMStart))
	r.POST("/vmstop", PostPreProcess(ctler.VMStop))
	r.POST("/vmpause", PostPreProcess(ctler.VMPause))
	r.POST("/vmresume", PostPreProcess(ctler.VMResume))
	r.POST("/vmrestart", PostPreProcess(ctler.VMRestart))
	r.POST("/vmdelete", PostPreProcess(ctler.VMDelete))

	srv := &http.Server{
		Addr:    "0.0.0.0:" + utils.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("port is used:%v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	logrus.Warn("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Warn("Exiting")
}
