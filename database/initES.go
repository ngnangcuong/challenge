package database

import (
	"github.com/elastic/go-elasticsearch/v8"

	"log"
	"fmt"
	"sync"
)

var (
	lockES = &sync.Mutex{}
	esClient *elasticsearch.Client
	username = "elastic"
	password = "c-awd64Fl6I3BvTv6oIf"
)

func GetESClient() *elasticsearch.Client {
	if esClient == nil {
		lockES.Lock()
		defer lockES.Unlock()

		if esClient == nil {
			cfg := elasticsearch.Config{
				Addresses: []string{"http://localhost:9200" ,"http://es_client:9200", },
				Username: username,
				Password: password,
			}

			es, err := elasticsearch.NewClient(cfg)
			if err != nil {
				log.Fatalln("Fail to get connection to elastic")
			}
			fmt.Println("Connect to elastic")

			return es
		}

		fmt.Println("Already connected to elastic")
	}

	fmt.Println("Already connected to elastic")

	return esClient
}