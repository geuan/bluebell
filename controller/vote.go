package controller

import (
	"bluebell/logic"
	"bluebell/modules"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票
type VoteData struct {
	// UserID 从请求中获取当前的用户
	PostID 	int64 `json:"post_id,string"`  // 帖子id
	Direction 	int `json:"direction,string"`  // 赞成票（1）反对票（-1）
}


func PostVoteController(c *gin.Context)  {
	// 参数校验
	p := new(modules.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)  //类型断言
		if !ok {
			ResponseError(c,CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))  // 翻译并去除掉错误提示中的结构体标识
		ResponseErrorWithMsg(c,CodeInvalidParam,errData)
		return
	}
	//获取当前请求的用户id
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c,CodeNeedLogin)
		return
	}

	//具体投票的业务逻辑
	if err := logic.VoteForPost(userID,p); err != nil {
		zap.L().Error("logic.VoteForPost(userID,p) failed",zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	ResponseSuccess(c,nil)
}