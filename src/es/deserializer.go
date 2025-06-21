package es

import (
	"encoding/json"
	"errors"
	"io"
	"muyu.com/reader_server_go/v1/src/logger"
)

func ParseEsResponse(body io.ReadCloser) (res []any, err error) {

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			logger.Error("[ParseEsResponse] body.Close() err=%+v", err)
		}
	}(body)

	var response map[string]any
	if err = json.NewDecoder(body).Decode(&response); err != nil {
		return
	}

	hits, ok := response["hits"].(map[string]any)
	if !ok {
		logger.Warning("[ParseEsResponse] body.hits is not map[string]any, result is %+v", response)
		return nil, errors.New("not found hits")
	}

	res, ok = hits["hits"].([]any)
	if !ok {
		logger.Warning("[ParseEsResponse] body.hits is not []any, result is %+v", hits)
		return nil, errors.New("not found inner hits")
	}

	return res, nil
}
