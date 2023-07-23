package model

import (
	"encoding/base64"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	UserName string `gorm:"type: varchar(20); not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	PassWord string `gorm:"type: varchar(20); not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role     int    `gorm:"type : int;DEFAULT:1" json:"role" validate:"required,gte=1" lable:"角色码"` //2 管理员 ， 3用户
}

// 查询用户是否存在
func CheckUser(username string) (code int) {
	var users User
	db.Select("id").Where("user_name = ?", username).First(&users)
	if users.ID > 0 {
		return errmsg.ERROR_USERNAME_USER //1001 用户名已存在
	}
	return errmsg.SUCCSE
}

// CheckUpUser 更新查询
func CheckUpUser(id int, name string) (code int) {
	var user User
	db.Select("id, user_name").Where("user_name = ?", name).First(&user)
	if user.ID == uint(id) {
		return errmsg.SUCCSE
	}
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USER //1001
	}
	return errmsg.SUCCSE
}

// 新增用户
func CreateUser(data *User) int {
	data.PassWord = ScryptPw(data.PassWord)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE // 200
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Limit(1).Where("ID = ?", id).Find(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

// GetUsers 查询用户列表
func GetUsers(username string, pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64

	if username != "" {
		db.Select("id,user_name,role,created_at").Where(
			"user_name LIKE ?", username+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		db.Model(&users).Where(
			"user_name LIKE ?", username+"%",
		).Count(&total)
		return users, total
	}
	db.Select("id,user_name,role,created_at").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	db.Model(&users).Count(&total)

	if err != nil {
		return users, 0
	}
	return users, total
}

// 编辑用户信息（密码以外）
func EdirUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["user_name"] = data.UserName
	maps["role"] = data.Role
	err := db.Model(&user).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// ChangePassword 修改密码
func ChangePassword(id int, data *User) int {
	var user User
	//var maps = make(map[string]interface{})
	//maps["password"] = data.Password

	err = db.Model(&user).Select("pass_word").Where("id = ?", id).Updates(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// 删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id = ? ", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// BeforeCreate 密码加密&权限控制
/*func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.PassWord = ScryptPw(u.PassWord)
	u.Role = 2
	return nil
}*/

/*func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.PassWord = ScryptPw(u.PassWord)
	return nil
}*/

// ScryptPw 生成密码
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 32, 4, 6, 66, 22, 222, 11}
	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(errmsg.ERROR_PASSWORD_WRONG)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}

// CheckLogin 后台登录验证
func CheckLogin(username string, password string) int {
	var user User
	db.Where("user_name = ?", username).First(&user)

	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.PassWord {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role != 1 {
		return errmsg.ERROR_USER_NO_RIGHT
	}
	return errmsg.SUCCSE
}

// CheckLoginFront 前台登录
func CheckLoginFront(username string, password string) int {
	var user User

	db.Where("user_name = ?", username).First(&user)

	if user.ID == 0 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPw(password) != user.PassWord {
		return errmsg.ERROR_PASSWORD_WRONG
	}
	return errmsg.SUCCSE
}
