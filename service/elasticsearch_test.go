package service

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var url = "http://localhost:9200"

func TestGetMappings(t *testing.T) {
	resp, err := GetMappings(url)
	assert.Nil(t, err)
	respJSON, err := json.Marshal(resp)
	assert.Nil(t, err)
	fmt.Println(string(respJSON))
}
