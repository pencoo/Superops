package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

var Ecli *elastic.Client

func InitEs() {
	var err error
	Ecli, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(E.EsUrl))
	if err != nil {
		panic("es client err! error: " + fmt.Sprint(err))
	}
}

func ElasticDo() {
	for {
		select {
		case esinfo := <-EsChan:
			_, err := Ecli.Index().
				Index(E.EsIndex + "-" + time.Now().Format("200601")).
				BodyJson(esinfo).
				Do(context.Background())
			if err != nil {
				i, _ := json.Marshal(esinfo)
				loginfo := WafLog{DTime: time.Now().Format("2006-01-02 15:04:05"), Client: LockIp, Type: "Elastic", Message: string(i)}
				loginfo.WafDoLog()
			}
		}
	}
}
