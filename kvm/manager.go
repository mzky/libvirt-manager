package kvm

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"libvirt-manager/utils"
	"os"
	"strings"
)

type VMConfig struct {
	VMName    string `json:"vmName"`
	Memory    string `json:"memory"`
	Vcpus     string `json:"vcpus"`
	DiskSize  string `json:"diskSize"` // G
	IsoName   string `json:"isoName"`  // CreateISO2VM
	QcowTmpl  string `json:"qcowTmpl"` // CreateVM2VM
	NetType   string `json:"netType"`
	NetPort   string `json:"netPort"`
	VmRemak   string `json:"remarks"`
	AutoStart bool   `json:"isAutoStart"`
}

type QcowInfo struct {
	Id            string `yaml:"Id"`
	Name          string `yaml:"Name"`
	UUID          string `yaml:"UUID"`
	OSType        string `yaml:"OS Type"`
	State         string `yaml:"State"`
	CPUs          string `yaml:"CPU(s)"`
	MaxMemory     string `yaml:"Max memory"`
	UsedMemory    string `yaml:"Used memory"`
	Persistent    string `yaml:"Persistent"`
	Autostart     string `yaml:"Autostart"`
	ManagedSave   string `yaml:"Managed save"`
	SecurityModel string `yaml:"Security model"`
	SecurityDOI   string `yaml:"Security DOI"`
}

// CreateQcowFile 创建虚拟机磁盘文件
func (vmc *VMConfig) CreateQcowFile() (string, error) {
	return utils.Exec("qemu-img", "create", "-f", "qcow2", utils.JoinVMPath(vmc.VMName), vmc.DiskSize+"G")
}

func (vmc *VMConfig) retNetInfo() string {
	nt := fmt.Sprintf("bridge=%s,model=virtio", vmc.NetPort)
	if vmc.NetType == utils.NetworkTypeNAT {
		nt = "network=default,model=virtio"
	}
	return nt
}

// CreateVM2VM 根据模板镜像创建虚拟机
func (vmc *VMConfig) CreateVM2VM() (string, error) {
	return utils.Exec("virt-install",
		"--name", vmc.VMName,
		"--memory", vmc.Memory,
		"--vcpus", vmc.Vcpus,
		"--cpu", "host",
		"--machine=virt",
		"--arch=aarch64",
		"--os-type", "linux",
		"--os-variant", "generic",
		"--disk", "path="+utils.JoinVMPath(vmc.VMName)+",format=qcow2,bus=virtio",
		"--network", vmc.retNetInfo(),
		"--graphics", "vnc,listen=0.0.0.0",
		"--import", // 使用镜像文件创建虚拟机
		"--noautoconsole",
		utils.IsAutoStart(vmc.AutoStart),
	)
}

/**
--network source=network: 指定网络的源类型，可以是"bridge", "network", "direct", "hostdev", 等等。
          bridge: 指定一个物理网桥。虚拟机的网络接口将被添加到这个网桥上，从而能够与物理网络通信。例如：--network bridge=br0。
          network: 指定一个已经在libvirt中定义的网络。这可以是一个 NAT 网络、一个桥接网络或者其他类型的网络。例如：--network network=my_network.
          direct: 直接绑定到主机的一个物理网络接口。这种模式通常用于高性能网络或者需要低延迟的场景。例如：--network source=direct,source_dev=eth0.
--network model=virtio: 指定网络模型，如"virtio", "e1000", "ne2k_pci", 等等。
--network address=IP_address/netmask: 指定虚拟机的IP地址和子网掩码。
--network dhcp=yes: 使用DHCP自动获取IP地址和其他网络配置。
*/

// CreateISO2VM 根据ISO镜像创建虚拟机
func (vmc *VMConfig) CreateISO2VM() (string, error) {
	return utils.Exec("virt-install",
		"--virt-type=kvm",
		"--name", vmc.VMName,
		"--memory", vmc.Memory,
		"--vcpus", vmc.Vcpus,
		"--cpu", "host",
		"--disk", "path="+utils.JoinVMPath(vmc.VMName)+",format=qcow2,bus=virtio",
		"--machine", "virt",
		"--arch", "aarch64",
		"--os-type", "linux",
		"--os-variant", "generic",
		"--network", vmc.retNetInfo(),
		"--cdrom", utils.JoinISOPath(vmc.IsoName),
		"--graphics", "vnc,listen=0.0.0.0",
		"--noautoconsole",
		utils.IsAutoStart(vmc.AutoStart),
	)
}

// VMDelete 删除虚拟机
func VMDelete(name string) error {
	if _, err := PublicVMControl("destroy", name); err != nil {
		return fmt.Errorf("failed to stop VM: %v", err)
	}

	if _, err := PublicVMControl("undefine", name, "--nvram"); err != nil {
		return fmt.Errorf("failed to delete VM: %v", err)
	}

	// 删除虚拟机磁盘文件
	if err := os.Remove(utils.JoinVMPath(name)); err != nil {
		return fmt.Errorf("failed to delete VM disk image: %v", err)
	}

	return nil
}

// VMStart 启动虚拟机
func VMStart(name string) (string, error) {
	return PublicVMControl("start", name)
}

// VMRestart 重启虚拟机
func VMRestart(name string) (string, error) {
	return PublicVMControl("reboot", name)
}

// VMPause 关闭虚拟机
func VMPause(name string) (string, error) {
	return PublicVMControl("suspend", name)
}

// VMStop 暂停/挂起虚拟机
func VMStop(name string) (string, error) {
	return PublicVMControl("shutdown", name)
}

// VMResume 恢复暂停/挂起的虚拟机
func VMResume(name string) (string, error) {
	return PublicVMControl("resume", name)
}

// PublicVMControl 公共vm控制
func PublicVMControl(arg ...string) (string, error) {
	cmd, err := utils.Exec("virsh", arg...)
	if err != nil {
		return "", fmt.Errorf("VM: %v", err)
	}
	return cmd, nil
}

type VmList struct {
	ID    string
	Name  string
	State string
}

// VMLists 获取虚拟机列表
func VMLists() ([]VmList, error) {
	output, err := PublicVMControl("list", "--all")
	if err != nil {
		return nil, fmt.Errorf("VM list cannot be obtained: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var vms []VmList
	for _, line := range lines[2:] {
		var vmList VmList
		lin := strings.Fields(line)
		vmList.ID = lin[0]
		vmList.Name = lin[1]
		vmList.State = lin[2]
		if len(lin) > 3 {
			vmList.State = lin[2] + lin[3]
		}
		vms = append(vms, vmList)
	}

	return vms, nil
}

// VMState 虚拟机状态
func VMState(name string) (string, error) {
	return PublicVMControl("domstate", name)
}

// VMInfo 虚拟机信息
func VMInfo(name string) (QcowInfo, error) {
	var vmc QcowInfo
	v, err := PublicVMControl("dominfo", name)
	if err != nil {
		return vmc, err
	}
	v = strings.ReplaceAll(v, "-", "") // 必须加，否则报错
	if e := yaml.Unmarshal([]byte(strings.TrimSpace(v)), &vmc); e != nil {
		fmt.Println("Error unmarshalling YAML:", e)
		return vmc, e
	}
	return vmc, nil
}

// GetVNCDisplay 查看指定虚拟的vnc端口
func GetVNCDisplay(vmName string) (string, error) {
	return PublicVMControl("vncdisplay", vmName)
}

// VMConvert 将虚拟机压缩成模板镜像
func VMConvert(src, dst string) (string, error) {
	imagePath := utils.JoinImagePath(dst)
	if utils.FileExist(imagePath) {
		return "", fmt.Errorf("VM file already exist")
	}
	return utils.Exec("qemu-img", "convert", "-c", "-O", "qcow2", utils.JoinVMPath(src), utils.JoinImagePath(dst))
}

// CloneImage2VM 从模板克隆虚拟机文件
func CloneImage2VM(src, dst string) error {
	vmPath := utils.JoinVMPath(dst)
	if utils.FileExist(vmPath) {
		return fmt.Errorf("VM file already exist")
	}
	return utils.CopyFile(utils.JoinImagePath(src), vmPath)
}

// Autostart 开机启动
func Autostart(name string, start bool) (string, error) {
	if start {
		return PublicVMControl("autostart", name)
	}
	return PublicVMControl("autostart", "--disable", name)
}
