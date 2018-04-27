package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-fightcoder/oj-judger/managers"
	"github.com/open-fightcoder/oj-judger/router/controllers/base"
)

func RegisterJudge(router *gin.RouterGroup) {
	router.POST("judge/test", httpHandlerJudgeTest)
	router.POST("judge/special", httpHandlerJudgeSpecial)
	router.POST("judge/default", httpHandlerJudgeDefault)
}

type JudgeParam struct {
	SubmitId int64 `form:"submit_id" json:"submit_id"`
}

func httpHandlerJudgeTest(c *gin.Context) {
	job := JudgeParam{}
	err := c.Bind(&job)
	if err != nil {
		panic(err)
	}

	err = managers.JudgeTest(job.SubmitId)
	if err != nil {
		c.JSON(http.StatusOK, base.Fail(err.Error()))
	}
	c.JSON(http.StatusOK, base.Success("ok"))
}

func httpHandlerJudgeSpecial(c *gin.Context) {
	job := JudgeParam{}
	err := c.Bind(&job)
	if err != nil {
		panic(err)
	}

	err = managers.JudgeSpecial(job.SubmitId)
	if err != nil {
		c.JSON(http.StatusOK, base.Fail(err.Error()))
	}
	c.JSON(http.StatusOK, base.Success("ok"))
}

func httpHandlerJudgeDefault(c *gin.Context) {
	job := JudgeParam{}
	err := c.Bind(&job)
	if err != nil {
		panic(err)
	}
	err = managers.JudgeDefault(job.SubmitId)
	if err != nil {
		c.JSON(http.StatusOK, base.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, base.Success("ok"))
}
