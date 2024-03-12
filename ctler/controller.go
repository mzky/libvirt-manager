package ctler

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"libvirt-manager/kvm"
	"libvirt-manager/utils"
	"net/http"
)

func Test(c *gin.Context) {
	var vmc kvm.VMConfig
	vmc.VMName = "ttt"
	vmc.Memory = "2048" // 默认单位M
	vmc.Vcpus = "2"
	vmc.DiskSize = "20" // 默认单位G
	vmc.IsoName = "openEuler-22.03-LTS-SP2-aarch64-dvd.iso"
	vmc.QcowTmpl = "temp2"             // 模板名
	vmc.NetType = utils.NetworkTypeNAT // bridge or network
	vmc.NetPort = "virbr0"
	s, err := kvm.VMInfo(vmc.VMName)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSTEST, err))
		return
	}
	c.JSON(http.StatusOK, utils.Success(s))
}

func CreateVm(c *gin.Context) {
	var vmc kvm.VMConfig
	body, _ := c.GetRawData()
	_ = json.Unmarshal(body, &vmc)
	//ip, _ := jsonparser.GetString(body, "ip")
	imagetype, _ := jsonparser.GetString(body, "imageType")
	/*vmc.VMName, _ = jsonparser.GetString(body, "vmname")
	  vmc.Memory, _ = jsonparser.GetString(body, "memory")
	  vmc.Vcpus, _ = jsonparser.GetString(body, "cpu")
	  vmc.DiskSize, _ = jsonparser.GetString(body, "disk")
	  vmc.AutoStart, _ = jsonparser.GetBoolean(body, "isautostart")
	  vmc.NetType = utils.NetworkTypeNAT // bridge or network
	  vmc.NetPort = "virbr0"
	  vmc.VmRemak, := jsonparser.GetString(body, "remarks")*/
	if imagetype == "iso" {
		s, err := vmc.CreateQcowFile()
		if err != nil {
			c.JSON(http.StatusOK, utils.CustomError(utils.STATUSCREATEVM, err))
			return
		}
		s, err = vmc.CreateISO2VM()
		if err != nil {
			c.JSON(http.StatusOK, utils.CustomError(utils.STATUSCREATEVM, err))
			return
		}
		fmt.Println(s)
		c.JSON(http.StatusOK, utils.Success(s))
	} else {
		err := kvm.CloneImage2VM(vmc.QcowTmpl, vmc.VMName)
		if err != nil {
			c.JSON(http.StatusOK, utils.CustomError(utils.STATUSCREATEVM, err))
			return
		}
		s, err := vmc.CreateISO2VM()
		if err != nil {
			c.JSON(http.StatusOK, utils.CustomError(utils.STATUSCREATEVM, err))
			return
		}
		fmt.Println(s)
		c.JSON(http.StatusOK, utils.Success(s))
	}
}

func GetVmList(c *gin.Context) {
	s, err := kvm.VMLists()
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSGETVMLIST, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func GetVMState(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	s, err := kvm.VMState(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSGETVMSTATE, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func VMStart(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	s, err := kvm.VMStart(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSVMSTART, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func VMStop(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	s, err := kvm.VMStop(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSVMSTOP, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func VMPause(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	s, err := kvm.VMPause(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSVMPAUSE, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func VMResume(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	s, err := kvm.VMResume(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSVMRESUME, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func VMRestart(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	s, err := kvm.VMRestart(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSVMRESTART, err))
		return
	}
	fmt.Println(s)
	c.JSON(http.StatusOK, utils.Success(s))
}

func VMDelete(c *gin.Context) {
	body, _ := c.GetRawData()
	name, _ := jsonparser.GetString(body, "name")
	err := kvm.VMDelete(name)
	if err != nil {
		c.JSON(http.StatusOK, utils.CustomError(utils.STATUSVMDELETE, err))
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func BasePkg() {
	var vmc kvm.VMConfig
	vmc.VMName = "ttt"
	vmc.Memory = "2048" // 默认单位M
	vmc.Vcpus = "2"
	vmc.DiskSize = "20" // 默认单位G
	vmc.IsoName = "openEuler-22.03-LTS-SP2-aarch64-dvd.iso"
	vmc.QcowTmpl = "temp2"             // 模板名
	vmc.NetType = utils.NetworkTypeNAT // bridge or network
	vmc.NetPort = "virbr0"             // 根据实际情况填写

	// 从ISO创建虚拟机
	fmt.Println(vmc.CreateQcowFile()) // 创建虚拟机磁盘文件(在创建虚拟机之前创建)
	fmt.Println(vmc.CreateISO2VM())   // 从ISO镜像创建虚拟机

	// 从模板创建虚拟机
	fmt.Println(kvm.CloneImage2VM(vmc.QcowTmpl, vmc.VMName)) // 从系统镜像模板文件克隆到默认虚拟机文件目录()
	fmt.Println(vmc.CreateVM2VM())                           // 根据镜像模板创建虚拟机

	s, _ := kvm.VMInfo(vmc.VMName) // 返回虚拟机信息（返回信息为结构体）
	fmt.Printf("%+v\n", s)

	fmt.Println(kvm.VMLists())           // 返回虚拟机列表
	fmt.Println(kvm.VMState(vmc.VMName)) // 返回虚拟机状态

	// 以下成对儿存在
	fmt.Println(kvm.VMStart(vmc.VMName)) // 启动虚拟机---需虚拟机为关机状态
	fmt.Println(kvm.VMStop(vmc.VMName))  // 关闭虚拟机---需虚拟机为开机状态

	fmt.Println(kvm.VMPause(vmc.VMName))  // 挂起虚拟机---需虚拟机为开机状态
	fmt.Println(kvm.VMResume(vmc.VMName)) // 恢复虚拟机---需虚拟机为挂起状态

	fmt.Println(kvm.VMRestart(vmc.VMName)) // 重启虚拟机---需要虚拟机为开机状态才可执行
	fmt.Println(kvm.VMDelete(vmc.VMName))  // 删除虚拟机，包括文件（注意删除后不能还原）

	fmt.Println(kvm.GetVNCDisplay(vmc.VMName))        // 虚拟机对应宿主机vnc的端口号，用于连接虚拟机（可使用web页面远程连接vnc访问虚拟机桌面）
	fmt.Println(kvm.VMConvert(vmc.VMName, "imgTemp")) // 将虚拟机压缩成模板镜像---需虚拟机为关机状态
	fmt.Println(kvm.CloneVM(vmc.VMName, "vmTemp"))    // 从存储池克隆已存在的虚拟机
	fmt.Println(kvm.Autostart(vmc.VMName, true))      // 虚拟机随宿主机开机自启动（参数2为true时开启自启动，false为禁用自启动）

	// 快照管理---注意快照所以操作，均需要虚拟机关机状态
	fmt.Println(kvm.CreateSnapshot(vmc.VMName, "snapshot")) // 创建虚拟机快照
	fmt.Println(kvm.DeleteSnapshot(vmc.VMName, "snapshot")) // 删除虚拟机快照
	fmt.Println(kvm.SnapshotList(vmc.VMName))               // 列出虚拟机快照
	fmt.Println(kvm.RevertSnapshot(vmc.VMName, "snapshot")) // 回滚快照

	// 虚拟机配置修改--以下配置需重启虚拟机生效
	fmt.Println(kvm.ModifyVMDisk(vmc.VMName, "20G"))  // 修改虚拟机磁盘大小，注意只能改大不能改小，单位支持G、M
	fmt.Println(kvm.ModifyVMCPU(vmc.VMName, "4"))     // 修改虚拟机CPU核数
	fmt.Println(kvm.ModifyVMMemory(vmc.VMName, "8G")) // 修改虚拟机内存大小，单位支持G、M

	var netInfo kvm.NetInfo
	netInfo.NetType = utils.NetworkTypeNAT
	netInfo.NetPort = "default"
	netInfo.VMName = "ttt"
	netInfo.IP = "192.168.122.120"

	fmt.Println(kvm.GetVMMac(vmc.VMName)) // 获取指定虚拟机的MAC
	fmt.Println(netInfo.AddNetIP())       // 添加IP（用于新建虚拟机时）
	fmt.Println(netInfo.ModifyNetIP())    // 修改IP（需要已有配置下的修改）
	fmt.Println(netInfo.DeleteNetIP())    // 删除IP绑定
	fmt.Println(netInfo.NetRestart())     // 以上配置重启虚拟网络才能生效
	fmt.Println(netInfo.NetList())        // 获取已分配IP的绑定信息
	fmt.Println(netInfo.DomainAddr())     // 获取指定虚拟机的网络信息

	fmt.Println(kvm.GetDomainXML(vmc.VMName))           // 获取指定虚拟机的xml配置（暂时没什么用，可通过修改xml内容来修改虚拟机各种配置）
	fmt.Println(kvm.SetInterfaceParameters(vmc.VMName)) // 修改虚拟机指定的参数
}
