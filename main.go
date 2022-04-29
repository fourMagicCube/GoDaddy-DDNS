package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var domain = "" //域名：不包括www.
var key = ""    //自己从godaddy api创建的key
var secret = "" //自己从godaddy api创建的secret
var logPath = "logs/"

var name = "AAAA" //ipv4:A; ipv6:AAAA
var domainType = "@"
var accept = "application/json"
var ip6Url = "https://api6.ipify.org/" //ipv4:https://api4.ipify.org/
var apiUrl = "https://api.godaddy.com/v1/domains/" + domain + "/records/" + name + "/" + domainType
var authorization = "sso-key " + key + ":" + secret
var tempDate string

func main() {
	for {
		after := time.After(time.Hour)
		go ddns()
		<-after
	}
}

func ddns() {

	//定义日志
	nowDate := time.Now().Format("2006-01-02")
	if tempDate != nowDate {
		tempDate = nowDate
	}
	fileName := logPath + "info-" + tempDate + ".log"

	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	logInfo := log.New(logFile, "[Info]", log.Ldate|log.Lmicroseconds|log.Llongfile)
	logError := log.New(logFile, "[Error]", log.Ldate|log.Lmicroseconds|log.Llongfile)

	//预热
	apiGetReq, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		logError.Println(err)
		return
	}
	apiGetReq.Header.Add("Authorization", authorization)
	apiGetReq.Header.Add("Accept", accept)

	//获取当前v6地址
	getIp6Url, err := http.Get(ip6Url)
	body, err := ioutil.ReadAll(getIp6Url.Body)
	if err != nil {
		logError.Println(err)
		return
	}
	var currentIp *string
	temp := string(body)
	currentIp = &temp
	if *currentIp == "" {
		logError.Println("请求出错： currentIp为空")
		return
	}

	//获取域名服务器上的ip地址
	client := http.DefaultClient
	apiGetResponse, err := client.Do(apiGetReq)
	body, err = ioutil.ReadAll(apiGetResponse.Body)
	if err != nil {
		logError.Println(err)
		return
	}
	var oldResponses []responseBody
	err = json.Unmarshal(body, &oldResponses)
	if err != nil {
		logError.Println(err)
		return
	}
	oldIp := &oldResponses[0].Data

	//对比两个ip,不相同的话就注册新的ip
	if *currentIp != *oldIp {
		logInfo.Println(" ip不相同，进行更新ddns更新ip")
		oldResponses[0].Data = *currentIp

		marshal, err := json.Marshal(oldResponses)
		if err != nil {
			logError.Println(err, "转换失败")
			return
		}
		bodyReader := bytes.NewReader(marshal)

		apiPutReq, err := http.NewRequest("PUT", apiUrl, bodyReader)
		if err != nil {
			logError.Println(err)
			return
		}
		apiPutReq.Header.Add("Authorization", authorization)
		apiPutReq.Header.Add("Content-Type", accept)

		apiPutResponse, err := client.Do(apiPutReq)
		if err != nil {
			logError.Println(err, "更新失败")
			return
		}
		if apiPutResponse.StatusCode == 200 {
			logInfo.Println(err, "DDNS更新成功")
		}
	} else {
		logInfo.Println(" 检测完毕，无需更改")
	}
}

type responseBody struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Ttl  int64  `json:"ttl"`
	Type string `json:"type"`
}
