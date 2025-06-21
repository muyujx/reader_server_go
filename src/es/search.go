package es

import (
	"muyu.com/reader_server_go/v1/src/logger"
	"muyu.com/reader_server_go/v1/src/util"
	"strings"
)

type SearchResult[T any] struct {
	Source    *T                  `json:"_source"`
	Highlight map[string][]string `json:"highlight"`
}

// Search 封装 es 搜索
func Search[T any](index string, queryStr string) ([]SearchResult[T], error) {
	res, err := Client.Search(
		Client.Search.WithIndex(index),
		Client.Search.WithBody(strings.NewReader(queryStr)),
	)
	if err != nil {
		return nil, err
	}
	resList, err := ParseEsResponse(res.Body)
	if err != nil {
		return nil, err
	}
	searchRes := make([]SearchResult[T], 0, len(resList))

	for _, item := range resList {
		resItem := SearchResult[T]{}

		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		source, ok := itemMap["_source"].(map[string]any)
		if ok {
			resSource := new(T)
			err = util.MapToStruct(source, resSource)
			if err == nil {
				resItem.Source = resSource
			} else {
				logger.Warning("[Search] parse source fail! source = %+v", source)
			}
		}

		if highlight, ok := itemMap["highlight"].(map[string]any); ok {
			resItem.Highlight = make(map[string][]string, len(highlight))

			for k, v := range highlight {
				if vList, ok := v.([]any); ok {
					resItem.Highlight[k] = make([]string, 0, len(vList))
					for _, v := range vList {
						if vStr, ok := v.(string); ok {
							resItem.Highlight[k] = append(resItem.Highlight[k], vStr)
						}
					}
				}
			}

		}

		searchRes = append(searchRes, resItem)
	}

	return searchRes, nil
}
