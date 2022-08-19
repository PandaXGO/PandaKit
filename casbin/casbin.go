package casbin

import (
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/starter"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"sync"
)

type CasbinS struct {
	ModelPath string
}

func (c *CasbinS) UpdateCasbin(tenantId string, roleKey string, casbinInfos []CasbinRule) error {
	c.ClearCasbin(0, tenantId, roleKey)
	rules := [][]string{}
	for _, v := range casbinInfos {
		rules = append(rules, []string{tenantId, roleKey, v.Path, v.Method})
	}
	e := c.Casbin()
	success, _ := e.AddPolicies(rules)
	biz.IsTrue(success, "存在相同api,添加失败,请联系管理员")
	return nil
}

func (c *CasbinS) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) {
	err := starter.Db.Table("casbin_rule").Model(&CasbinRule{}).Where("v2 = ? AND v3 = ?", oldPath, oldMethod).Updates(map[string]any{
		"v2": newPath,
		"v3": newMethod,
	}).Error
	biz.ErrIsNil(err, "修改api失败")
}

func (c *CasbinS) GetPolicyPathByRoleId(tenantId, roleKey string) (pathMaps []CasbinRule) {
	e := c.Casbin()
	list := e.GetFilteredPolicy(0, tenantId, roleKey)
	for _, v := range list {
		pathMaps = append(pathMaps, CasbinRule{
			Path:   v[2],
			Method: v[3],
		})
	}
	return pathMaps
}

func (c *CasbinS) ClearCasbin(v int, p ...string) bool {
	e := c.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success

}

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func (c *CasbinS) Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(starter.Db)
		biz.ErrIsNil(err, "新建权限适配器失败")
		syncedEnforcer, err = casbin.NewSyncedEnforcer(c.ModelPath, a)
		biz.ErrIsNil(err, "新建权限适配器失败")
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}
