package client

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"os"
	"time"
)

type Client struct {
	URL string
}

func NewClient(url string) *Client {
	return &Client{
		URL: url,
	}
}

func (client *Client) NewElasticClient() (*elastic.Client, error) {
	elasticClient, err := elastic.NewClient(
		elastic.SetURL(client.URL),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeout(60*time.Second),
		elastic.SetHealthcheckInterval(60*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		return nil, err
	}
	return elasticClient, nil
}

func (client *Client) Ping() (*elastic.PingResult, int, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, 0, err
	}
	return elasticClient.Ping(client.URL).Do(context.Background())
}

func (client *Client) Version() (string, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return "", err
	}
	return elasticClient.ElasticsearchVersion(client.URL)
}

func (client *Client) ClusterStats() (*elastic.ClusterStatsResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.ClusterStats().Do(context.Background())
}

func (client *Client) ClusterState() (*elastic.ClusterStateResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.ClusterState().Do(context.Background())
}

func (client *Client) CatHealth() (elastic.CatHealthResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.CatHealth().Do(context.Background())
}

func (client *Client) CatAliases() (elastic.CatAliasesResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.CatAliases().Do(context.Background())
}

func (client *Client) CatAllocation() (elastic.CatAllocationResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.CatAllocation().Do(context.Background())
}

func (client *Client) CatCount() (elastic.CatCountResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.CatCount().Do(context.Background())
}

func (client *Client) CatIndices() (elastic.CatIndicesResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.CatIndices().Do(context.Background())
}

func (client *Client) IndexExists(index string) (bool, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return false, err
	}
	exists, err := elasticClient.IndexExists(index).Do(context.Background())
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (client *Client) CreateIndex(index string) (bool, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return false, err
	}
	indicesCreateResult, err := elasticClient.CreateIndex(index).Do(context.Background())
	if err != nil {
		return false, err
	}
	if !indicesCreateResult.Acknowledged {
		return false, fmt.Errorf("create index[%v] error. Acknowledged[%v]", indicesCreateResult.Index, indicesCreateResult.Acknowledged)
	}
	return true, nil
}

func (client *Client) DeleteIndex(index string) (bool, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return false, err
	}
	indicesDeleteResponse, err := elasticClient.DeleteIndex(index).Do(context.Background())
	if err != nil {
		return false, err
	}

	elasticClient.ClusterStats()
	if !indicesDeleteResponse.Acknowledged {
		return false, fmt.Errorf("delete index[%v] error. Acknowledged[%v]", index, indicesDeleteResponse.Acknowledged)
	}
	return true, nil
}

func (client *Client) MultiGet(multiGetArray []MultiGetDTO) (*elastic.MgetResponse, error) {
	if len(multiGetArray) == 0 {
		return nil, fmt.Errorf("bulk data length[%v] cannot be zero", len(multiGetArray))
	}
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	mgetService := elasticClient.MultiGet()
	for _, multiGet := range multiGetArray {
		mgetService = mgetService.Add(elastic.NewMultiGetItem().Index(multiGet.Index).Type(multiGet.Type).Id(multiGet.Id))
	}
	return mgetService.Do(context.Background())
}

func (client *Client) Bulk(bulkArray []BulkDTO) (*elastic.BulkResponse, error) {
	if len(bulkArray) == 0 {
		return nil, fmt.Errorf("bulk data length[%v] cannot be zero", len(bulkArray))
	}
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	mgetResponse, err := client.MultiGet(batchBulkToMultiGet(bulkArray))
	if err != nil {
		return nil, err
	}
	bulk := elasticClient.Bulk()
	for _, bulkDTO := range bulkArray {
		if checkExists(mgetResponse.Docs, bulkDTO.Id) {
			if bulkDTO.IsDelete {
				request := elastic.NewBulkDeleteRequest().Index(bulkDTO.Index).Type(bulkDTO.Type).Id(bulkDTO.Id)
				bulk.Add(request)
			} else {
				request := elastic.NewBulkUpdateRequest().Index(bulkDTO.Index).Type(bulkDTO.Type).Id(bulkDTO.Id).Doc(bulkDTO.Data)
				bulk.Add(request)
			}
		} else {
			if !bulkDTO.IsDelete {
				if bulkDTO.Id == "" {
					bulkDTO.Id = getUUID()
				}
				request := elastic.NewBulkIndexRequest().Index(bulkDTO.Index).Type(bulkDTO.Type).Id(bulkDTO.Id).Doc(bulkDTO.Data)
				bulk.Add(request)
			}
		}
	}
	return bulk.Do(context.Background())
}

func (client *Client) GetMappings() (map[string]interface{}, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.GetMapping().Index().Do(context.Background())
}

func (client *Client) GetMapping(index string) (map[string]interface{}, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.GetMapping().Index(index).Do(context.Background())
}
func (client *Client) Update(index string, id string, doc interface{}) (*elastic.UpdateResponse, error) {
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return nil, err
	}
	return elasticClient.Update().Index(index).Id(id).Upsert(doc).Do(context.Background())
}

func (client *Client) Delete(index string, id string) (int64, error) {
	termQuery := elastic.NewTermQuery("id", id)
	elasticClient, err := client.NewElasticClient()
	if err != nil {
		return 0, err
	}
	deleteResponse, err := elasticClient.DeleteByQuery(index).Query(termQuery).Do(context.Background())
	if err != nil {
		return 0, err
	}
	return deleteResponse.Deleted, nil
}
