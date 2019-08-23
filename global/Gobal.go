package global

import (
	"context"
	"github.com/Deansquirrel/goToolMSSql"
)

const (
	//PreVersion = "1.0.6 Build20190821"
	//TestVersion = "0.0.0 Build20190101"
	Version = "1.0.7 Build20190823"
	Type    = "ServiceSupportHelper"
)

var (
	HeartBeatCron      = "0 * * * * ?"
	RefreshSvrInfoCron = "0 0/5 * * * ?"
)

var Ctx context.Context
var Cancel func()
var DbConfig *goToolMSSql.MSSqlConfig

var (
	ClientId      = ""
	ClientType    = ""
	ClientVersion = ""
	HttpAddress   = ""
	HostName      = ""
	InternetIp    = ""
	DbId          = -1
	DbType        = -1
)
