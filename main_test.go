package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegex(t *testing.T) {
	s := "asdf.gz"

	actual := isGzipFile(s)

	assert.True(t, actual)
}

func TestSplitDomain(t *testing.T) {
	fileName := "ssl.gomaji.com-access_log-20190414"

	actual := getDomain(fileName)

	expected := "ssl.gomaji.com"

	assert.Equal(t, expected, actual)
}

func TestParseRowData(t *testing.T) {
	row := "101.136.192.233 - - [14/Apr/2019:00:03:51 +0800] \"GET /pcode.php?act=check_pcode&plat=mweb&pcode=gm-f6ty88&mobile_phone=&email=&product_id=223956&total_price=888&discount_amt=0 HTTP/1.1\" 200 133 -"

	actual := parseRowData(row)

	expected := accessLog{
		sourceIP: "101.136.192.233",
		method:   "GET",
		route:    "/pcode.php",
		dateTime: "14/Apr/2019:00:03:51 +0800",
		params: map[string]string{
			"act":          "check_pcode",
			"plat":         "mweb",
			"pcode":        "gm-f6ty88",
			"mobile_phone": "",
			"email":        "",
			"product_id":   "223956",
			"total_price":  "888",
			"discount_amt": "0",
		},
	}

	assert.Equal(t, expected, actual)
}

func TestCase1(t *testing.T) {
	row := "101.136.192.233 - - [14/Apr/2019:00:03:51 +0800] \"GET /pcode.php?act=check_pcode&plat=mweb&pcode=gm-f6ty88&mobile_phone=&email=&product_id=223956&total_price=888&discount_amt=0 HTTP/1.1\" 200 133 -"

	parse := parseRowData(row)

	fmt.Printf("params:%v\n", parse.params)
	fmt.Printf("len of params:%v\n", len(parse.params))

	if len(parse.params) > 0 {
		fmt.Printf("success\n")
	}

	a := false
	assert.True(t, a)

}
