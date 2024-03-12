package main

import (
	"libvirt-manager/db"
	"libvirt-manager/log"
	"libvirt-manager/router"
	"libvirt-manager/utils"
)

func main() {
	log.InitLog("debug", "kvm.log") // 根据实际情况修改日志路径
	utils.Flag()
	db.DBInit()
	router.StartServer()

	//ctler.BasePkg()
}
