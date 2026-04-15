# event-bus
事件总线, 发送/订阅消息

## channel用法
1. channel有通道方向(db.RouteCategory), Subscribe/PublishTo
2. Subscribe程序需要监听的通道: 加载后会执行监听
3. PublishTo发布通道: 程序可以发布消息到此通道中
4. RouteCategory+Route唯一(强制约束)
5. channel中有cipher, 支持通信加密
6. channel中配置callback, 实现外部业务系统集成
7. 加解密算法参考session
8. 加密密钥(公钥/私钥)在channel初始化时加载