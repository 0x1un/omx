package omx

import (
	"bytes"
	"fmt"
	"github.com/0x1un/go-zabbix"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"os"
	"strings"
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

type CsvTemplate struct {
	Hostname  string `csv:"主机名"`
	HostIP    string `csv:"主机IP"`
	Groups    string `csv:"主机组"`
	Templates string `csv:"主机模板"`
	Proxy     string `csv:"主机代理"`
	NotUsed   string `csv:"-"`
}

func ConvertXLSX2Csv(xlsxFile string) ([]byte, error) {


	return nil, nil
}

func (z *ZabbixConn) AddHostFromCsvFile(csvFiles ...string) {
	for _, csvFile := range csvFiles {
		data, err := ioutil.ReadFile(csvFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		data = bytes.Trim(data, "\xef\xbb\xbf")
		var csvTemplates []*CsvTemplate
		if err := gocsv.UnmarshalBytes(data, &csvTemplates); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		for _, csvValue := range csvTemplates {
			z.AddHost(csvValue.Hostname, csvValue.HostIP, csvValue.Groups, csvValue.Proxy, strings.Split(csvValue.Templates, ","))
		}
	}
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
	groupID, err := z.GetGroupIDByName(group)
	if err != nil {
		// TODO create new group if groupid not exists
		fmt.Fprintln(os.Stderr, "主机组查询错误:", err)
		return
	}
	hostParams.Groups = []zabbix.Group{
		{
			GroupID: groupID,
		},
	}

	proxyID, err := z.GetProxiesIDByName(proxy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "代理查询错误:", err)
		return
	}
	hostParams.ProxyIDS = proxyID

	templateIDS, err := z.GetTemplateIDByName(templates)
	if err != nil {
		fmt.Fprintln(os.Stderr, "模板查询错误:", err)
		return
	}
	for _, tid := range templateIDS {
		hostParams.Templates = append(hostParams.Templates, zabbix.TemplateObj{
			TemplateID: tid,
		})
	}
	_, err = z.CreateHost(hostParams)
	if err != nil {
		fmt.Fprintln(os.Stderr, "创建主机失败: "+err.Error())
		return
	}
	fmt.Println(hostname, ip, "新增成功")
}

func (z *ZabbixConn) GetTemplateIDByName(names []string) ([]string, error) {
	params := zabbix.TemplateGetParams{}
	params.Output = "extend"
	respData, err := z.GetTemplates(params)
	if err != nil {
		return nil, err
	}
	templateIDS := make([]string, 0)
	for _, name := range names {
		if name == "" {
			continue
		}

		for _, template := range respData {
			if strings.Trim(name, "") == template.Name {
				templateIDS = append(templateIDS, template.Templateid)
			}
		}
	}
	if len(templateIDS) == 0 {
		return nil, zabbix.ErrNotFound
	}
	return templateIDS, nil
}
