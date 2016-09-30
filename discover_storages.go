package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type response struct {
	//XMLName xml.Name `xml:"RESPONSE"`
	OBJECT []object
}

type object struct {
	PROPERTY []property
	TYPE     string `xml:"basetype,attr"`
}

type property struct {
	NAME  string `xml:"name,attr"`
	VALUE string `xml:",chardata"`
}

func (m *response) xmlGetDecode(r []byte) int {
	if err := xml.Unmarshal(r, &m); err != nil {
		fmt.Println(err)
	}
	return len(m.OBJECT)
}

func request(r string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("{\"data\":[]}")
		os.Exit(0)
	}
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return bodyByte
}

func requestCookie(r string, key string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", r, nil)
	wbisessionkey := "wbisessionkey=" + key
	req.Header.Add("cookie", wbisessionkey)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("{\"data\":[]}")
		os.Exit(0)
	}
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return bodyByte
}

func GetKey(ip string) (string, error) {
	urlGet, err := url.Parse("http://template/api/login/6bf512b52c242dcda04d7fdff8072f70")
	if err != nil {
		return "", err
	}

	urlGet.Host = ip
	respSearch := request(urlGet.String())

	m := &response{}
	if m.xmlGetDecode(respSearch) == 0 {
		err := fmt.Errorf(ip + " - Failed to get info")
		return "", err
	}
	b, err := json.Marshal(m)

	re := regexp.MustCompile(`[0-9a-f]{32}`)
	wbisessionkey := re.FindAllString(string(b), 1)
	return wbisessionkey[0], nil
}

func GetAndParse(ip string, key string) error {
	urlGet, err := url.Parse("http://template/api/show/volumes")
	if err != nil {
		return err
	}

	urlGet.Host = ip
	respSearch := requestCookie(urlGet.String(), key)

	m := &response{}
	if m.xmlGetDecode(respSearch) == 0 {
		err := fmt.Errorf(ip + " - Failed to get info")
		return err
	}
	b, err := json.Marshal(m)
	re := regexp.MustCompile(`\{"PROPERTY":\[\{"NAME":"durable-id","VALUE":".{2}"\},\{"NAME":"virtual-disk-name","VALUE"`)
	//test := re.FindAllString(string(b), -1)
	str := re.ReplaceAllLiteralString(string(b), "{\"{#NAME}\"")
	//fmt.Println(str)
	re = regexp.MustCompile(`\},\{"PROPERTY"[^\]]*][^\]]*`)
	newstr := re.ReplaceAllLiteralString(str, "}")
	//fmt.Println(newstr)

	r := strings.NewReplacer("OBJECT", "data", "},{\"NAME\":", ",", ",\"VALUE\"", "", "}],\"TYPE\":\"volumes\"", "",
		"storage-pool-name", "{#POOLNAME}", "reserved-size-in-pages", "{#RESSIZEPAGES}", "volume-name",
		"{#VOLNAME}", "total-size-numeric", "{#TOTALSIZENUM}", "total-size", "{#TOTALSIZE}",
		"allocated-size-numeric", "{#ALOCSIZENUM}", "allocated-size", "{#ALLOCSIZE}",
		"read-ahead-size-numeric", "{#READAHEADSIZENUM}", "read-ahead-size", "{#READAHEADRSIZE}",
		"size-numeric", "{#SIZENUM}", "size", "{#SIZE}", "preferred-owner-numeric", "{#PREFOWNERNUM}",
		"preferred-owner", "{#PREFOWNER}", "owner-numeric", "{#OWNERNUM}", "owner", "{#OWNER}",
		"serial-number", "{#SERIAL}", "write-policy-numeric", "{#WPOLICYNUM}", "write-policy", "{#WPOLICY}",
		"cache-optimization-numeric", "{#CACHEOPTNUM}", "cache-optimization", "{#CACHEOPT}",
		"volume-type-numeric", "{#VOLUMETYPENUM}", "volume-type", "{#VOLUMETYPE}",
		"volume-class-numeric", "{VOLCLASSNUM}", "volume-class", "{#VOLCLASS}",
		"profile-preference-numeric", "{#PROFILEPREFNUM}", "profile-preference", "{#PROFILEPREF}",
		"snapshot", "{#SNAPSHOT}", "volume-qualifier-numeric", "{#VOLUMEQNUM}", "volume-qualifier", "{#VOLUMEQ}",
		"blocks", "{#BLOCKS}", "capabilities", "{#CAPAB}", "volume-parent", "{#VOLPARENT}", "snap-pool", "{#SPOOL}",
		"replication-set", "{#RSET}", "attributes", "{#ATTR}", "virtual-disk-serial", "{#VDSERIAL}",
		"volume-description", "{#DESC}", "wwn", "{#WWN}", "progress-numeric", "{#PROGRESSRNUM}",
		"progress", "{#PROGRESS}", "container-name", "{#CNAME}", "container-serial", "{#CSERIAL}",
		"allowed-storage-tiers-numeric", "{#ASTN}", "allowed-storage-tiers", "{#AST}",
		"threshold-percent-of-pool", "{#TRESHP}", "allocate-reserved-pages-first-numeric", "{#ARPFNUM}",
		"allocate-reserved-pages-first", "{#ARPF}", "zero-init-page-on-allocation-numeric", "{#ZIPOANUM}",
		"zero-init-page-on-allocation", "{#ZIPOA}", "raidtype-numeric", "{#RAIDTYPENUM}", "raidtype", "{#RAIDTYPE}",
		"pi-format-numeric", "{#PIFORMATNUM}", "pi-format", "{#PIFORMAT}", "health-reason", "{#HEALTHREASON}",
		"health-recommendation", "{#HEALTHRECOM}", "health-numeric", "{#HEALTHNUM}", "health", "{#HEALTH}",
		"volume-group", "{#VOLGROUP}", "group-key", "{#GROUPKEY}")
	json := r.Replace(newstr)

	fmt.Println(json)
	return nil
}

func main() {
	arg1 := os.Args[1]
	key, err := GetKey(arg1)
	if err != nil {
		fmt.Println(err)
	}
	err = GetAndParse(arg1, key)
	if err != nil {
		fmt.Println(err)
	}

}
