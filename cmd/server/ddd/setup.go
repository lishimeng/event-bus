package ddd

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/event-bus/cmd/server/proc"
	"github.com/lishimeng/event-bus/internal/channel"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/domains/sysCfg"
	"github.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/event-bus/internal/tls/cypher"
	"github.com/lishimeng/event-bus/providers/RocketMqProvider"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/container"
)

func AfterWeb(ctx context.Context) (err error) {

	return
}

func BeforeWeb(ctx context.Context) (err error) {

	sysCfg.EnableCache()

	err = loadChannels(ctx) // 加载每个通道
	if err != nil {
		return
	}

	subscribeTopics := channel.GetManager().SubscribeTopics()
	publishTopics := channel.GetManager().PublishTopics()

	err = initRmq(ctx, publishTopics, subscribeTopics) // 在channel加载之前
	if err != nil {
		return
	}
	// rmq不需要显示调用subscribe
	// TODO 切换mqtt时需要调度器
	return
}

// 加载channel
func loadChannels(_ context.Context) (err error) {

	var list []db.ChannelConfig
	_, err = app.GetOrm().Context.QueryTable(new(db.ChannelConfig)).All(&list)
	if err != nil {
		return
	}
	log.Info("load channels [%d]", len(list))
	defer func() {
		log.Info("load channels done")
	}()
	for _, item := range list {
		log.Info("load channel %s[%s]:%s, tls:%d", item.Code, item.Name, item.Category.String(), item.UseSecurity)
		var ch message.Channel
		ch, err = channel.LoadChannel(item)
		if err != nil {
			log.Info("load channel fail, %s[%s]", item.Code, item.Name)
			return
		}
		log.Info("load channel success, %s[%s]", ch.Code, ch.Name)
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

func initRmq(ctx context.Context, publishTopics []string, subscribeTopics []string) (err error) {
	// rmq需要在启动时做subscribe操作,不支持动态订阅. publish也需要预先汇总全部出口topic, 实际的subscribe接口不生效
	log.Info("setup rmq")
	var cfg RocketMqProvider.RmqConfig
	err = sysCfg.GetSysConfig(db.SysRmqConfig, &cfg)
	if err != nil {
		return
	}

	log.Info("setup rmq provider")
	rmqProvider := RocketMqProvider.New(ctx, cfg, publishTopics, subscribeTopics)
	container.Add(&rmqProvider)
	engine := proc.NewEngine(rmqProvider)
	rmqProvider.SetMessageListener(engine.OnMessage)
	proc.EngineInstance = engine // 初始化message engine
	return
}
