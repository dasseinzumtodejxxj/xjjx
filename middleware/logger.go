package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	filePath := "log/log.log"
	linkName := "latest_log.log"
	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel) //日志级别
	//日志分割
	logWriter, _ := rotalog.New(
		filePath+"%Y%m%d.log",
		rotalog.WithMaxAge(7*24*time.Hour), //
		rotalog.WithRotationTime(24*time.Hour),
		rotalog.WithLinkName(linkName))
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(Hook)
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds())/1000000.0)))
		hostName, err := os.Hostname() //客户端请求过来的名字
		if err != nil {
			hostName = "unknown"
		}
		statusCode := c.Writer.Status()    //状态码
		clientIp := c.ClientIP()           //IP地址
		userAgent := c.Request.UserAgent() //客户端用什么浏览器请求
		dataSize := c.Writer.Size()        //数据大小
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method   //请求方法
		path := c.Request.RequestURI //请求路径
		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
