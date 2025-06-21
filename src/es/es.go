package es

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v9"
	"muyu.com/reader_server_go/v1/src/config"
	"muyu.com/reader_server_go/v1/src/logger"
	"net/http"
)

var Client *elasticsearch.Client

func InitEs() {
	con := config.Config.Es

	cfg := elasticsearch.Config{
		Addresses: []string{
			con.Host,
		},
		Username: con.Username,
		Password: con.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	var err error
	Client, err = elasticsearch.NewClient(cfg)

	if Client == nil {
		logger.Panic("es init fail! err = %+v", err)
	}

	if err != nil {
		logger.Panic("es init fail! err = %+v", err)
	}

}
