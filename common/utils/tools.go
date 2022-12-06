package utils

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
)

func GetSubnetIp() (ip string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
						return
					}
				}
			}
		}
	}
	err = errors.New("can't find subnet")
	return
}

func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

func GetNowServerRootPath() string {
	fPath, _ := os.Getwd()

	configPath := flag.String("c", fPath, "root path")

	flag.Parse()

	return *configPath
}
