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
		log.Printf("[InitElasticsearch] Error creating client: %s", err)
		return nil, nil // Return nil client, but no fatal error to allow app to start
	}

	// Verify connection
	res, err := es.Info()
	if err != nil {
		log.Printf("[InitElasticsearch] Connection failed (running without ES): %s", err)
		return nil, nil
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[InitElasticsearch] Connection Error: %s", res.String())
		return nil, nil
	}

	log.Println("[InitElasticsearch] Connected successfully")
	return es, nil
}
