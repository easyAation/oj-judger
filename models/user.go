package models

import (
	. "github.com/open-fightcoder/oj-judger/common/store"
)

type User struct {
	Id           int64  `form:"id" json:"id"`
	AccountId    int64  `form:"accountId" json:"accountId"`       //账号Id
	UserName     string `form:"userName" json:"userName"`         //用户名
	NickName     string `form:"nickName" json:"nickName"`         //昵称
	Sex          string `form:"sex" json:"sex"`                   //性别
	Avator       string `form:"avator" json:"avator"`             //头像
	Blog         string `form:"blog" json:"blog"`                 //博客地址
	Git          string `form:"git" json:"git"`                   //Git地址
	Description  string `form:"description" json:"description"`   //个人描述
	Birthday     string `form:"birthday" json:"birthday"`         //生日
	DailyAddress string `form:"dailyAddress" json:"dailyAddress"` //日常所在地：省、市
	StatSchool   string `form:"statSchool" json:"statSchool"`     //当前就学状态(小学及以下、中学学生、大学学生、非在校生)
	SchoolName   string `form:"schoolName" json:"schoolName"`     //学校名称
}

//增加
func (this User) Create(user *User) (int64, error) {
	_, err := OrmWeb.Insert(user) //第一个参数为影响的行数
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

//删除
func (this User) Remove(id int64) error {
	user := User{}
	_, err := OrmWeb.Id(id).Delete(user)
	return err
}

//修改
func (this User) Update(user *User) error {
	_, err := OrmWeb.AllCols().ID(user.Id).Update(user)
	return err
}

//查询
func (this User) GetById(id int64) (*User, error) {
	user := new(User)
	has, err := OrmWeb.Id(id).Get(user) //第一个为 bool 类型，表示是否查找到记录

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

func (this User) GetByAccountId(accountId int64) (*User, error) {
	user := new(User)
	has, err := OrmWeb.Where("account_id = ?", accountId).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

func (this User) QueryByName(nickname string) ([]*User, error) {
	userList := make([]*User, 0)
	err := OrmWeb.Where("nick_name like ?", "%"+nickname+"%").Find(&userList)
	if err != nil {
		return nil, err
	}
	return userList, nil
}
