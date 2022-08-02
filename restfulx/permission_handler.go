package restfulx

import (
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/casbin"
	"github.com/XM-GO/PandaKit/config"
	"github.com/XM-GO/PandaKit/token"
	"github.com/dgrijalva/jwt-go"
	"strconv"
)

type Permission struct {
	NeedToken  bool // 是否需要token
	NeedCasbin bool // 是否进行权限  api路径权限验证
}

func (p *Permission) WithNeedToken(needToken bool) *Permission {
	p.NeedToken = needToken
	return p
}

func (p *Permission) WithNeedCasBin(needCasBin bool) *Permission {
	p.NeedCasbin = needCasBin
	return p
}

func PermissionHandler(rc *ReqCtx) error {
	permission := rc.RequiredPermission
	// 如果需要的权限信息不为空，并且不需要token，则不返回错误，继续后续逻辑
	if permission != nil && !permission.NeedToken {
		return nil
	}
	tokenStr := rc.Request.Request.Header.Get("X-TOKEN")
	// header不存在则从查询参数token中获取
	if tokenStr == "" {
		tokenStr = rc.Request.QueryParameter("token")
	}
	if tokenStr == "" {
		return biz.PermissionErr
	}
	j := token.NewJWT("", []byte(config.Conf.Jwt.Key), jwt.SigningMethodHS256)
	loginAccount, err := j.ParseToken(tokenStr)
	if err != nil || loginAccount == nil {
		return biz.PermissionErr
	}
	rc.LoginAccount = loginAccount

	if !permission.NeedCasbin {
		return nil
	}
	e := casbin.Casbin()
	// 判断策略中是否存在
	tenantId := strconv.Itoa(int(rc.LoginAccount.TenantId))
	success, err := e.Enforce(tenantId, loginAccount.RoleKey, rc.Request.Request.URL.Path, rc.Request.Request.Method)
	if !success {
		return biz.CasbinErr
	}

	return nil
}
