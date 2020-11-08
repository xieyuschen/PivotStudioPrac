package errmsg

//处理错误

const SUCCSE = 200
const ERROR = 500

//code=1000:用户模块错误
const ERROR_USERNAME_USED = 1001
const ERROR_PASSWORD_WRONG = 1002
const ERROR_USER_NOT_EXIST = 1003
const ERROR_TOKEN_EXIST = 1004
const ERROR_TOKEN_RUNTIME = 1005
const ERROR_TOKEN_WRONG = 1006
const ERROR_TOKEN_TYPE_WRONG = 1007

//code=2000：文章模块错误
const ERROR_ARTICLE_NOT_EXIST = 2001


var codeMsg = map[int]string{
	SUCCSE : "OK",
	ERROR : "FAIL",
	ERROR_USERNAME_USED : "用户名已存在",
	ERROR_PASSWORD_WRONG : "密码错误",
	ERROR_USER_NOT_EXIST : "用户不存在",
	ERROR_TOKEN_EXIST : "TOKEN不存在",
	ERROR_TOKEN_RUNTIME : "TOKEN已过期",
	ERROR_TOKEN_WRONG : "TOKEN错误",
	ERROR_TOKEN_TYPE_WRONG:  "TOKEN格式错误",
	ERROR_ARTICLE_NOT_EXIST:  "文章不存在",
}

func GetErrMsg(code int)string{
	//返回错误信息
	return codeMsg[code]
}

