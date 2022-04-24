package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type dailyThreePostMessage struct {
	Sfzx       int    `json:"sfzx"` //是否在校
	GeoApiInfo string `json:"geo_api_info"`
	Address    string `json:"address"`
	Area       string `json:"area"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Tw         int    `json:"tw"`      // 体温
	Sfyzz      int    `json:"sfyzz"`   // 是否有症状
	Sfcyglq    int    `json:"sfcyglq"` // 是否处于隔离期
	Ymtys      int    `json:"ymtys"`   // 一码通颜色
	Qtqk       string `json:"qtqk"`    // 其他信息
}

func DailyThreeHandler() {
	msg := AssembleDailyThreePostMessage()
	req, err := http.NewRequest("POST", DAILY_UP_URL, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if res.StatusCode != 200 {
		log.Println("Failed Status code of response is not 200")
		data, _ := io.ReadAll(res.Body)
		log.Fatalf(string(data))
	}
}

func AssembleDailyThreePostMessage() []byte {
	type interJSON struct {
		FormattedAddress string `json:"formattedAddress"`
		AddressComponent struct {
			Country      string        `json:"country"`
			Province     string        `json:"province"`
			Citycode     string        `json:"citycode"`
			City         string        `json:"city"`
			District     string        `json:"district"`
			Township     []interface{} `json:"township"`
			Neighborhood struct {
				Name []interface{} `json:"name"`
				Type []interface{} `json:"type"`
			} `json:"neighborhood"`
			Building struct {
				Name []interface{} `json:"name"`
				Type []interface{} `json:"type"`
			} `json:"building"`
			Adcode string `json:"adcode"`
			Street string `json:"street"`
		} `json:"addressComponent"`
		Number   []interface{} `json:"number"`
		Location string        `json:"location"`
		Level    string        `json:"level"`
	}
	var rawMsg dailyThreePostMessage
	// 其他信息:"qtqk": ""
	// 体温:"tw": 0,
	// 是否有症状:"sfyzz": 0,
	// 是否处于隔离期:"sfcyglq": 0,
	// 一码通颜色:"ymtys": 0,
	// 是否在校   "sfzx": 1,
	rawMsg.Sfyzz = 0
	if inSchool {
		rawMsg.Sfzx = 1
	} else {
		rawMsg.Sfzx = 0
	}
	rawMsg.Tw = 0
	rawMsg.Sfcyglq = 0
	rawMsg.Ymtys = 0
	rawMsg.Qtqk = ""
	// Geometrical Information
	rawMsg.Province = locationGEO.Geocodes[0].Province
	rawMsg.City = locationGEO.Geocodes[0].City
	rawMsg.Address = locationGEO.Geocodes[0].FormattedAddress
	rawMsg.Area = locationGEO.Geocodes[0].Province + " " + locationGEO.Geocodes[0].City + " " + locationGEO.Geocodes[0].District
	// GeoAPIInfo
	tempJSON := interJSON{
		FormattedAddress: locationGEO.Geocodes[0].FormattedAddress,
		AddressComponent: struct {
			Country      string        `json:"country"`
			Province     string        `json:"province"`
			Citycode     string        `json:"citycode"`
			City         string        `json:"city"`
			District     string        `json:"district"`
			Township     []interface{} `json:"township"`
			Neighborhood struct {
				Name []interface{} `json:"name"`
				Type []interface{} `json:"type"`
			} `json:"neighborhood"`
			Building struct {
				Name []interface{} `json:"name"`
				Type []interface{} `json:"type"`
			} `json:"building"`
			Adcode string `json:"adcode"`
			Street string `json:"street"`
		}{
			Country:  locationGEO.Geocodes[0].Country,
			Province: locationGEO.Geocodes[0].Province,
			Citycode: locationGEO.Geocodes[0].Citycode,
			City:     locationGEO.Geocodes[0].City,
			District: locationGEO.Geocodes[0].District,
			Township: locationGEO.Geocodes[0].Township,
			Adcode:   locationGEO.Geocodes[0].Adcode,
			Street:   locationGEO.Geocodes[0].Street,
		},
		Location: locationGEO.Geocodes[0].Location,
		Level:    locationGEO.Geocodes[0].Level,
	}
	tempData, err := json.Marshal(tempJSON)
	if err != nil {
		log.Fatalf(err.Error())
	}
	rawMsg.GeoApiInfo = string(tempData[:])
	var msgJSON []byte
	msgJSON, err = json.Marshal(rawMsg)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return msgJSON
}
