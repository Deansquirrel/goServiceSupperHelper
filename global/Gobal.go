package global

import (
	"context"
	"github.com/Deansquirrel/goToolMSSql"
)

const (
	//PreVersion = "1.0.1 Build20190807"
	//TestVersion = "0.0.0 Build20190101"
	Version = "1.0.2 Build20190807"
	Type    = "ServiceSupportHelper"
)

var (
	HeartBeatCron        = "0 * * * * ?"
	RefreshSvrV3InfoCron = "0 0/5 * * * ?"
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
