package goServiceSupportHelper

import (
	"context"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
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

	waitForClientId()

	//HeartBeat
	go func() {
		for {
			err := goToolCron.AddFunc(
				"HeartBeat",
				global.HeartBeatCron,
				NewJob().FormatSSJob("HeartBeat", jobHeartBeat),
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

func GetVersion() string {
	return global.Version
}

func GetType() string {
	return global.Type
}

type SvrType string

const (
	SVRNONE SvrType = "svrNone"
	SVRV3   SvrType = "svrV3"
	SVRZ5   SvrType = "svrZ5"
)

func SetOtherInfo(dbConfig *goToolMSSql.MSSqlConfig,
	dbType int,
	svrType SvrType) {
	waitForClientId()
	global.DbConfig = dbConfig
	global.DbType = dbType

	if global.DbConfig == nil {
		return
	}
	go refreshDbId(global.DbConfig, global.DbType)
	switch svrType {
	case SVRV3:
		go func() {
			for {
				err := goToolCron.AddFunc(
					"RefreshSvrV3Info",
					global.RefreshSvrV3InfoCron,
					NewJob().FormatSSJob("RefreshSvrV3Info", jobRefreshSvrV3Info),
					panicHandle)
				if err == nil {
					break
				} else {
					time.Sleep(time.Second * 10)
				}
			}
		}()
	default:
		log.Warn(fmt.Sprintf("unexpected type: %s", string(svrType)))
	}
}

func panicHandle(v interface{}) {
	log.Error(fmt.Sprintf("panicHandle: %s", v))
}

func getClientId() string {
	for {
		if global.ClientType == "" {
			log.Warn("client type can not be empty")
			time.Sleep(time.Second * 10)
			continue
		}
		id, err := goToolEnvironment.GetClientId(global.ClientType)
		if err != nil {
			log.Warn(fmt.Sprintf("get client id err: %s", err.Error()))
			time.Sleep(time.Second * 10)
			continue
		}
		return id
	}
}

//刷新global.InternetIp
func refreshInternetIp() {
	waitForClientId()
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
	waitForClientId()
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
	waitForClientId()
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
	waitForClientId()
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

func jobHeartBeat(id string) {
	waitForClientId()
	err := NewHeartBeat().HeartBeatUpdate()
	if err != nil {
		log.Error(err.Error())
		err = JobErrRecord(id, err.Error())
		if err != nil {
			log.Error(fmt.Sprintf("record err err: %s", err.Error()))
		}
	}
	return
}

func jobRefreshSvrV3Info(id string) {
	waitForClientId()
	coId, coAb, coCode, coUserAb, coUserCode, coFunc, err :=
		goToolSVRV3.GetZlCompany(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		jobRefreshSvrV3InfoHandleErr(id, err)
		return
	}
	svName, svVer, svDate, err := goToolSVRV3.GetXtSelfVer(goToolMSSqlHelper.ConvertDbConfigTo2000(global.DbConfig))
	if err != nil {
		jobRefreshSvrV3InfoHandleErr(id, err)
		return
	}
	err = NewClient().RefreshSvrV3Info(
		global.ClientId,
		coId, coAb, coCode, coUserAb, coUserCode, coFunc,
		svName, svVer, svDate)
	if err != nil {
		jobRefreshSvrV3InfoHandleErr(id, err)
		return
	}
}

func jobRefreshSvrV3InfoHandleErr(id string, err error) {
	log.Error(err.Error())
	err = JobErrRecord(id, err.Error())
	if err != nil {
		log.Error(fmt.Sprintf("record err err: %s", err.Error()))
	}
}

func waitForClientId() {
	for {
		if global.ClientId == "" {
			time.Sleep(time.Second)
			continue
		}
		break
	}
}
