package logic

import (
	"bluebell/dao/mysql"
	"bluebell/modules"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

//存放业务逻辑的代码

func SignUp(p *modules.ParamSignUp)  (err error) {
	//1、判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//2、生成UID
	UserID := snowflake.GenID()
	//构造一个User实例
	user := &modules.User{
		UserID:   UserID,
		Username: p.Username,
		Password: p.Password,
	}

	//3、保存进数据库
	return mysql.InsertUser(user)
	
}

func Login(p *modules.ParamLogin)  (token string,err error) {
	user := &modules.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递是一个指针，就能拿到user.UserID
	if err :=  mysql.Login(user); err != nil {
		return  "", err
	}
	//生成jwt
	return  jwt.GenToken(user.UserID,user.Username)
}