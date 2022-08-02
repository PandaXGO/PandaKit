package config

import (
	"flag"
	"fmt"
	"github.com/XM-GO/PandaKit/biz"
	"github.com/XM-GO/PandaKit/utils"
	"path/filepath"
)

var Conf *Config

func InitConfig(configFilePath string) *Config {
	// 获取启动参数中，配置文件的绝对路径
	path, _ := filepath.Abs(configFilePath)
	startConfigParam = &CmdConfigParam{ConfigFilePath: path}
	// 读取配置文件信息
	yc := &Config{}
	if err := utils.LoadYml(startConfigParam.ConfigFilePath, yc); err != nil {
		panic(any(fmt.Sprintf("读取配置文件[%s]失败: %s", startConfigParam.ConfigFilePath, err.Error())))
	}
	// 校验配置文件内容信息
	yc.Valid()

	Conf = yc
	return yc

}

// 启动配置参数
type CmdConfigParam struct {
	ConfigFilePath string // -e  配置文件路径
}

// 启动可执行文件时的参数
var startConfigParam *CmdConfigParam

// yaml配置文件映射对象
type Config struct {
	App        *App        `yaml:"app"`
	Server     *Server     `yaml:"server"`
	Jwt        *Jwt        `yaml:"jwt"`
	Redis      *Redis      `yaml:"redis"`
	Mysql      *Mysql      `yaml:"mysql"`
	Postgresql *Postgresql `yaml:"postgresql"`
	Casbin     *Casbin     `yaml:"casbin"`
	Gen        *Gen        `yaml:"gen"`
	Log        *Log        `yaml:"log"`
}

// 配置文件内容校验
func (c *Config) Valid() {
	biz.IsTrue(c.Jwt != nil, "配置文件的[jwt]信息不能为空")
	c.Jwt.Valid()
}

// 获取执行可执行文件时，指定的启动参数
func getStartConfig() *CmdConfigParam {
	configFilePath := flag.String("e", "./config.yml", "配置文件路径，默认为可执行文件目录")
	flag.Parse()
	// 获取配置文件绝对路径
	path, _ := filepath.Abs(*configFilePath)
	sc := &CmdConfigParam{ConfigFilePath: path}
	return sc
}
