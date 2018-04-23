package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-fightcoder/oj-judger/router/controllers/base"
)

func RegisterJudge (router *gin.RouterGroup) {
	router.POST("judge/test", httpHandlerJudgeTest)
	router.POST("judge/special", httpHandlerJudgeSpecial)
	router.POST("judge/default", httpHandlerJudgeDefault)
}

type JudgeParam struct {
	SubmitId    int64 `form:"submit_id" json:"submit_id"`
}

func httpHandlerJudgeTest(c *gin.Context) {
	job := JudgeParam{}
	err := c.Bind(&job)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, base.Success())
}

func httpHandlerJudgeSpecial(c *gin.Context) {
	job := JudgeParam{}
	err := c.Bind(&job)
	if err != nil {
		panic(err)
	}
}

func httpHandlerJudgeDefault(c *gin.Context) {
	job := JudgeParam{}
	err := c.Bind(&job)
	if err != nil {
		panic(err)
	}
}