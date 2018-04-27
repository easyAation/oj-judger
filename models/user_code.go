package models

import (
	. "github.com/open-fightcoder/oj-judger/common/store"
)

type UserCode struct {
	Id        int64  `form:"id" json:"id"`
	ProblemId int64  `form:"problem_id" json:"problem_id"` //题目ID
	UserId    int64  `form:"user_id" json:"user_id"`       //用户ID
	SaveCode  string `form:"saveCode" json:"saveCode"`     //保存代码
	Language  string `form:"language" json:"language"`     //代码语言
}

func UserCodeCreate(userCode *UserCode) (int64, error) {
	return OrmWeb.Insert(userCode)
}

func UserCodeRemove(id int64) error {
	_, err := OrmWeb.Id(id).Delete(&UserCode{})
	return err
}

func UserCodeUpdate(userCode *UserCode) error {
	_, err := OrmWeb.AllCols().ID(userCode.Id).Update(userCode)
	return err
}

func UserCodeGetById(id int64) (*UserCode, error) {
	userCode := new(UserCode)
	has, err := OrmWeb.Id(id).Get(userCode)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return userCode, nil
}

func UserCodeGetUserCode(userId, problemId int64) (*UserCode, error) {
	userCode := new(UserCode)
	has, err := OrmWeb.Where("user_id=?", userId).And("problem_id=?", problemId).Get(userCode)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return userCode, nil
}
