package managers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/minio/minio-go"
	"github.com/open-fightcoder/oj-judger/common/g"
	"github.com/open-fightcoder/oj-judger/common/store"
	"github.com/open-fightcoder/oj-judger/judge"
	"github.com/open-fightcoder/oj-judger/models"
)

func JudgeTest(submitId int64) judge.Result {
	// 获取提交信息：代码，语言，用户输入
	// 编译
	// 运行
	// 写入结果
	return judge.Result{}
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
	//fmt.Println("sleep before")
	//time.Sleep(10 * time.Second)
	//fmt.Println("sleep after")
	fmt.Println("start")
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

	fmt.Println("workdir")
	workDir, err := createWorkDir("default", submitId, submit.UserId)
	if err != nil {
		err = fmt.Errorf("create workDir %s failure: %s", workDir, err.Error())
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	problem, err := models.ProblemGetById(submit.ProblemId)
	if err != nil {
		err = fmt.Errorf("get problem failure: %s", err.Error())
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	if problem == nil {
		err = fmt.Errorf("get problem %d failure: col not found", submit.ProblemId)
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}

	fmt.Println("getcode")
	err = getCode(submit.Code, workDir)
	if err != nil {
		err = fmt.Errorf("get code file %s failure: %s", submit.Code, err.Error())
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}
	fmt.Println("getcase")
	err = getCase(problem.CaseData, workDir)
	if err != nil {
		err = fmt.Errorf("get case file %s failure: %s", problem.CaseData, err.Error())
		return judge.Result{
			ResultCode: judge.SystemError,
			ResultDes:  err.Error(),
		}
	}

	// 编译中
	callResult(judge.Result{
		ResultCode: judge.Compiling,
	})

	j := judge.NewJudge(submit.Language)
	result := j.Compile(workDir, submit.Code)
	if result.ResultCode != 0 {
		// 编译失败
		callResult(result)
		return result
	}

	// 运行中
	callResult(result)

	totalResult := judge.Result{
		ResultCode:    judge.Accepted,
		ResultDes:     "",
		RunningMemory: 0,
		RunningTime:   0,
	}

	caseList := getCaseList(workDir + "/case")
	for _, name := range caseList {
		result = j.Run(workDir+"/user.bin",
			workDir+"/case/"+name+".in",
			workDir+"/"+name+".user")
		if result.ResultCode != judge.Normal {
			fmt.Printf("Running Error :%#v\n", result)
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
			totalResult = result
			break
		}
	}

	callResult(totalResult)

	return totalResult
}

func callResult(result judge.Result) {
	//submit := &models.Submit {
	//	Result:result.ResultCode,
	//	ResultDes:result.ResultDes,
	//	RunningTime:result.RunningTime,
	//	RunningMemory:result.RunningMemory,
	//}

	//return models.SubmitUpdate(submit)
	fmt.Printf("%#v\n", result)
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
		panic(err)
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
	cmd := exec.Command("diff", "-B", userOutput, caseOutput)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("diff err:", err)
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
		panic("getCurrentPath: " + err.Error())
	}
	return dir
}
