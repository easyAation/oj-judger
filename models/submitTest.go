package models

import . "github.com/open-fightcoder/oj-judger/common/store"

type SubmitTest struct {
	Id            int64  `form:"id" json:"id"`
	UserId        int64  `form:"userId" json:"userId"`
	Language      string `form:"language" json:"language"`
	SubmitTime    int64  `form:"submitTime" json:"submitTime"`
	RunningTime   int    `form:"runningTime" json:"runningTime"`
	RunningMemory int    `form:"runningMemory" json:"runningMemory"`
	Result        int    `form:"result" json:"result"`
	Input         string `form:"input" json:"input"`
	ResultDes     string `form:"resultDes" json:"resultDes"`
	Code          string `form:"code" json:"code"`
}

func SubmitTestCreate(submitTest *SubmitTest) (int64, error) {
	_, err := OrmWeb.Insert(submitTest) //第一个参数为影响的行数
	if err != nil {
		return 0, err
	}
	return submitTest.Id, nil
}

func SubmitTestRemove(id int64) error {
	_, err := OrmWeb.Id(id).Delete(&SubmitTest{})
	return err
}

func SubmitTestUpdate(submitTest *SubmitTest) error {
	_, err := OrmWeb.AllCols().ID(submitTest.Id).Update(submitTest)
	return err
}

func SubmitTestGetById(id int64) (*SubmitTest, error) {
	submitTest := new(SubmitTest)
	has, err := OrmWeb.Id(id).Get(submitTest)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return submitTest, nil
}
func SubmitTestGetByUserId(userId int64, currentPage int, perPage int) ([]*SubmitTest, error) {
	submitTestList := make([]*SubmitTest, 0)
	err := OrmWeb.Where("user_id=?", userId).Limit(perPage, (currentPage-1)*perPage).Find(&submitTestList)
	if err != nil {
		return nil, err
	}
	return submitTestList, nil
}
