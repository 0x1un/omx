package omx

import (
	"fmt"
	"testing"
)

func TestParseAINIFile(t *testing.T) {
	IniHostObj.ParseAINIFile("testdata/ips.txt")
	//fmt.Println(a.Data[0].Hosts)
	fmt.Println(IniHostObj.GetHostnameByIP("10.21.11.7"))
	fmt.Println(IniHostObj.GetAllIP())
}
