package omx

import (
	"crypto/tls"
	"fmt"
	"github.com/0x1un/go-zabbix"
	"log"
	"net/http"
)


type ZabbixConn struct {
	*zabbix.Session
}

func NewZabbixConn(url , username, password string) *ZabbixConn {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	cache := zabbix.NewSessionFileCache().SetFilePath("./zabbix_session")
	session, err := zabbix.CreateClient(url).
		WithCache(cache).
		WithHTTPClient(client).
		WithCredentials(username, password).
		Connect()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	version, err := session.GetVersion()

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("Connected to Zabbix API v%s\n", version)
	return &ZabbixConn{session}
}

func (z *ZabbixConn) ListHostIP() []string {
	params := zabbix.HostInterfaceGetParams{}
	params.Output = []string{"ip"}
	ifaces, err := z.HostInterfaceGet(params)
	if err != nil {
		panic(err)
	}
	ret := make([]string, 0)
	for _, iface := range ifaces {
		if iface.IP == "" {
			continue
		}
		ret = append(ret, iface.IP)
	}
	return ret
}