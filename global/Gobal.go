package global

import (
	"context"
	"github.com/Deansquirrel/goToolMSSql"
)

const (
	//PreVersion = "0.0.0 Build20190101"
	//TestVersion = "0.0.0 Build20190101"
	Version              = "0.0.0 Build20190101"
	Type                 = "ServiceSupportHelper"
	HeartBeatCron        = "0/10 * * * * ?"
	RefreshSvrV3InfoCron = "0/15 * * * * ?"
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
	IsSvrV3       = false
)
