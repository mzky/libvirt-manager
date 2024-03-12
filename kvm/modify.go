package kvm

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"libvirt-manager/utils"
	"strconv"
	"strings"
)

type NetInfo struct {
	VMName  string `json:"name"`
	IP      string `json:"ip"`
	NetType string `json:"netType"`
	NetPort string `json:"netPort"`
}

// ModifyVMMemory 修改虚拟机内存
func ModifyVMMemory(name, memory string) (string, error) {
	memorySize, e := utils.StringToUint64(memory)
	if e != nil {
		return "", e
	}
	if err := SetVMMaxMemory(name, memorySize); err != nil {
		return "", err
	}
	// 命令方式不支持max内存修改,需要先修改最大内存后才可修改当前内存
	return PublicVMControl("setmem", name, strconv.FormatUint(memorySize, 10), "--config")
}

// ModifyVMCPU 修改虚拟机CPU #重启后生效
func ModifyVMCPU(name, vcpu string) (string, error) {
	if _, err := PublicVMControl("setvcpus", name, vcpu, "--config"); err != nil {
		return "", err
	}
	return PublicVMControl("setvcpus", name, "--maximum", vcpu, "--config")
}

// ModifyVMDisk 修改虚拟机磁盘大小
func ModifyVMDisk(name, newSize string) (string, error) {
	return utils.Exec("qemu-img", "resize", utils.JoinVMPath(name), strings.ToUpper(newSize))
}

// GetVMMac 获取指定虚拟机的MAC
func GetVMMac(name string) (string, error) {
	dom, err := PublicVMControl("domiflist", name)
	if err != nil {
		return "", err
	}
	line := strings.Split(dom, "\n")[2]
	mac := strings.Fields(line)[4]
	return mac, nil
}

// JoinHostXML 拼接网络配置的XML
func JoinHostXML(name string, ip string) string {
	mac, err := GetVMMac(name)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return fmt.Sprintf("<host mac='%s' name='%s' ip='%s' />", mac, name, ip)
}

// AddNetIP 添加虚拟机IP地址
/*
修改IP原理，从虚拟机获取MAC，通过配置虚拟路由器的MAC和IP绑定（这里使用默认的default，后期如果手动创建网络，需要根据实际情况修改）
添加DHCP配置将这个MAC绑定一个IP，虚拟机启动时将自动获取到此IP，注意这里需要重启一些虚拟路由器，也就是default
virsh net-list
取MAC：
virsh dumpxml name | grep 'mac address'
virsh domiflist ttt
virsh domifaddr ttt
修改路由器配置：
virsh net-edit default
重启虚拟路由器：
virsh net-destroy default
virsh net-start default
启动和自启动默认网络：
virsh net-start default
virsh net-autostart default
*/
func (ni *NetInfo) AddNetIP() (string, error) {
	xml := JoinHostXML(ni.VMName, ni.IP)
	return PublicVMControl("net-update", ni.NetPort, "add", "ip-dhcp-host", xml, "--live", "--config")
}

// ModifyNetIP 修改虚拟机IP地址
func (ni *NetInfo) ModifyNetIP() (string, error) {
	xml := JoinHostXML(ni.VMName, ni.IP)
	return PublicVMControl("net-update", ni.NetPort, "modify", "ip-dhcp-host", xml, "--live", "--config")
}

// DeleteNetIP 删除虚拟机IP地址
func (ni *NetInfo) DeleteNetIP() (string, error) {
	xml := JoinHostXML(ni.VMName, ni.IP)
	return PublicVMControl("net-update", ni.NetPort, "delete", "ip-dhcp-host", xml, "--live", "--config")
}

// NetRestart 重启默认default网络，后期根据实际修改
func (ni *NetInfo) NetRestart() error {
	if _, err := PublicVMControl("net-destroy", ni.NetPort); err != nil {
		return err
	}
	if _, err := PublicVMControl("net-start", ni.NetPort); err != nil {
		return err
	}
	return nil
}

// 查看已绑定MAC/IP的虚拟机 virsh net-dhcp-leases default
/*
# virsh net-dhcp-leases default
Expiry Time           MAC               Protocol       IP address      Hostname   Client ID or DUID
---------------------------------------------------------------------------------------------------------
2024-01-11 11:58:26   52:54:00:e3:01:35   ipv4       192.168.122.100/24   ttt        -

# virsh domifaddr ttt
 Name       MAC address          Protocol     Address
-------------------------------------------------------------------------------
 vnet0      52:54:00:e3:01:35    ipv4         192.168.122.100/24
*/

type Leases struct {
	ExpiryTime string
	MAC        string
	Protocol   string
	Address    string
	Hostname   string
}

// NetList 虚拟路由器DHCP分配情况
func (ni *NetInfo) NetList() ([]Leases, error) {
	var lesList []Leases
	leases, err := PublicVMControl("net-dhcp-leases", ni.NetPort)
	if err != nil {
		return lesList, err
	}

	lines := strings.Split(strings.TrimSpace(leases), "\n")
	for _, line := range lines[2:] {
		var les Leases
		lin := strings.Fields(line)
		les.ExpiryTime = lin[0] + " " + lin[1]
		les.MAC = lin[2]
		les.Protocol = lin[3]
		les.Address = lin[4]
		les.Hostname = lin[5]
		lesList = append(lesList, les)
	}

	return lesList, nil
}

type AddrInfo struct {
	Name     string
	MAC      string
	Protocol string
	Address  string
}

// DomainAddr 虚拟机网络配置
func (ni *NetInfo) DomainAddr() ([]AddrInfo, error) {
	var addrList []AddrInfo
	addr, err := PublicVMControl("domifaddr", ni.VMName)
	if err != nil {
		return addrList, err
	}
	lines := strings.Split(strings.TrimSpace(addr), "\n")
	for _, line := range lines[2:] {
		var ai AddrInfo
		lin := strings.Fields(line)
		ai.Name = lin[0]
		ai.MAC = lin[1]
		ai.Protocol = lin[2]
		ai.Address = lin[3]
		addrList = append(addrList, ai)
	}

	return addrList, nil
}

/*
nmcli con add ifname br0 type bridge con-name br0
nmcli con add type bridge-slave ifname enp125s0f0 master br0
nmcli c modify br0 ipv4.addresses '192.168.227.31/24'  ipv4.gateway '192.168.227.254' bridge.stp no
nmcli con up br0
nmcli -f bridge con show br0

创建bridge.xml内容
<network>
  <name>bridge</name>
  <forward mode='bridge'/>
  <model type='virtio'/>
  <bridge name='br0'/>
</network>

删除br
virsh net-destroy bridge
从xml创建br
virsh net-define bridge.xml
启动
virsh net-start bridge
自启
virsh net-autostart bridge

virsh net-list --all
*/
