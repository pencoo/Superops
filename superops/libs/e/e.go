package e

const (
	SUCCESS            = 200   //请求成功
	ERROR              = 400   //请求失败
	FAILED_PARAMS_NUM  = 40001 //参数数量错误
	FAILED_PARAMS_TYPE = 40002 //参数类型错误
	FAILED_PARAMS_ANAL = 40003 //参数解析失败
	FAILED_INSERT_DB   = 40004 //数据库写入失败
	FAILED_UPDATE_DB   = 40005 //数据库更新失败
	FAILED_SELECT_DB   = 40006 //数据查询失败
	FAILED_PARAMS_NULL = 40007 //有必须参数为空

	FAILED_AUTH       = 41001 //认证错误
	FAILED_AUTH_ROLE  = 41002 //权限认证错误
	FAILED_AUTH_TOKEN = 41003 //token认证错误
	FAILED_POWER      = 41004 //权限不足
)

var MsgMap = map[int]string{
	SUCCESS:            "ok",
	ERROR:              "数据解析失败",
	FAILED_PARAMS_NUM:  "请求参数数量错误",
	FAILED_PARAMS_TYPE: "请求参数类型错误",
	FAILED_PARAMS_ANAL: "请求数据绑定失败",
	FAILED_INSERT_DB:   "数据库写入失败",
	FAILED_UPDATE_DB:   "数据库更新失败",
	FAILED_SELECT_DB:   "数据查询失败",
	FAILED_PARAMS_NULL: "有必须参数为空",
	FAILED_AUTH:        "认证失败",
	FAILED_AUTH_ROLE:   "角色权限认证错误",
	FAILED_AUTH_TOKEN:  "Token解析错误",
	FAILED_POWER:       "权限不足",
}

func GetMsg(i int) string {
	msg, ok := MsgMap[i]
	if ok {
		return msg
	} else {
		return MsgMap[400]
	}
}
