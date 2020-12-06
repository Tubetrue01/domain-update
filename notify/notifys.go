package notify

import "org.tubetrue01/domain-update/config"

// Notify 定义通知抽象层
type Notify interface {
	// DoNotify 通知方法
	DoNotify(*config.Config, interface{})

	// DoNotifyBefore 发送通知前的处理操作
	DoNotifyBefore(*config.Config, interface{})
}
