/**
 * Created by shiyi on 2017/12/21.
 * Email: shiyi@fightcoder.com
 */

package judge

const (
	Normal              int = 0  //默认
	Waiting             int = 1  //WT	用户程序正在排队等待测试
	Compiling           int = 2  //CING	编译中
	Running             int = 3  //RING	运行中
	Accepted            int = 4  //AC	用户程序输出正确的结果
	WrongAnswer         int = 5  //WA	用户程序输出错误的结果
	CompilationError    int = 6  //CE	用户程序编译错误
	TimeLimitExceeded   int = 7  //TLE	用户程序运行时间超过题目的限制
	MemoryLimitExceeded int = 8  //MLE	用户程序运行内存超过题目的限制
	OutputLimitExceeded int = 9  //OLE	用户程序输出的结果大大超出正确答案的长度
	RuntimeError        int = 10 //RE	用户程序发生运行时错误
	SystemError         int = 11 //SE	用户程序不能被评测系统正常运行
)

type Result struct {
	ResultCode    int
	ResultDes     string
	RunningTime   int64 //耗时(ms)
	RunningMemory int64 //所占空间
}
