package client

import (
	"github.com/olivere/elastic"
	"github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

func batchBulkToMultiGet(bulkArray []BulkDTO) []MultiGetDTO {
	var array []MultiGetDTO
	for _, bulk := range bulkArray {
		array = append(array, MultiGetDTO{Id: bulk.Id, Index: bulk.Index, Type: bulk.Type})
	}
	return array
}

func checkExists(resultArray []*elastic.GetResult, id string) bool {
	for _, result := range resultArray {
		if result.Id == id {
			return result.Found
		}
	}
	return false
}

var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getUUID() string {
	var id string
	idUUID, err := uuid.NewV4()
	if err != nil {
		id = randStringRunes(32)
	} else {
		id = idUUID.String()
	}
	return strings.ToUpper(strings.ReplaceAll(id, "-", ""))
}
