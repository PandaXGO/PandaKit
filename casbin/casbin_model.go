package casbin

type CasbinRule struct {
	Ptype    string `json:"ptype" gorm:"column:ptype"`
	TenantId string `json:"tenantId" gorm:"column:v0"`
	RoleKey  string `json:"roleKey" gorm:"column:v1"`
	Path     string `json:"path" gorm:"column:v2"`
	Method   string `json:"method" gorm:"column:v3"`
	V4       string `json:"v4" gorm:"column:v4"`
	V5       string `json:"v5" gorm:"column:v5"`
	Id       int    `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
}

func (CasbinRule) TableName() string {
	return "casbin_rule"
}
