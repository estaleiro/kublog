package outputs

import (
	"context"
	"log"
	"time"

	"kublog/config"

	elastic "gopkg.in/olivere/elastic.v5"
	"k8s.io/api/core/v1"
)

// Kublog is a structure used for serializing/deserializing data in Elasticsearch
type Kublog struct {
	Operation  string    `json:"operation"`
	Namespace  string    `json:"namespace"`
	Name       string    `json:"name"`
	ActionTime time.Time `json:"action_time"`
	PodIP      string    `json:"pod_ip"`
	HostIP     string    `json:"host_ip"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"kublog":{
			"properties":{
				"operation":{
					"type":"keyword"
				},
				"namespace":{
					"type":"keyword"
				},
				"name":{
					"type":"keyword"
				},
				"action_time":{
					"type":"date"
				},
				"pod_ip":{
					"type":"keyword"
				},
				"host_ip":{
					"type":"keyword"
				}
			}
		}
	}
}`

// IndexKublog index logs in ElasticSearch
func IndexKublog(pod *v1.Pod, operation string, config config.BaseConfig) {

	kublog := Kublog{
		Operation:  operation,
		Namespace:  pod.ObjectMeta.Namespace,
		Name:       pod.ObjectMeta.Name,
		ActionTime: pod.ObjectMeta.CreationTimestamp.Time,
		PodIP:      pod.Status.PodIP,
		HostIP:     pod.Status.HostIP}

	// Starting with elastic.v5, you must pass a context to execute each service
	d, _ := time.ParseDuration(config.ElasticSearchConfig.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	var clientOptions []elastic.ClientOptionFunc

	i, _ := time.ParseDuration(config.ElasticSearchConfig.HealthCheckInterval)
	clientOptions = append(clientOptions,
		elastic.SetSniff(config.ElasticSearchConfig.EnableSniffer),
		elastic.SetURL(config.ElasticSearchConfig.Hosts...),
		elastic.SetHealthcheckInterval(i),
	)

	if i == 0 {
		clientOptions = append(clientOptions,
			elastic.SetHealthcheck(false),
		)
	}

	client, err := elastic.NewClient(clientOptions...)
	if err != nil {
		// Handle error
		log.Println(err)
	}

	indexName := config.ElasticSearchConfig.Indexname

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		// Handle error
		log.Println(err)
	}
	if !exists {
		// Create a new index.
		_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			log.Println(err)
		}
	}

	_, err = client.Index().
		Index(indexName).
		Type("group").
		BodyJson(kublog).
		Do(ctx)
	if err != nil {
		// Handle error
		log.Printf("Operation %s was not indexed due to error: %v \n", operation, err)
	}

	log.Printf("Indexed operation %s to index %s", kublog.Operation, indexName)
}
