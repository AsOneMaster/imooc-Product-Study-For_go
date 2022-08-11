package common

import (
	"errors"
	"net"
)

func GetInterfaceIP() (string, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrList {
		//检查IP地址是否回环地址
		//address.(*net.IPNet)  将address由addr类型转换成IPNet型 类型断言会检查 addr 的动态类型是否满足 IPNet。 满足赋值 不满足nil
		// IsLoopback()方法作用为 ip 是否为环回地址。
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			//To4 将 IPv4 地址 ip 转换为 4 字节表示形式。如果 ip 不是 IPv4 地址，则 To4 返回 nil。
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("获取地址异常！！！")
}
