package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

func InitElasticsearch(cfg ElasticSearch) (*elasticsearch.Client, error) {
	configElastic := elasticsearch.Config{
		Addresses: []string{cfg.Host},
	}
	es, err := elasticsearch.NewClient(configElastic)
	if err != nil {
		log.Fatalf("[InitElasticsearch] Error initializing Elasticsearch: %s", err)
		return nil, err
	}

	return es, nil
}
