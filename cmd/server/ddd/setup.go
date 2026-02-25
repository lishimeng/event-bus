package ddd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"

	"gitee.com/lishimeng/event-bus/cmd/server/proc"
	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/providers/RocketMqProvider"
	"gitee.com/lishimeng/event-bus/providers/RocketMqProvider/proxy"
	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/container"
)

var sysConfigs map[string]db.SysConfig

func AfterWeb(ctx context.Context) (err error) {

	return
}

func BeforeWeb(ctx context.Context) (err error) {

	sysConfigs = make(map[string]db.SysConfig)
	sysConfigs, err = loadSysConfigs()
	if err != nil {
		return
	}

	err = initRmq(ctx) // 在channel加载之前

	err = loadLocalSecret(ctx) // 加载本地密钥
	if err != nil {
		return
	}
	err = loadChannels(ctx) // 加载每个通道
	if err != nil {
		return
	}
	return
}

// 加载channel
func loadChannels(_ context.Context) (err error) {

	var list []db.ChannelConfig
	_, err = app.GetOrm().Context.QueryTable(new(db.ChannelConfig)).All(&list)
	if err != nil {
		return
	}
	log.Info("load channels %d", len(list))
	for _, item := range list {
		err = channel.LoadChannel(item)
		if err != nil {
			log.Info("load channel fail, %s[%s]", item.Code, item.Name)
			return
		}
		log.Info("load channel success, %s[%s]", item.Code, item.Name)
	}
	return
}

func loadLocalSecret(_ context.Context) (err error) {
	var list []db.SysConfig
	_, err = app.GetOrm().Context.
		QueryTable(new(db.SysConfig)).
		Filter("name", db.SysLocalSecret).All(&list)
	if err != nil {
		return
	}
	log.Info("load local_secret:%d", len(list))
	if len(list) == 0 {
		log.Info("no local_secret")
		return
	}
	var conf = list[0]
	bs, err := base64.StdEncoding.DecodeString(conf.Config)
	if err != nil {
		return
	}
	var localSecret db.LocalSecurity
	err = json.Unmarshal(bs, &localSecret)
	if err != nil {
		return
	}
	proc.LocalCipher.RsaPriKey, err = cypher.LoadPrivateKey([]byte(localSecret.RsaKey))
	if err != nil {
		return
	} // 本地只需要私钥即可(加密操作)
	proc.UserLocalCipher = true

	return
}

func initRmq(ctx context.Context) (err error) {
	log.Info("setup rmq")
	var cfg RocketMqProvider.RmqConfig
	err = getSysConfig(db.SysRmqConfig, &cfg)
	if err != nil {
		return
	}
	client := proxy.New(ctx,
		proxy.WithEndpoint(cfg.Endpoint),
		proxy.WithAuth(cfg.AppId, cfg.Secret),
		proxy.WithPublisherConfigs(cfg.Publisher.MessageGroup, cfg.Publisher.Topics...),
	)
	container.Add(&client)

	log.Info("setup rmq provider")
	rmqProvider := RocketMqProvider.New(client, cfg)
	container.Add(&rmqProvider)
	engine := proc.NewEngine(rmqProvider)
	proc.EngineInstance = engine // 初始化message engine
	return
}

func getSysConfig(key string, ptr any) (err error) {
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
