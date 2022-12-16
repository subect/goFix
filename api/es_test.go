package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"goFix/config"
	"goFix/model"
	"testing"
	"time"
)

type Weibo struct {
	User     string                `json:"user"`               // 用户
	Message  string                `json:"message"`            // 微博内容
	Retweets int                   `json:"retweets"`           // 转发数
	Image    string                `json:"image,omitempty"`    // 图片
	Created  time.Time             `json:"created,omitempty"`  // 创建时间
	Tags     []string              `json:"tags,omitempty"`     // 标签
	Location string                `json:"location,omitempty"` //位置
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

// 创建索引
func TestCreatIndex(t *testing.T) {
	config.InitConfig()
	ctx := context.Background()
	esClient := model.GetEsOpt()
	exist, err := esClient.IndexExists("es_hx").Do(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !exist {
		_, err := esClient.CreateIndex("es_hx").BodyString(mapping).Do(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("es_hx exists")
	}
}

//插入文档
func TestInsertDoc(t *testing.T) {
	config.InitConfig()
	ctx := context.Background()
	esClient := model.GetEsOpt()
	weibo := Weibo{
		User:     "lxl",
		Message:  "明天是星期天",
		Retweets: 2,
		Image:    "http://localhost:8080/images/678.jpg",
		Created:  time.Now(),
		Tags:     []string{"开心", "茄子"},
	}
	put, err := esClient.Index().Index("es_hx").Id("2").BodyJson(weibo).Do(ctx)
	fmt.Println(put, err)
}

// 获取文档
func TestGetDoc(t *testing.T) {
	config.InitConfig()
	ctx := context.Background()
	esClient := model.GetEsOpt()
	get, err := esClient.Get().Index("es_hx").Id("1").Do(ctx)
	if err != nil {
		panic(err)
	}
	if get.Found {
		fmt.Printf("文档id:%s -->版本号:%d -->索引名:%s\n", get.Id, get.Version, get.Index)
	}
	msg := Weibo{}
	dataJson, _ := get.Source.MarshalJSON()
	json.Unmarshal(dataJson, &msg)
	fmt.Printf("msg:%+v\n", msg)
}

//批量获取文档
func TestMultiGet(t *testing.T) {
	config.InitConfig()
	ctx := context.Background()
	esClient := model.GetEsOpt()
	result, err := esClient.MultiGet().Add(elastic.NewMultiGetItem().Index("es_hx").Id("1")).
		Add(elastic.NewMultiGetItem().Index("es_hx").Id("2")).
		Do(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, doc := range result.Docs {
		weibo := Weibo{}
		data, err := doc.Source.MarshalJSON()
		if err != nil {
			fmt.Println(err.Error())
		}
		json.Unmarshal(data, &weibo)
		fmt.Printf("weibo:%+v\n", weibo)
	}
}

//更新文档
func TestUpdateDoc(t *testing.T) {
	//config.InitConfig()
	//ctx := context.Background()
	//esClient := model.GetEsOpt()
}
