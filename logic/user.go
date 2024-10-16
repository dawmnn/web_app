package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

// 存放业务逻辑代码
func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户不存在
	//fmt.Println("cccccccccccc")
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//fmt.Println("bbbbbbbbb")
	//2.生成UID
	userID := snowflake.GenID()
	user := &models.User{
		Username: p.Username,
		UserID:   userID,
		Password: p.Password,
	}
	//构造一个User实例
	//3.保存进数据库
	err = mysql.InsertUser(user)
	return err
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就能拿到user.UserID
	if err = mysql.Login(user); err != nil {
		return
	}
	//生成JWT
	//user.UserID
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
