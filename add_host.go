package omx

import (
	"fmt"
	"github.com/0x1un/go-zabbix"
	"os"
)

func (z *ZabbixConn) GetGroupIDByName(name string) (string, error) {
	params := zabbix.HostgroupGetParams{}
	groups, err := z.GetHostgroups(params)
	if err != nil {
		return "", err
	}
	for _, group := range groups {
		if name == group.Name {
			return group.GroupID, nil
		}
	}
	return "", zabbix.ErrNotFound
}

func (z *ZabbixConn) AddHost(hostname, ip, group, proxy string, templates []string) {
	hostParams := zabbix.CreateHostRequest{}
	hostParams.Host = hostname
	hostParams.VisibleName = hostname
	hostParams.Status = zabbix.StatusEnabled
	hostParams.Interfaces = []zabbix.Interface{
		{
			Main:  1,
			Type:  zabbix.TypeAgent,
			Port:  10050,
			Bulk:  1,
			Useip: 1,
			IP:    ip,
		},
	}
	groupID, err  := z.GetGroupIDByName(group)
	if err != nil {
		// TODO create new group if groupid not exists
		fmt.Fprintln(os.Stderr, err)
		return
	}
	hostParams.Groups = []zabbix.Group {
		{
			GroupID: groupID,
		},
	}
}
