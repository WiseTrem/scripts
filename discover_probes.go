package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type monData struct {
	DATA []mon `xml:"mon"`
}

type mon struct {
	NAME string `xml:"name,attr"`
	DST  string `xml:"dst_addr,attr"`
	BIT  string `xml:"bitrate,attr"`
	CCE  string `xml:"cc_errors,attr"`
	DM   string `xml:"dst_mac,attr"`
	MLR  string `xml:"sum_mlr_1m,attr"`
}

func (m *monData) xmlGetDecode(r []byte) int {
	if err := xml.Unmarshal(r, &m); err != nil {
		fmt.Println(err)
	}
	return len(m.DATA)
}

func request(r string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return bodyByte
}

func GetAndParse(ip string) error {
	urlGet, err := url.Parse("http://template/probe/data/mon")
	if err != nil {
		return err
	}

	urlGet.Host = ip
	respSearch := request(urlGet.String())

	m := &monData{}
	if m.xmlGetDecode(respSearch) == 0 {
		err := fmt.Errorf(ip + " - Failed to get info")
		return err
	}

	b, err := json.Marshal(m)

	r := strings.NewReplacer("DATA", "data", "NAME", "{#NAME}", "DST", "{#ADDR}", "BIT", "{#BITRATE}", "CCE", "{#CCERR}", "DM", "{#DM}", "MLR", "{#MLR}")
	json := r.Replace(string(b))
	fmt.Println(json)
	return nil
}

func main() {
	arg1 := os.Args[1]
	err := GetAndParse(arg1)
	if err != nil {
		fmt.Println(err)
	}

}
