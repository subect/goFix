package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"goFix/model"
)

const mapping = `
{
  "mappings": {
    "properties": {
      "user": {
        "type": "keyword"
      },
      "message": {
        "type": "text"
      },
      "image": {
        "type": "keyword"
      },
      "created": {
        "type": "date"
      },
      "tags": {
        "type": "keyword"
      },
      "location": {
        "type": "geo_point"
      },
      "suggest_field": {
        "type": "completion"
      }
    }
  }
}`

func EsTd(c *gin.Context) {
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
