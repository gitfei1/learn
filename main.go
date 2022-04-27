package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"learn/util"
	"net/http"
	"strings"
	"time"
)

var Auth string
var User string
var Pass string
var DeviceList []string
func main() {
	DeviceList = make([]string,0)
	User = "admin"
	Pass = "123456"
	r := gin.Default()
	r.POST("/login",Login)
	r.POST("/c-interface/alarm-report/device-blacklist/set",Set)
	r.POST("/c-interface/alarm-report/device-blacklist/get",Get)
	err := r.Run("0.0.0.0:10010") // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println("run err",err)
		return
	}
}
func Login(c *gin.Context)  {
	//登录主要做两件事 一个校验用户名密码 一个生成token
	var req util.LoginRequest
	if err := c.BindJSON(&req);err != nil{
		fmt.Println("login bind json err",err,"  ",req.UserName,req.Password)
		c.JSON(http.StatusOK,util.Response{
			Code: util.PARAMSERR,
			Message: "请求参数有误",
		})
	}
	//密码是md5之后的值
	md5result := fmt.Sprintf("%x",md5.Sum([]byte(Pass)))
	if req.Password == md5result && req.UserName == User  {
		fmt.Println("user and pass correct !")
		var response util.LoginResponse
		Auth = GetToken()
		response.AccessToken = Auth
		response.ExpiresIn = 600
		fmt.Println("token:",Auth)
		go func() {
			time.Sleep(time.Second*time.Duration(response.ExpiresIn))
			fmt.Println("time is end ")
			Auth = ""
		}()
		c.JSON(http.StatusOK,util.Response{
			Code: util.SUCCESS,
			Message: "ok",
			Data: response,
		})
	}else {
		fmt.Println("err user ",req.UserName, req.Password,User,md5result)
		c.JSON(http.StatusOK,util.Response{
			Code: util.PARAMSERR,
			Message: "请求参数有误",
		})
	}

}

func Set(c *gin.Context)  {
	//设置黑名单主要做一件事 将收到的黑名单设备id保存起来 全局的(考虑保存到数据库)
	if Auth == "" {
		fmt.Println("need login ")
		c.JSON(http.StatusOK,util.Response{
			Code: util.NEEDLOGIN,
			Message: "需要登录",
		})
	}else if Auth == c.Request.Header.Get("Authorization"){
		fmt.Println("auth success ")
		var req util.SetRequest
		if err := c.BindJSON(&req);err!= nil {
			fmt.Println("set bind json err",err)
			c.JSON(http.StatusOK,util.Response{
				Code: util.PARAMSERR,
				Message: "请求参数有误",
			})
		}else {
			DeviceList = append(DeviceList,req.GidBlacklist...)
			fmt.Println("the length of the device ",len(DeviceList))
			c.JSON(http.StatusOK,util.Response{
				Code: util.SUCCESS,
				Message: "ok",
			})
		}
	}else {
		fmt.Println("auth fail ")
		c.JSON(http.StatusOK,util.Response{
			Code: util.NOAUTH,
			Message: "没有访问权限",
		})
	}

}

func Get(c *gin.Context)  {
	// 获取黑名单主要做一件事 查询数据库返回黑名单
	tempAuth := strings.Split(c.Request.Header.Get("Authorization")," ")
	fmt.Println("auth:",tempAuth)
	if Auth == "" || len(tempAuth) < 2{
		fmt.Println("need login ")
		c.JSON(http.StatusOK,util.Response{
			Code: util.NEEDLOGIN,
			Message: "需要登录",
		})
	}else if Auth == tempAuth[1]{
		fmt.Println("auth success ")
		var req util.GetRequest
		if err := c.BindJSON(&req);err!= nil {
			fmt.Println("get bind json err",err)
			c.JSON(http.StatusOK,util.Response{
				Code: util.PARAMSERR,
				Message: "请求参数有误",
			})
		}else {
			fmt.Println("set length of the device ",req.Limit,req.Offset)
			var getresponse util.GetResponse
			if req.Limit == 0 {
				result := DeviceList[req.Offset:]
				getresponse.Total = uint(len(result))
				getresponse.GidBlacklist = result
			}else {
				result := DeviceList[req.Offset:req.Offset+req.Limit]
				getresponse.Total = uint(len(result))
				getresponse.GidBlacklist = result
			}
			c.JSON(http.StatusOK,util.Response{
				Code: util.SUCCESS,
				Message: "ok",
				Data: getresponse,
			})
		}
	}else {
		fmt.Println("auth fail ",tempAuth)
		c.JSON(http.StatusOK,util.Response{
			Code: util.NOAUTH,
			Message: "没有访问权限",
		})
	}
}

func GetToken()  string{
	//做一件事 生成token
	t := time.Now().Unix()
	h := md5.Sum([]byte(fmt.Sprintf("%#v",t)))
	return fmt.Sprintf("%x",h)
}