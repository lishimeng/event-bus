package sysCfg

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/event-bus/internal/db"
)

var cacheEnable = false
var sysConfigs map[string]db.SysConfig

func init() {
	sysConfigs = map[string]db.SysConfig{}
}

func EnableCache() {
	cacheEnable = true
}

func GetSysConfig(key string, ptr any) (err error) {
	conf, ok := getCacheSysConfig(key)
	if !ok {
		err = errors.New("no rmq config")
		return
	}
	bs, err := base64.StdEncoding.DecodeString(conf.Config)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, ptr)
	return
}

func getCacheSysConfig(name string) (conf db.SysConfig, exist bool) {
	var err error
	if cacheEnable { // 使用cache
		conf, exist = sysConfigs[name]
	}
	if exist {
		return
	}
	conf, exist = sysConfigs[name]
	sysConfigs, err = loadSysConfigs()
	if err != nil {
		exist = false
		return
	}
	conf, exist = sysConfigs[name]
	return
}

func loadSysConfigs() (configs map[string]db.SysConfig, err error) {
	configs = map[string]db.SysConfig{}
	var list []db.SysConfig
	_, err = app.GetOrm().Context.QueryTable(new(db.SysConfig)).All(&list)
	if err != nil {
		return
	}
	for _, item := range list {
		configs[item.Name] = item
	}
	return
}
