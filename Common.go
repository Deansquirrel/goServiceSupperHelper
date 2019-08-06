package goServiceSupportHelper

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/Deansquirrel/goToolCron"
	"github.com/Deansquirrel/goToolEnvironment"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/Deansquirrel/goToolMSSqlHelper"
	"github.com/Deansquirrel/goToolSVRV3"
	"strings"
	"time"
)

import log "github.com/Deansquirrel/goToolLog"

type Params struct {
	HttpAddress   string
	ClientType    string
	ClientVersion string
	Ctx           context.Context
	Cancel        func()
}

func InitParam(p *Params) {
	global.HttpAddress = strings.Trim(p.HttpAddress, " ")
	global.ClientType = strings.Trim(p.ClientType, " ")
	global.ClientVersion = strings.Trim(p.ClientVersion, " ")
	global.Ctx = p.Ctx
	global.Cancel = p.Cancel
	global.ClientId = getClientId()

	//HeartBeat
	go func() {
		for {
			err := goToolCron.AddFunc(
				"heartbeat",
				global.HeartBeatCron,
				NewJob().FormatSSJob("heartbeat", jobHeartBeat),
				panicHandle)
			if err == nil {
				break
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}()
	go refreshClientInfo()
	go refreshHostName()
	go refreshInternetIp()
}

func SetHeartBeatCron(spec string) {
	spec = strings.Trim(spec, " ")
	if spec != "" {
		global.HeartBeatCron = spec
	}
}

func SetRefreshSvrV3InfoCron(spec string) {
	spec = strings.Trim(spec, " ")
	if spec != "" {
		global.RefreshSvrV3InfoCron = spec
	}
}

func SetOtherInfo(dbConfig *goToolMSSql.MSSqlConfig,
	dbType int,
	isSvrV3 bool) {
	global.DbConfig = dbConfig
	global.DbType = dbType
	go func() {
		if global.DbConfig == nil {
			return
		}
		go refreshClientInfo()
		go refreshDbId(global.DbConfig, global.DbType)
		global.IsSvrV3 = isSvrV3
		if global.IsSvrV3 {
			go func() {
				for {
					err := goToolCron.AddFunc(
						"refreshSvrV3Info",
						global.RefreshSvrV3InfoCron,
						NewJob().FormatSSJob("refreshSvrV3Info", jobRefreshSvrV3Info),
						panicHandle)
					if err == nil {
						break
					} else {
						time.Sleep(time.Second * 10)
					}
				}
			}()
		}
	}()
}

func panicHandle(v interface{}) {
	log.Error(fmt.Sprintf("panicHandle: %s", v))
}

func getClientId() string {
	if global.ClientType == "" {
		time.Sleep(time.Second * 10)
		return getClientId()
	}
	biosSn, _ := goToolEnvironment.BIOSSerialNumber()
	diskSn, _ := goToolEnvironment.DiskDriverSerialNumber()
	return strings.ToUpper(goToolCommon.Md5([]byte(global.ClientType + biosSn + diskSn)))
}

//刷新global.InternetIp
func refreshInternetIp() {
	for {
		ip, err := goToolEnvironment.GetInternetAddr()
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.InternetIp = ip
		refreshClientInfo()
		break
	}
}

//刷新global.HostName
func refreshHostName() {
	for {
		hostName, err := goToolEnvironment.GetHostName()
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		global.HostName = hostName
		refreshClientInfo()
		break
	}
}

func refreshClientInfo() {
	for {
		dbName := ""
		if global.DbConfig != nil {
			dbName = global.DbConfig.DbName
		}
		err := NewClient().RefreshClientInfo(
			global.ClientId,
			global.ClientType,
			global.ClientVersion,
			global.HostName,
			global.DbId,
			dbName,
			global.InternetIp)
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		break
	}
}

//刷新global.DbId
func refreshDbId(dbConfig *goToolMSSql.MSSqlConfig, dbType int) {
root:
	for {
		switch dbType {
		case 0:
			dbId, err := goToolMSSqlHelper.GetDbId(dbConfig)
			if err != nil {
				time.Sleep(time.Minute * 5)
				continue
			} else {
				global.DbId = dbId
				break root
			}
		case 1:
			dbId, err := goToolMSSqlHelper.GetDbId2000(goToolMSSqlHelper.ConvertDbConfigTo2000(dbConfig))
			if err != nil {
				time.Sleep(time.Minute * 5)
				continue
			} else {
				global.DbId = dbId
				break root
			}
		default:
			return
		}
	}
	refreshClientInfo()
}

func jobHeartBeat() {
	err := NewHeartBeat().HeartBeatUpdate()
	if err != nil {
		log.Error(err.Error())
	}
	return
}

func jobRefreshSvrV3Info() {
	coId, coAb, coCode, coUserAb, coUserCode, coFunc, err :=
		goToolSVRV3.GetZlCompany(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		log.Error(err.Error())
		return
	}
	svName, svVer, svDate, err := goToolSVRV3.GetXtSelfVer(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = NewClient().RefreshSvrV3Info(
		global.ClientId,
		coId, coAb, coCode, coUserAb, coUserCode, coFunc,
		svName, svVer, svDate)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
