package admin

import (
	"github.com/lishimeng/app-starter/server"
)

func Router(root server.Router) {

	localSecret := root.Path("/local_secret")
	localSecret.Get("/", apiGetLocalSecret)
	localSecret.Post("/", apiConfigLocalSecret)

	ch := root.Path("/channel")    //通道管理
	ch.Get("/", apiListChannels)   //获取通道列表
	ch.Post("/", apiChannelConfig) //添加通道

	rmq := root.Path("/rmq_config") //Rmq配置管理
	rmq.Get("/", apiGetRmqConfig)   //获取Rmq配置
	rmq.Post("/", apiSaveRmqConfig) //保存Rmq配置

	rec := root.Path("/records")     //记录管理
	rec.Get("/", apiListDataRecords) //获取记录列表
}
