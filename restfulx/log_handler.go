package restfulx

import (
	"encoding/json"
	"fmt"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/logger"
	"github.com/XM-GO/PandaKit/utils"
	"reflect"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

type LogInfo struct {
	LogResp     bool   // 是否记录返回结果
	Description string // 请求描述
}

func NewLogInfo(description string) *LogInfo {
	return &LogInfo{Description: description, LogResp: false}
}

func (i *LogInfo) WithLogResp(logResp bool) *LogInfo {
	i.LogResp = logResp
	return i
}

func LogHandler(rc *ReqCtx) error {
	li := rc.LogInfo
	if li == nil {
		return nil
	}

	lfs := logrus.Fields{}
	if la := rc.LoginAccount; la != nil {
		lfs["uid"] = la.UserId
		lfs["uname"] = la.UserName
	}

	req := rc.Request.Request
	lfs[req.Method] = req.URL.Path

	if err := rc.Err; err != nil {
		logger.Log.WithFields(lfs).Error(getErrMsg(rc, err))
		return nil
	}
	logger.Log.WithFields(lfs).Info(getLogMsg(rc))
	return nil
}

func getLogMsg(rc *ReqCtx) string {
	msg := rc.LogInfo.Description + fmt.Sprintf(" ->%dms", rc.timed)
	if !utils.IsBlank(reflect.ValueOf(rc.ReqParam)) {
		rb, _ := json.Marshal(rc.ReqParam)
		msg = msg + fmt.Sprintf("\n--> %s", string(rb))
	}

	// 返回结果不为空，则记录返回结果
	if rc.LogInfo.LogResp && !utils.IsBlank(reflect.ValueOf(rc.ResData)) {
		respB, _ := json.Marshal(rc.ResData)
		msg = msg + fmt.Sprintf("\n<-- %s", string(respB))
	}
	return msg
}

func getErrMsg(rc *ReqCtx, err any) string {
	msg := rc.LogInfo.Description
	if !utils.IsBlank(reflect.ValueOf(rc.ReqParam)) {
		rb, _ := json.Marshal(rc.ReqParam)
		msg = msg + fmt.Sprintf("\n--> %s", string(rb))
	}

	var errMsg string
	switch t := err.(type) {
	case *biz.BizError:
		errMsg = fmt.Sprintf("\n<-e errCode: %d, errMsg: %s", t.Code(), t.Error())
	case error:
		errMsg = fmt.Sprintf("\n<-e errMsg: %s\n%s", t.Error(), string(debug.Stack()))
	case string:
		errMsg = fmt.Sprintf("\n<-e errMsg: %s\n%s", t, string(debug.Stack()))
	}
	return (msg + errMsg)
}
