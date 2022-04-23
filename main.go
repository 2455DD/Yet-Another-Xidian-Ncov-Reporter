package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DAILY_UP_URL    = "https://xxcapp.xidian.edu.cn/site/ncov/xidiandailyup"
	HEALTH_CARD_URL = "https://xxcapp.xidian.edu.cn/ncov/wap/default/index"
)

// Parameter
var (
	username string
	password string
	location string
	option   string
	mode     string //控制填报方式
	apiKey   string
	client   = http.Client{
		Timeout: time.Second * 15, // Maximum of 2 secs
	}
)

//
var (
	locationGEO geoLocation
)

var locationMap = map[string]string{
	"xian_south": "xian_south.json",
	"xian_north": "xian_north.json",
	"guangzhou":  "guangzhou.json",
}

type geoLocation struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Geocodes []struct {
		FormattedAddress string        `json:"formatted_address"`
		Country          string        `json:"country"`
		Province         string        `json:"province"`
		Citycode         string        `json:"citycode"`
		City             string        `json:"city"`
		District         string        `json:"district"`
		Township         []interface{} `json:"township"`
		Adcode           string        `json:"adcode"`
		Street           string        `json:"street"`
		Number           string        `json:"number"`
		Location         string        `json:"location"`
		Level            string        `json:"level"`
	} `json:"geocodes"`
}

// JSON MAP
var ()

func Authentication() {
	type responseJSON struct {
		E int    `json:"e"`
		M string `json:"m"`
		D struct {
		} `json:"d"`
	}
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	resp, err := client.Post("https://xxcapp.xidian.edu.cn/uc/wap/login/check", "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalf("Authentication Failed!\n%v", err.Error())
	}
	respJSON := &responseJSON{}
	json.NewDecoder(resp.Body).Decode(respJSON)
	if resp.StatusCode != 200 {
		log.Fatalf("Authentication Failed! %v\n", respJSON.M)
	}
	if respJSON.E != 0 {
		log.Fatalf("Authentication Failed! %v\n", respJSON.M)
	}
	log.Printf("Authentication Succeed! %v\n", respJSON.M)
}

func getJson(myClient *http.Client, url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		log.Printf("HTTP Request from %v Failure:%v\n", url, err.Error())
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		log.Printf("ERR: Get non-json reply from %v\nStatusCode:%v Body:%v", url, r.StatusCode, bodyString)
	}
	return err
}

func ReadFromCmd() {
	flag.StringVar(&username, "u", "", "学号,不得为空")
	flag.StringVar(&password, "p", "", "密码,不得为空")
	flag.StringVar(&location, "location", "xian_south", "校区,默认为西电南校区(长安校区)")
	flag.StringVar(&option, "option", "", "其他地区,仅在location==others时有效")
	flag.StringVar(&apiKey, "key", "", "API_KEY")
	flag.StringVar(&mode, "mode", "daily3", "模式: 晨午晚检为daily3, 健康卡为 hc")
	flag.Parse()
	if username == "" || password == "" || (mode != "daily3" && mode != "hc") ||
		(location != "xian_south" && location != "xian_north" && location != "others" && location != "guangzhou") {
		flag.Usage()
		os.Exit(-1)
	}
	if location == "others" {
		if option == "" || apiKey == "" {
			flag.Usage()
			os.Exit(-1)
		}
		err := getJson(&client,
			fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%v&address=%v", apiKey, option),
			&locationGEO)
		log.Println(locationGEO)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		data, err := os.ReadFile(locationMap[location])
		if err != nil {
			log.Fatalf(err.Error())
		}
		err = json.Unmarshal(data, &locationGEO)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

}
func main() {
	ReadFromCmd()
	//Authentication()
	defer client.CloseIdleConnections()

}
