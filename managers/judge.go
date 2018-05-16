package managers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"encoding/json"

	"github.com/mholt/archiver"
	"github.com/minio/minio-go"
	"github.com/open-fightcoder/oj-judger/common/g"
	"github.com/open-fightcoder/oj-judger/common/store"
	"github.com/open-fightcoder/oj-judger/judge"
	"github.com/open-fightcoder/oj-judger/models"
	"github.com/open-fightcoder/oj-judger/redis"
	log "github.com/sirupsen/logrus"
)

func JudgeTest(submitId int64) judge.Result {
	// 获取提交信息：代码，语言，用户输入
	// 编译
	// 运行
	// 写入结果
	log.Debugf("submit:%d start judge test", submitId)
	submit, err := models.SubmitTestGetById(submitId)
	if err != nil {
		err = fmt.Errorf("get submit %d failure: %s", submitId, err.Error())
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	if submit == nil {
		err = fmt.Errorf("get submit %d failure: col not found", submitId)
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	log.Debugf("submit:%d create workdir", submitId)
	workDir, err := createWorkDir("test", submitId, submit.UserId)
	if err != nil {
		err = fmt.Errorf("create workDir %s failure: %s", workDir, err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callTestResult(submit, res)
		return res
	}
	log.Debugf("submit:%d get code", submitId)
	err = getCode(submit.Code, workDir)
	if err != nil {
		err = fmt.Errorf("get code file %s failure: %s", submit.Code, err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callTestResult(submit, res)
		return res
	}
	err = ioutil.WriteFile(workDir+"/input", []byte(submit.Input), 0644)
	if err != nil {
		err = fmt.Errorf("get input %s failure: %s", submit.Input, err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callTestResult(submit, res)
		return res
	}
	callTestResult(submit, judge.Result{
		ResultCode: judge.Compiling,
	})
	log.Debugf("submit:%d start compile", submitId)
	j := judge.NewJudge(submit.Language)
	result := j.Compile(workDir, submit.Code)
	if result.ResultCode != 0 {
		// 编译失败
		callTestResult(submit, result)
		return result
	}
	log.Infof("%d start run", submitId)
	// 运行中
	callTestResult(submit, judge.Result{
		ResultCode: judge.Running,
	})
	result = j.Run(workDir+"/user.bin",
		workDir+"/input",
		workDir+"/output",
		int64(5000),
		int64(100000))
	log.Infof("%d run result %#v", submitId, result)
	if result.ResultCode != judge.Normal {
		callTestResult(submit, result)
		return result
	}
	data, err := ioutil.ReadFile(workDir + "/output")
	if err != nil {
		err = fmt.Errorf("get output failure: %s", err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callTestResult(submit, res)
		return res
	}
	result.ResultCode = judge.Accepted
	result.ResultDes = string(data)

	callTestResult(submit, result)
	return result
}

func JudgeSpecial(submitId int64) judge.Result {
	// 获取提交信息
	// 编译
	// 运行
	// 执行标准输入运行，得到标准输出
	// 获取题目信息
	// 编译特判断
	// 将标准输出作为特判程序输入
	// 拿到判断结果
	// 写入结果

	return judge.Result{}
}

func JudgeDefault(submitId int64) judge.Result {
	log.Debugf("submit:%d start judge default", submitId)
	submit, err := models.SubmitGetById(submitId)
	if err != nil {
		err = fmt.Errorf("get submit %d failure: %s", submitId, err.Error())
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	if submit == nil {
		err = fmt.Errorf("get submit %d failure: col not found", submitId)
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	log.Debugf("submit:%d create workdir", submitId)
	workDir, err := createWorkDir("default", submitId, submit.UserId)
	if err != nil {
		err = fmt.Errorf("create workDir %s failure: %s", workDir, err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callDefaResult(submit, res)
		return res
	}
	log.Debugf("submit:%d get problem", submitId)
	problem, err := models.ProblemGetById(submit.ProblemId)
	if err != nil {
		err = fmt.Errorf("get problem failure: %s", err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callDefaResult(submit, res)
		return res
	}
	if problem == nil {
		err = fmt.Errorf("get problem %d failure: col not found", submit.ProblemId)
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callDefaResult(submit, res)
		return res
	}
	log.Debugf("submit:%d get code", submitId)
	err = getCode(submit.Code, workDir)
	if err != nil {
		err = fmt.Errorf("get code file %s failure: %s", submit.Code, err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callDefaResult(submit, res)
		return res
	}
	log.Debugf("submit:%d get case", submitId)
	err = getCase(problem.CaseData, workDir)
	if err != nil {
		err = fmt.Errorf("get case file %s failure: %s", problem.CaseData, err.Error())
		res := judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
		callDefaResult(submit, res)
		return res
	}

	callDefaResult(submit, judge.Result{
		ResultCode: judge.Compiling,
	})
	log.Debugf("submit:%d start compile", submitId)
	j := judge.NewJudge(submit.Language)
	result := j.Compile(workDir, submit.Code)
	if result.ResultCode != 0 {
		// 编译失败
		callDefaResult(submit, result)
		return result
	}
	log.Infof("%d start run", submitId)
	// 运行中
	callDefaResult(submit, judge.Result{
		ResultCode: judge.Running,
	})

	totalResult := judge.Result{
		ResultCode:    judge.Accepted,
		ResultDes:     "",
		RunningMemory: 0,
		RunningTime:   0,
	}

	caseList := getCaseList(workDir + "/case")
	log.Infof("%d case list %#v", submitId, caseList)
	for _, name := range caseList {
		result = j.Run(workDir+"/user.bin",
			workDir+"/case/"+name+".in",
			workDir+"/"+name+".user",
			int64(problem.TimeLimit),
			int64(problem.MemoryLimit))
		log.Infof("%d run result %#v", submitId, result)
		if result.ResultCode != judge.Normal {
			totalResult = result
			break
		}

		if result.RunningMemory > totalResult.RunningMemory {
			totalResult.RunningMemory = result.RunningMemory
		}

		if result.RunningTime > totalResult.RunningTime {
			totalResult.RunningTime = result.RunningTime
		}

		diff := compare(workDir+"/"+name+".user", workDir+"/case/"+name+".out")
		if diff != "" {
			result.ResultCode = judge.WrongAnswer
			result.ResultDes = diff
			totalResult = result
			break
		}
	}

	callDefaResult(submit, totalResult)
	return totalResult
}

type SubmitCount struct {
	Accepted            int64 `json:"accepted"`
	FailNum             int64 `json:"fail_num"`
	WrongAnswer         int64 `json:"wrong_answer"`
	CompilationError    int64 `json:"compilation_error"`
	TimeLimitExceeded   int64 `json:"time_limit_exceeded"`
	MemoryLimitExceeded int64 `json:"memory_limit_exceeded"`
	OutputLimitExceeded int64 `json:"output_limit_exceeded"`
	RuntimeError        int64 `json:"runtime_error"`
	SystemError         int64 `json:"system_error"`
}

type ProblemCount struct {
	AcNum    int64 `json:"ac_num"`
	TotalNum int64 `json:"total_num"`
}

func callTestResult(submit *models.SubmitTest, result judge.Result) {
	submit.Result = result.ResultCode
	submit.ResultDes = result.ResultDes
	submit.RunningTime = result.RunningTime
	submit.RunningMemory = result.RunningMemory

	length := len(submit.ResultDes)
	if length > 999 {
		length = 999
	}
	submit.ResultDes = string([]byte(submit.ResultDes)[:length])

	log.Infof("%d call test result %#v", submit.Id, result)
	err := models.SubmitTestUpdate(submit)
	if err != nil {
		log.Error("call test result failure:", err.Error())
	}
	return
}

func isAc(submit *models.Submit) bool {
	submitList, _ := models.SubmitGetByUserId(submit.UserId)
	for _, sub := range submitList {
		if sub.Result == judge.Accepted {
			return true
		}
	}
	return false
}

func callDefaResult(submit *models.Submit, result judge.Result) {
	submit.Result = result.ResultCode
	submit.ResultDes = result.ResultDes
	submit.RunningTime = result.RunningTime
	submit.RunningMemory = result.RunningMemory

	length := len(submit.ResultDes)
	if length > 999 {
		length = 999
	}
	submit.ResultDes = string([]byte(submit.ResultDes)[:length])

	log.Infof("%d call defalut result %#v", submit.Id, result)
	err := models.SubmitUpdate(submit)
	if err != nil {
		log.Error("call defalut result failure:", err.Error())
	}
	if submit.Result == judge.Accepted {
		if !isAc(submit) {
			redis.PersonWeekRankUpdate(1, submit.UserId)
			redis.PersonMonthRankUpdate(1, submit.UserId)
			redis.RankListUpdate(1, submit.UserId)
			redis.ProblemStatusSet(submit.UserId, submit.ProblemId, 1)
		}
	}

	if submit.Result > judge.Running {
		redis.ProblemStatusSet(submit.UserId, submit.ProblemId, 2)

		jsonStr, err := redis.ProblemCountGet(submit.ProblemId)
		if err != nil {
			log.Error("problemCount get failure:", err.Error())
			return
		}
		var problemCount ProblemCount
		err = json.Unmarshal([]byte(jsonStr), &problemCount)
		if err != nil {
			log.Error("problemCount get failure:", err.Error())
			return
		}
		problemCount.TotalNum += 1
		if submit.Result == judge.Accepted {
			problemCount.AcNum += 1
		}
		data, err := json.Marshal(problemCount)
		if err != nil {
			log.Error("problemCount set failure:", err.Error())
			return
		}
		flag := redis.ProblemCountSet(submit.ProblemId, string(data))
		if !flag {
			log.Error("problemCount set failure")
		}

		jsonStr, err = redis.SubmitCountGet(submit.UserId)
		if err != nil {
			log.Error("submitcount get failure:", err.Error())
			return
		}
		var submitCount SubmitCount
		err = json.Unmarshal([]byte(jsonStr), &submitCount)
		if err != nil {
			log.Error("submitcount get failure:", err.Error())
			return
		}
		switch submit.Result {
		case judge.Accepted:
			submitCount.Accepted += 1
		case judge.WrongAnswer:
			submitCount.WrongAnswer += 1
		case judge.CompilationError:
			submitCount.CompilationError += 1
		case judge.TimeLimitExceeded:
			submitCount.TimeLimitExceeded += 1
		case judge.MemoryLimitExceeded:
			submitCount.MemoryLimitExceeded += 1
		case judge.OutputLimitExceeded:
			submitCount.OutputLimitExceeded += 1
		case judge.RuntimeError:
			submitCount.RuntimeError += 1
		case judge.SystemError:
			submitCount.SystemError += 1
		}

		data, err = json.Marshal(submitCount)
		if err != nil {
			log.Error("submitcount set failure:", err.Error())
			return
		}
		flag = redis.SubmitCountSet(submit.UserId, string(data))
		if !flag {
			log.Error("submitcount set failure")
		}
	}

	return
}

func getCode(code string, workDir string) error {
	err := store.MinioClient.FGetObject(g.Conf().Minio.CodeBucket,
		code, workDir+"/"+code, minio.GetObjectOptions{})
	return err
}

func getCase(cs string, workDir string) error {
	err := store.MinioClient.FGetObject(g.Conf().Minio.CaseBucket,
		cs, workDir+"/case.zip", minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	err = archiver.Zip.Open(workDir+"/case.zip", workDir+"/case")
	return err
}

func getCaseList(path string) []string {
	dir_list, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error(err)
		return nil
	}

	caseList := make([]string, 0)

	for _, v := range dir_list {
		if v.IsDir() {
			continue
		}
		name := v.Name()
		if name[len(name)-3:] == ".in" {
			caseList = append(caseList, name[:len(name)-3])
		}
	}

	return caseList
}

func compare(userOutput string, caseOutput string) string {
	// -N 将不存在的文件作为空白量比较
	// -B 忽略空行
	// --speed-large-files 假设是大文件、很多分散的小变化
	cmd := exec.Command("diff", "-NB", userOutput, caseOutput)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Info("diff err:", err)
		return string(output)
	}

	return ""
}

func createWorkDir(judgeType string, submitId int64, userId int64) (string, error) {
	dir := fmt.Sprintf("%s/work/%s/%d_%d", getCurrentPath(), judgeType, submitId, userId)
	err := os.MkdirAll(dir, 0777)
	return dir, err
}

func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("getCurrentPath: " + err.Error())
	}
	return dir
}
