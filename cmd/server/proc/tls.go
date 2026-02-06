package proc

import "gitee.com/lishimeng/event-bus/internal/message"

var UserLocalCipher = false

// LocalCipher 本地密钥
var LocalCipher message.ChannelCipher
