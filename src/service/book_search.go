package service

import (
	"fmt"
	"muyu.com/reader_server_go/v1/src/es"
	"muyu.com/reader_server_go/v1/src/es/model"
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/serializer"
)

type SearchBookOnTypeSvc struct {
	Query string
}

type SearchBookOnTypeRes struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Highlight string `json:"highlight"`
}

func (svc *SearchBookOnTypeSvc) SearchBookOnType() serializer.Response {
	query := svc.Query
	queryTemplate := `
			{
			  "query": {
				"multi_match": {
				  "query": "%s",
				  "type": "bool_prefix",
				  "fields": [
					"name",
					"name._2gram",
					"name._3gram"
				  ]
				}
			  },

			  "highlight": {
				"fields": {
				  "name": {
					"matched_fields": ["name._index_prefix"]
				  }
				}
			  },

			  "size": 10
			}
	`

	queryStr := fmt.Sprintf(queryTemplate, query)

	res, err := es.Search[model.BookEs]("book", queryStr)
	if err != nil {
		logger.Error("search book err: %+v", err)
		return serializer.UnknownErr
	}

	resList := make([]SearchBookOnTypeRes, 0, len(res))

	for _, item := range res {
		resItem := SearchBookOnTypeRes{
			Id:     item.Source.Id,
			Name:   item.Source.Name,
			Author: item.Source.Author,
		}

		if hArr, ok := item.Highlight["name"]; ok && len(hArr) > 0 {
			resItem.Highlight = hArr[0]
		}

		resList = append(resList, resItem)
	}

	return serializer.SuccessRes(resList)
}
