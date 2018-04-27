package models

import (
	. "github.com/open-fightcoder/oj-judger/common/store"
)

type SysConfig struct {
	Id       int64  `form:"id" json:"id"`
	SysKey   string `form:"sysKey" json:"sysKey"`     //键
	SysValue string `form:"sysValue" json:"sysValue"` //值
}

func SysConfigCreate(sysConfig *SysConfig) (int64, error) {
	return OrmWeb.Insert(sysConfig)
}

func SysConfigRemove(id int64) error {
	_, err := OrmWeb.Id(id).Delete(&SysConfig{})
	return err
}

func SysConfigUpdate(sysConfig *SysConfig) error {
	_, err := OrmWeb.AllCols().ID(sysConfig.Id).Update(sysConfig)
	return err
}

func SysConfigGetById(id int64) (*SysConfig, error) {
	sysConfig := new(SysConfig)
	has, err := OrmWeb.Id(id).Get(sysConfig)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return sysConfig, nil
}
