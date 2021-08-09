package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/0x1un/omx"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"strings"
)

/*
	omx zabbix check host ips.txt
*/

var (
	App = cli.NewApp()
)

var (ZbxConn *omx.ZabbixConn)

func encodeB64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func decodeB64(s string) (string, string, string) {
	dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println(".x_zbx文件损坏")
		return "", "", ""
	}
	ret := strings.Split(string(dec), " ")
	if len(ret) == 3 {
		return ret[0], ret[1], ret[2]
	}
	fmt.Println(".x_zbx文件损坏")
	return "","", ""
}

func readBin() (string, string, string) {
	data, err := ioutil.ReadFile(".x_zbx")
	if err != nil {
		fmt.Println(err)
		return "", "", ""
	}
	return decodeB64(string(data))
}

func writeBin(s []byte) {
	fp, err := os.OpenFile(".x_zbx", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {panic(err)}
	defer fp.Close()
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, s)
	fp.Write(buf.Bytes())
}

func zbxInit() {
	validate := func(input string) error {
		if len(input) < 6 {
			return errors.New("密码长度必须大于等于6")
		}
		return nil
	}

	url := promptui.Prompt{
		Label: "zabbix地址",
		Validate: nil,
	}

	username := promptui.Prompt{
		Label: "用户名",
		Validate: nil,
	}

	passwd := promptui.Prompt{
		Label:    "密码",
		Validate: validate,
		Mask:     '*',
	}

	retUrl, _ := url.Run()
	retUsername, _ := username.Run()
	retPasswd, _:= passwd.Run()
	enc := fmt.Sprintf("%s %s %s", retUrl, retUsername, retPasswd)
	writeBin([]byte(encodeB64(enc)))
}

func init() {
	App.Commands = []*cli.Command {
		{
			Name: "init",
			Usage: "初始化",
			Action: func(context *cli.Context) error {
				zbxInit()
				return nil
			},
		},
		{
			Name: "zabbix",
			Aliases: []string{"zbx", "zb"},
			Usage: "Zabbix操作命令",
			Subcommands: []*cli.Command{
				{
					Name: "check",
					Aliases: []string{"chk", "ck"},
					Usage: "检查命令",
					Subcommands: []*cli.Command{
						{
							Name: "host",
							Usage: "检查host",
							Flags: []cli.Flag{
								&cli.StringSliceFlag{
									Name:        "file",
									Usage:       "host文件名, 可接收多个, 以英文逗号分割",
									Required:    true,
								},
							},
							Action: func(c *cli.Context) error {
								if ZbxConn == nil {
									url, uname, passwd := readBin()
									ZbxConn = omx.NewZabbixConn(url, uname, passwd)
								}
								ZbxConn.CheckHostFromFile(c.StringSlice("file")...)
								return nil
							},
						},
						{
							Name: "ip",
							Usage: "输入IP地址, 以空格分割",
							Action: func(c *cli.Context) error {
								if ZbxConn == nil {
									url, uname, passwd := readBin()
									ZbxConn = omx.NewZabbixConn(url, uname, passwd)
								}
								ZbxConn.CheckHost(c.Args().Slice()...)
								return nil
							},
						},
					},
				},
			},
		},
	}
}