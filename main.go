package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	database "UrlFilter/db"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type accessLog struct {
	method   string
	route    string
	domain   string
	sourceIP string
	dateTime string
	params   map[string]string
}

func main() {
	db, err := gorm.Open("mysql", "root:test@(localhost:4444)/temp?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		panic("failed to connect database")
	}
	defer db.Close()

	fmt.Println("Connect Success")

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&database.AccessLog{}, &database.AccessLogParam{})

	fmt.Println("Auto migrate complete")

	// base := "./../../../20190415/"
	base := "/Users/1330/Dean/workspace/work/script/UrlFilter/logs/"

	files, err := ioutil.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileName := f.Name()
		fmt.Printf("目前檔案：%v\n", fileName)

		if isGzipFile(fileName) {
			fmt.Printf("解壓縮檔案：%v\n", fileName)
			unGzip(base + fileName)
		}

		fileName = fileName[0:len(fileName)]
		fmt.Printf("寫入檔案 %v\n", fileName)

		file, err := os.Open(base + fileName)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log := parseRowData(scanner.Text())

			a := &database.AccessLog{
				IP:          log.sourceIP,
				Method:      log.method,
				Domain:      getDomain(fileName),
				Route:       log.route,
				TriggeredAt: log.dateTime,
			}
			if err := db.Create(a).Error; err != nil {
				panic(err)
			}

			fmt.Printf("寫入資料 log id:%v \n", a.ID)

			if len(log.params) > 0 {
				// create AccessLogParam
				for index, item := range log.params {
					p := &database.AccessLogParam{
						AccessLogID: a.ID,
						Key:         index,
						Value:       item,
					}
					if err := db.Create(p).Error; err != nil {
						panic(err)
					}
					fmt.Printf("寫入資料 log -- param id:%v \n", p.ID)
				}
			}
		}
	}
	db.Close()

	fmt.Println("寫入完成")
}

func getDomain(fileName string) string {
	split := strings.Split(fileName, "-")
	return split[0]
}

func iterDir(base string) {
	files, err := ioutil.ReadDir(base)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func isGzipFile(fileName string) bool {
	matched, _ := regexp.MatchString(`.*\.gz`, fileName)

	return matched
}

func unGzip(fileName string) {
	// fmt.Println(fileName)
	cmd := exec.Command("gzip", "-d", fileName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
}

func parseRowData(row string) (a accessLog) {
	ipRe := regexp.MustCompile(`(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]) `)
	mainDataRe := regexp.MustCompile(`"(.*)\"`)
	dateRe := regexp.MustCompile(`\[(.*)\]`)

	ips := ipRe.FindAllString(row, -1)
	dates := dateRe.FindAllString(row, -1)
	mainDatas := mainDataRe.FindAllString(row, -1)

	// fmt.Printf("ips:%v,dates:%v,maindatas:%v\n", ips, dates, mainDatas)

	mainData := mainDatas[0]

	urlString := strings.Split(mainData, " ")
	u, err := url.Parse(urlString[1])
	if err != nil {
		fmt.Println(err)
		panic("failed to parse url")
	}
	queries := u.Query()

	// fmt.Printf("queries:%v\n", queries)

	params := make(map[string]string)
	for key, value := range queries {
		params[key] = value[0]
	}

	a.sourceIP = strings.TrimSpace(ips[0])
	a.dateTime = dates[0]
	a.dateTime = a.dateTime[1 : len(a.dateTime)-1]
	a.route = u.Path
	a.method = urlString[0]
	a.method = a.method[1:]
	a.params = params

	return
}
