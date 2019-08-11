package service

import (
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sillyhatxu/elasticsearch-ui/client"
	"github.com/sillyhatxu/elasticsearch-ui/dto"
)

func Ping(url string) (*elastic.PingResult, int, error) {
	esClient := client.NewClient(url)
	return esClient.Ping()
}

func Version(url string) (string, error) {
	esClient := client.NewClient(url)
	return esClient.Version()
}

func Health(url string) (elastic.CatHealthResponse, error) {
	esClient := client.NewClient(url)
	return esClient.CatHealth()
}

func ClusterStats(url string) (*elastic.ClusterStatsResponse, error) {
	esClient := client.NewClient(url)
	return esClient.ClusterStats()
}

func Indices(url string) (elastic.CatIndicesResponse, error) {
	esClient := client.NewClient(url)
	return esClient.CatIndices()
}

func GetMappings(url string) ([]dto.MappingsDTO, error) {
	esClient := client.NewClient(url)
	mappingsMap, err := esClient.GetMappings()
	if err != nil {
		return nil, err
	}
	var mappings []dto.MappingsDTO
	for k, v := range mappingsMap {
		mappings = append(mappings, dto.MappingsDTO{
			Index:    k,
			Mappings: getMappings(v),
		})
	}
	return mappings, nil
}

func getMappings(v interface{}) []dto.MappingsDetailDTO {
	value, ok := v.(map[string]interface{})
	var mappingDetailArray []dto.MappingsDetailDTO
	if ok {
		mappingsValue, ok := value["mappings"].(map[string]interface{})
		if ok {
			propertiesValue, ok := mappingsValue["properties"].(map[string]interface{})
			if ok {
				for k, v := range propertiesValue {
					typeSrc, isField := "", true
					typeValue, ok := v.(map[string]interface{})
					if ok {
						if typeValue["properties"] != nil {
							isField = false
							for _, m := range getMappingsObj(k, typeValue["properties"]) {
								mappingDetailArray = append(mappingDetailArray, dto.MappingsDetailDTO{
									Field: m.Field,
									Type:  m.Type,
								})
							}
						} else {
							typeSrc = fmt.Sprintf("%v", typeValue["type"])
						}
					}
					if isField {
						mappingDetailArray = append(mappingDetailArray, dto.MappingsDetailDTO{
							Field: k,
							Type:  typeSrc,
						})
					}
				}
			}
		}
	}
	if mappingDetailArray == nil {
		mappingDetailArray = make([]dto.MappingsDetailDTO, 0)
	}
	return mappingDetailArray
}

func getMappingsObj(previouKsey string, value interface{}) []dto.MappingsDetailDTO {
	var mappingsDetailArray []dto.MappingsDetailDTO
	propertiesValue, ok := value.(map[string]interface{})
	if ok {
		for k, v := range propertiesValue {
			typeSrc, isField := "", true
			typeValue, ok := v.(map[string]interface{})
			if ok {
				if typeValue["properties"] != nil {
					isField = false
					for _, m := range getMappingsObj(k, typeValue["properties"]) {
						mappingsDetailArray = append(mappingsDetailArray, dto.MappingsDetailDTO{
							Field: m.Field,
							Type:  m.Type,
						})
					}
				} else {
					typeSrc = fmt.Sprintf("%v", typeValue["type"])
				}
			}
			if isField {
				mappingsDetailArray = append(mappingsDetailArray, dto.MappingsDetailDTO{
					Field: fmt.Sprintf("%s.%s", previouKsey, k),
					Type:  typeSrc,
				})
			}
		}
	}
	return mappingsDetailArray
}
