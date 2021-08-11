package omx

import (
	"fmt"
	"github.com/relex/aini"
	"os"
)

var (
	IniHostObj = &IniHost{}
)

type IniHost struct {
	Data []*aini.InventoryData
}


func (i *IniHost) ParseAINIFile(filenames ...string) *IniHost {
	inventories := make([]*aini.InventoryData, 0)
	for _, filename := range filenames {
		inventory, err := aini.ParseFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		inventories = append(inventories, inventory)
	}
	i.Data = inventories
	return i
}

func (i *IniHost) GetAllIP() []string {
	ips := make([]string, 0)
	for _, inventory := range i.Data {
		for k:= range inventory.Hosts {
			ips = append(ips, k)
		}
	}
	return ips
}

func (i *IniHost) GetHostnameByIP(ip string) string {
	for _, inventory := range i.Data {
		for k, v := range inventory.Hosts {
			if hostname, ok := v.Vars["hostname"]; ip == k && ok  {
				return hostname
			}
		}
	}
	return ""
}