package models

import (
	"math"
	"strings"

	. "github.com/open-fightcoder/oj-judger/common/store"
)

type Problem struct {
	Id                 int64  `form:"id" json:"id"`
	Flag               int    `form:"flag" json:"flag"`                             //1-普通题目 2-用户题目
	Status             int    `form:"status" json:"status"`                         //申请状态
	UserId             int64  `form:"user_id" json:"user_id"`                       //题目提供者
	Difficulty         string `form:"difficulty" json:"difficulty"`                 //题目难度
	CaseData           string `form:"caseData" json:"caseData"`                     //测试数据
	Title              string `form:"title" json:"title"`                           //题目标题
	Description        string `form:"description" json:"description"`               //题目描述
	InputDes           string `form:"inputDes" json:"inputDes"`                     //输入描述
	OutputDes          string `form:"outputDes" json:"outputDes"`                   //输出描述
	InputCase          string `form:"inputCase" json:"inputCase"`                   //测试输入
	OutputCase         string `form:"outputCase" json:"outputCase"`                 //测试输出
	Hint               string `form:"hint" json:"hint"`                             //题目提示(可以为对样例输入输出的解释)
	TimeLimit          int    `form:"timeLimit" json:"timeLimit"`                   //时间限制
	MemoryLimit        int    `form:"memoryLimit" json:"memoryLimit"`               //内存限制
	Tag                int64  `form:"tag" json:"tag"`                               //题目标签
	IsSpecialJudge     bool   `form:"isSpecialJudge" json:"isSpecialJudge"`         //是否特判
	SpecialJudgeSource string `form:"specialJudgeSource" json:"specialJudgeSource"` //特判程序源代码
	SpecialJudgeType   string `form:"specialJudgeType" json:"specialJudgeType"`     //特判程序源代码类型
	Code               string `form:"code" json:"code"`                             //标准程序
	LanguageLimit      string `form:"languageLimit" json:"languageLimit"`           //语言限制
	Remark             string `form:"remark" json:"remark"`                         //备注
}

func ProblemCreate(problem *Problem) (int64, error) {
	return OrmWeb.Insert(problem)
}

func ProblemRemove(id int64) error {
	_, err := OrmWeb.Id(id).Delete(&Problem{})
	return err
}

func ProblemUpdate(problem *Problem) error {
	_, err := OrmWeb.AllCols().ID(problem.Id).Update(problem)
	return err
}

func ProblemGetById(id int64) (*Problem, error) {
	problem := new(Problem)
	has, err := OrmWeb.Id(id).Get(problem)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return problem, nil
}

func ProblemGetByUserId(userId int64, currentPage int, perPage int) ([]*Problem, error) {
	problemList := make([]*Problem, 0)
	err := OrmWeb.Where("user_id=?", userId).Limit(perPage, (currentPage-1)*perPage).Find(&problemList)
	if err != nil {
		return nil, err
	}
	return problemList, nil
}

func ProblemCountByUserId(userId int64) (int64, error) {
	problem := &Problem{}
	count, err := OrmWeb.Where("user_id=?", userId).Count(problem)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getOriginList(origin string) []int64 {
	originMap := map[string]int64{
		"HDU":    1,
		"CodeVs": 2,
	}
	origins := []int64{}
	if origin != "" {
		strs := strings.Split(origin, ",")
		for i := 0; i < len(strs); i++ {
			origins = append(origins, originMap[strs[i]])
		}
	}
	return origins
}

func getNum(tag string) int {
	tagArr := map[string]int{
		"分治":   0,
		"贪心":   1,
		"字符串":  2,
		"动态规划": 3,
		"搜索":   4,
		"线性结构": 5,
		"链表":   6,
		"堆结构":  7,
	}
	num := 0
	if tag != "" {
		strs := strings.Split(tag, ",")
		for i := 0; i < len(strs); i++ {
			num += int(math.Pow(2, float64(tagArr[strs[i]])))
		}
	}
	return num
}

func ProblemGetIdsByConds(origins string, tag string) ([]*Problem, error) {
	session := OrmWeb.NewSession()
	originIds := getOriginList(origins)
	if origins != "" {
		session.In("user_id", originIds)
	}
	tagNum := getNum(tag)
	if tagNum != 0 {
		session.Where("tag & ? > 0", tagNum)
	}
	problemList := make([]*Problem, 0)

	err := session.Cols("id").Find(&problemList)
	if err != nil {
		return nil, err
	}
	return problemList, nil
}

func ProblemGetProblem(origins string, tag string, sortKey string, isAscKey string, currentPage int, perPage int) ([]*Problem, error) {
	session := OrmWeb.NewSession()
	originIds := getOriginList(origins)
	if origins != "" {
		session.In("user_id", originIds)
	}
	tagNum := getNum(tag)
	if tagNum != 0 {
		session.Where("tag & ? > 0", tagNum)
	}
	if isAscKey == "asc" {
		session.Asc(sortKey).Limit(perPage, (currentPage-1)*perPage)
	} else {
		session.Desc(sortKey).Limit(perPage, (currentPage-1)*perPage)
	}
	problemList := make([]*Problem, 0)

	err := session.Find(&problemList)
	if err != nil {
		return nil, err
	}
	return problemList, nil
}

func ProblemCountProblem(origins string, tag string) (int64, error) {
	session := OrmWeb.NewSession()
	originIds := getOriginList(origins)
	if origins != "" {
		session.In("user_id", originIds)
	}
	tagNum := getNum(tag)
	if tagNum != 0 {
		session.Where("tag & ? > 0", tagNum)
	}
	problem := &Problem{}
	count, err := session.Count(problem)
	if err != nil {
		return 0, err
	}
	return count, nil
}
