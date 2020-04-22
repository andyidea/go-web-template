package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type TiandituGeocoderData struct {
	Status        string `json:"status"`
	Msg           string `json:"msg"`
	SearchVersion string `json:"searchVersion"`
	Location      struct {
		Lon   float64 `json:"lon"`
		Level string  `json:"level"`
		Lat   float64 `json:"lat"`
	} `json:"location"`
}

func TiandituGeocoder(keyWord string) (*TiandituGeocoderData, error) {
	if keyWord == "" {
		return &TiandituGeocoderData{}, nil
	}
	url := "http://api.tianditu.gov.cn/geocoder?ds={\"keyWord\":\"" + keyWord + "\"}&tk=9f788b93b7b6787ab0f08523234785d7"

	log.Debug().Msg(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Info().Msg(string(body))

	var data TiandituGeocoderData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if data.Status != "0" {
		return nil, errors.New("天地图接口报错：" + data.Status + data.Msg)
	}

	return &data, nil
}

type TencentGeocoderData struct {
	Status uint64                `json:"status"`
	Msg    string                `json:"message"`
	Count  uint64                `json:"count"`
	Data   []TencentGeocoderItem `json:"data"`
}

type TencentGeocoderItem struct {
	Title    string `json:"title"`
	Location struct {
		Lng float64 `json:"lng"`
		Lat float64 `json:"lat"`
	}
	AdInfo struct {
		Adcode   uint64 `json:"adcode"`
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
	} `json:"ad_info"`
	Address string `json:"address"`
}

func TencentGeocoder(address, city string) (*TencentGeocoderData, error) {
	if address == "" {
		return &TencentGeocoderData{}, nil
	}
	url := "https://apis.map.qq.com/ws/place/v1/search?boundary=region(" + city + ",0)&page_size=1&page_index=1&keyword=" + address + "&orderby=_distance&key=P3WBZ-Z2KRW-RDYR7-RMQNZ-7OMZ6-HMBWZ"
	//url := "https://apis.map.qq.com/ws/geocoder/v1/?city=上海&address=" + address + "&key=P3WBZ-Z2KRW-RDYR7-RMQNZ-7OMZ6-HMBWZ"

	log.Debug().Msg(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Info().Msg(string(body))

	var data TencentGeocoderData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if data.Status != 0 {
		return nil, errors.New("腾讯地图接口报错：" + string(data.Status) + data.Msg)
	}

	return &data, nil
}

func main() {

}
