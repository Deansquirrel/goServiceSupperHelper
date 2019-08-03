package goServiceSupportHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goServiceSupportHelper/global"
	"github.com/Deansquirrel/goServiceSupportHelper/object"
	"github.com/Deansquirrel/goToolCommon"
	"io/ioutil"
	"net/http"
	"strings"
)

func FormatSSJob(jobKey string, cmd func()) func() {
	return func() {
		jobId := NewJobId()
		_ = JobRecordStart(jobId, global.ClientId, jobKey)
		defer func() {
			_ = JobRecordEnd(jobId, global.ClientId, jobKey)
		}()
		cmd()
	}
}

func NewJobId() string {
	return strings.Replace(goToolCommon.Guid(), "-", "", -1)
}

func JobRecordStart(jobId, clientId, jobKey string) error {
	if strings.Trim(global.HttpAddress, " ") == "" {
		return errors.New("HttpAddress is empty")
	}
	d := object.JobRecordRequest{
		JobId:    jobId,
		ClientId: clientId,
		JobKey:   jobKey,
	}
	bd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/job/start", global.HttpAddress),
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

func JobRecordEnd(jobId, clientId, jobKey string) error {
	if strings.Trim(global.HttpAddress, " ") == "" {
		return errors.New("HttpAddress is empty")
	}
	d := object.JobRecordRequest{
		JobId:    jobId,
		ClientId: clientId,
		JobKey:   jobKey,
	}
	bd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	resp, err := http.Post(
		fmt.Sprintf("%s/job/end", global.HttpAddress),
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
