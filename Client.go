package goServiceSupportHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
	"github.com/Deansquirrel/goServiceSupportHelper/object"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type client struct {
}

func NewClient() *client {
	return &client{}
}

func (c *client) RefreshClientInfo(
	clientId string, clientType string, clientVersion string, hostName string, dbId int,
	dbName string, internetIP string) error {
	if strings.Trim(global.HttpAddress, " ") == "" {
		return errors.New("HttpAddress is empty")
	}
	d := object.ClientInfoRequest{
		ClientId:      clientId,
		ClientType:    clientType,
		ClientVersion: clientVersion,
		HostName:      hostName,
		DbId:          dbId,
		DbName:        dbName,
		InternetIP:    internetIP,
	}
	bd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/client/info", global.HttpAddress),
		"application/json",
		bytes.NewReader(bd))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rd object.Response
	err = json.Unmarshal(body, &rd)
	if err != nil {
		return err
	}
	if rd.ErrCode != 200 {
		return errors.New(rd.ErrMsg)
	}
	return nil
}

func (c *client) RefreshSvrV3Info(
	clientId string,
	coId int, coAb string, coCode string, coUserAb string, coUserCode string, coFunc int,
	svName string, svVer string, svDate time.Time) error {
	if strings.Trim(global.HttpAddress, " ") == "" {
		return errors.New("HttpAddress is empty")
	}
	d := object.SvrV3InfoRequest{
		ClientId:   clientId,
		CoId:       coId,
		CoAb:       coAb,
		CoCode:     coCode,
		CoUserAb:   coUserAb,
		CoUserCode: coUserCode,
		CoFunc:     coFunc,
		SvName:     svName,
		SvVer:      svVer,
		SvDate:     svDate,
	}
	bd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/client/svrv3", global.HttpAddress),
		"application/json",
		bytes.NewReader(bd))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rd object.Response
	err = json.Unmarshal(body, &rd)
	if err != nil {
		return err
	}
	if rd.ErrCode != 200 {
		return errors.New(rd.ErrMsg)
	}
	return nil
}
