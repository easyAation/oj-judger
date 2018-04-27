package common

import (
	"github.com/open-fightcoder/oj-judger/common/g"
	"github.com/open-fightcoder/oj-judger/common/store"
)

func Init(cfgFile string) {
	g.LoadConfig(cfgFile)
	g.InitLog()
	store.InitMysql()
	//store.InitRedis()
	store.InitMinio()
}

func Close() {
	store.CloseMysql()
	//store.CloseRedis()
}
