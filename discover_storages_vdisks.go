package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatal(err)
	}
	b, err := json.Marshal(m)

	re := regexp.MustCompile(`"response","VALUE":"[^"]*`)
	ans := re.FindString(string(b))
	re = regexp.MustCompile(`"response","VALUE":"`)
	wbiSessionKey := re.ReplaceAllLiteralString(ans, "")
	if wbiSessionKey == "Authentication Unsuccessful" {
		err := fmt.Errorf(ip + " - Authentication Unsuccessful")
		log.Fatal(err)
	}
	return wbiSessionKey, nil
}

func GetAndParse(ip string, key string) error {
	urlGet, err := url.Parse("http://template/api/show/Vdisks")
	if err != nil {
		return err
	}

	urlGet.Host = ip
	respSearch := requestCookie(urlGet.String(), key)

	m := &response{}
	if m.xmlGetDecode(respSearch) == 0 {
		err := fmt.Errorf(ip + " - Failed to get info")
		log.Fatal(err)
	}
	b, err := json.Marshal(m)
	re := regexp.MustCompile(`\{"PROPERTY":\[\{"NAME":"name","VALUE":"`)
	//fmt.Println(string(b)")
	str := re.ReplaceAllLiteralString(string(b), "{\"{#NAME}\":\"")
	//fmt.Println(str)
	re = regexp.MustCompile(`],"TYPE":"virtual-disks"\}(,\{"PROPERTY"[^\]]*][^\]]*)?`)
	newstr := re.ReplaceAllLiteralString(str, "")
	//fmt.Println(newstr)

	r := strings.NewReplacer(
		"OBJECT", "data",
		"},{\"NAME\":", ",",
		",\"VALUE\"", "",
		"lun", "{#LUN}",
		"freespace-numeric", "{#FREESPACENUM}",
		"freespace", "{#FREESPACE}",
		"diskcount", "{#DISKCOUNT}",
		"total-size-numeric", "{#TOTALSIZENUM}",
		"total-size", "{#TOTALSIZE}",
		"chunksize", "{#CHUNKSIZE}",
		"allocated-size-numeric", "{#ALOCSIZENUM}",
		"allocated-size", "{#ALLOCSIZE}",
		"sparecount", "{#SPACECOUNT}",
		"min-drive-size-numeric", "{#MINDRSIZESIZENUM}",
		"min-drive-size", "{#MINDRSIZE}",
		"create-date-numeric", "{#CREATEDATENUM}",
		"create-date", "{#CREATEDATE}",
		"cache-read-ahead-numeric", "{#CACHEREADAHEADNUM}",
		"cache-read-ahead", "{#CACHEREADAHEAD}",
		"cache-flush-period", "{#CACHEFLUSHPERIOD}",
		"read-ahead-enabled-numeric", "{#RAENABLEDNUM}",
		"read-ahead-enabled", "{#RAENABLED}",
		"write-back-enabled-numeric", "{#WBENABLEDNUM}",
		"write-back-enabled", "{#WBENABLED}",
		"status-numeric", "{#STATUSNUM}",
		"status", "{#STATUS}",
		"job-running", "{#JOB RUNNING}",
		"current-job-numeric", "{#CURJOBNUM}",
		"current-job-completion", "{#CURJOBCOMPLETION}",
		"current-job", "{#CURJOB}",
		"num-array-partitions", "{#NUMARRPART}",
		"largest-free-partition-space-numeric", "{#LFREEPARTSPACENUM}",
		"largest-free-partition-space", "{#LFREEPARTSPACE}",
		"num-drives-per-low-level-array", "{#NUMDRPERLOWARR}",
		"num-expansion-partitions", "{#NEP}",
		"num-partition-segments", "{#NPS}",
		"new-partition-lba-numeric", "{#NPLBAN}",
		"new-partition-lba", "{#NPLBA}",
		"array-drive-type-numeric", "{#DRIVETYPENUM}",
		"array-drive-type", "{#DRIVETYPE}",
		"is-job-auto-abortable-numeric", "{#IJAAN}",
		"is-job-auto-abortable", "{#IJAA}",
		"disk-dsd-enable-vdisk-numeric", "{#DDEVN}",
		"disk-dsd-enable-vdisk", "{#DDEV}",
		"disk-dsd-delay-vdisk", "{#DDDV}",
		"scrub-duration-goal", "{#SDG}",
		"read-ahead-size-numeric", "{#READAHEADSIZENUM}",
		"read-ahead-size", "{#READAHEADRSIZE}",
		"size-numeric", "{#SIZENUM}",
		"size", "{#SIZE}",
		"preferred-owner-numeric", "{#PREFOWNERNUM}",
		"preferred-owner", "{#PREFOWNER}",
		"owner-numeric", "{#OWNERNUM}",
		"owner", "{#OWNER}",
		"serial-number", "{#SERIAL}",
		"blocks", "{#BLOCKS}",
		"capabilities", "{#CAPAB}",
		"volume-parent", "{#VOLPARENT}",
		"snap-pool", "{#SPOOL}",
		"replication-set", "{#RSET}",
		"raidtype-numeric", "{#RAIDTYPENUM}",
		"raidtype", "{#RAIDTYPE}",
		"health-blame-numeric", "{#HEALTHBLAMENUM}",
		"health-blame", "{#HEALTHBLAME}",
		"health-reason", "{#HEALTHREASON}",
		"health-recommendation", "{#HEALTHRECOM}",
		"health-numeric", "{#HEALTHNUM}",
		"health", "{#HEALTH}",
		"number-of-reads", "{#NUMREADS}",
		"number-of-writes", "{#NUMWRITES}",
		"total-sectors-read", "{#SECTORREAD}",
		"total-sectors-written", "{#SECTORWRITTN}",
		"target-id", "{#TARGETID}",
		"total-data-transferred-numeric", "{#DATATRANSNUM}",
		"total-data-transferred", "{#DATATRANS}",
		"total-bytes-per-sec-numeric", "{#BPSNUM}",
		"total-bytes-per-sec", "{#BPS}")
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
