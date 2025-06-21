package es

import (
	"encoding/json"
	"fmt"
	"muyu.com/reader_server_go/v1/src/es/model"
	"testing"
)

func Test_Es(t *testing.T) {

	query := `
			{
			  "query": {
				"multi_match": {
				  "query": "java go",
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
			  }
			}
	`

	res, err := Search[model.BookEs]("book", query)
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range res {

		b, _ := json.MarshalIndent(item, "", "  ")
		fmt.Println(string(b))

	}

}
