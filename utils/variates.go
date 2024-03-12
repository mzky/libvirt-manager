package utils

const (
	NetworkTypeBR  = "bridge"  // BR
	NetworkTypeNAT = "network" // NAT

	VMStateRunning  = "running"  // 虚拟机正在运行。
	VMStatePaused   = "paused"   // 虚拟机已被暂停。
	VMStateShutdown = "shutdown" // 虚拟机被正常关闭。
	VMStateShutoff  = "shutoff"  // 虚拟机被强制关闭或者由于错误而关闭。
	VMStateCrashed  = "crashed"  // 虚拟机崩溃。
	VMStateUnknown  = "unknown"  // 状态未知。

	VMAutoStart = "--autostart"
	VMNullArg   = "--force"

	// 错误吗示例
	noErrorMsg       = "unknown error"
	STATUSTEST       = 1900
	STATUSGETVMLIST  = 1901
	STATUSGETVMSTATE = 1902
	STATUSVMSTART    = 1903
	STATUSVMSTOP     = 1904
	STATUSVMPAUSE    = 1905
	STATUSVMRESUME   = 1906
	STATUSVMRESTART  = 1907
	STATUSVMDELETE   = 1908
	STATUSCREATEVM   = 1909
)

var (
	ISOPath   = "/home/ISO"               // ISO 镜像目录，可修改
	VMPath    = "/var/lib/libvirt/images" // 默认虚拟机文件位置，不建议修改
	ImagePath = "/home/image"             // 模板镜像路径，可修改

	Port = "9100" // 服务端口
)

var messages = map[int]string{
	STATUSTEST:       "test",
	STATUSGETVMLIST:  "获取虚拟机列表失败",
	STATUSGETVMSTATE: "获取虚拟机状态失败",
	STATUSVMSTART:    "启动虚拟机失败",
	STATUSVMSTOP:     "停止虚拟机失败",
	STATUSVMPAUSE:    "挂起虚拟机失败",
	STATUSVMRESUME:   "恢复虚拟机失败",
	STATUSVMRESTART:  "重启虚拟机失败",
	STATUSVMDELETE:   "删除虚拟机失败",
	STATUSCREATEVM:   "创建虚拟机失败",
}

// StatusText ret error msg
func StatusText(code int) string {
	msg, ok := messages[code]
	if !ok {
		msg = noErrorMsg
	}
	return msg
}
