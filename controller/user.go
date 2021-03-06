package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/modules"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)


//SignUpHandler  处理注册请求的函数
func SignUpHandler(c *gin.Context)  {
	//1、获取参数和参数校验
	p := new(modules.ParamSignUp)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		zap.L().Error("SignUp with invalid param",zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs,ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c,CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2、业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed",zap.Error(err))
		if errors.Is(err,mysql.ErrorUserExist){
			ResponseError(c,CodeUserExist)
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3、返回响应
	ResponseSuccess(c, nil)
	return
}

func LoginHandler(c *gin.Context)  {
	// 1、获取请求参数及参数校验

	p := new(modules.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		zap.L().Error("Login with invalid param",zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs,ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c,CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c,CodeInvalidParam,removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2、业务逻辑处理
	token,err := logic.Login(p);
	if err != nil {
		zap.L().Error("logic.Login failed",zap.String("username",p.Username),zap.Error(err))
		if errors.Is(err,mysql.ErrorUserNotExist){
			ResponseError(c,CodeUserNotExist)
		}
		ResponseError(c,CodeInvalidPassword)
		return
	}

	// 3、返回响应
	ResponseSuccess(c,token)
	return
}
