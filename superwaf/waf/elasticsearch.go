package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
)

var Ecli *elastic.Client

func InitEs() {
	var err error
	Ecli, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(Esurl))
	if err != nil {
		panic("es client err! error: " + fmt.Sprint(err))
	}
}

func ElasticDo() {
	for {
		select {
		case esinfo := <-EsChan:
			_, err := Ecli.Index().
				Index(EsIndex + "-" + time.Now().Format("200601")).
				BodyJson(esinfo).
				Do(context.Background())
			if err != nil {
				i, _ := json.Marshal(esinfo)
				loginfo := WafLog{DTime: time.Now().Format("2006-01-02 15:04:05"), Client: esinfo.Client, Type: "Elastic", Message: string(i)}
				loginfo.WafDoLog()
			}
		}
	}
}

type WafLog struct {
	DTime   string `json:"dtime"`
	Client  string `json:"client"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (wlog *WafLog) WafDoLog() {
	client := &http.Client{}
	req, _ := json.Marshal(wlog)
	reqinfo := bytes.NewBuffer(req)
	request, _ := http.NewRequest("POST", "http://"+Ip+"/api/v1/waflog", reqinfo)
	_, _ = client.Do(request)
}
