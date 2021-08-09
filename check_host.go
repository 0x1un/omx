package omx

import (
	"fmt"
	"strings"
)

func (z *ZabbixConn) CheckHost(ip ...string) {
	ips := z.ListHostIP()
	for _, i := range ip {
		if !Contains(ips, i) {
			fmt.Printf("%s 不存在zabbix中\n", i)
		}
	}
}

func (z *ZabbixConn) CheckHostFromFile(files ...string) {
	ipCache := make([]string, 0)
	for _, file := range files {
		ipCache = append(ipCache, ReadFileLines(file)...)
	}
	ips := z.ListHostIP()
	strBuffer := strings.Builder{}
	for _, ip := range ipCache {
		if !Contains(ips, ip) {
			strBuffer.WriteString(ip + " 不存在zabbix中\n")
		}
	}
	if strBuffer.Len() != 0 {
		fmt.Println(strBuffer.String())
	} else {
		fmt.Println("这看起来一切正常")
	}

}