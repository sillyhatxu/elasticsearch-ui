package client

import (
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var url = "http://localhost:9200"

func TestClient_Version(t *testing.T) {
	client := NewClient(url)
	version, err := client.Version()
	assert.Nil(t, err)
	fmt.Println(version)
}

func TestClient_Ping(t *testing.T) {
	client := NewClient(url)
	ping, code, err := client.Ping()
	assert.Nil(t, err)
	fmt.Println(code)
	pingJSON, err := json.Marshal(ping)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Ping"))
	fmt.Println(string(pingJSON))
}

//fielddata memory_size_in_bytes
// query_cache
func TestClient_ClusterState(t *testing.T) {
	client := NewClient(url)
	state, err := client.ClusterStats()
	assert.Nil(t, err)
	stateJSON, err := json.Marshal(state)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "ClusterState"))
	fmt.Println(string(stateJSON))
}

//CatHealth "status":"yellow"  total shards: "shards":"3"
//total size : store.size
//CatIndices index count:count
//shards
func TestClient_Cat(t *testing.T) {
	client := NewClient(url)
	response, err := client.CatHealth()
	assert.Nil(t, err)
	responseJSON, err := json.Marshal(response)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "CatHealth"))
	fmt.Println(string(responseJSON))

	responseAliases, err := client.CatAliases()
	assert.Nil(t, err)
	responseAliasesJSON, err := json.Marshal(responseAliases)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "CatAliases"))
	fmt.Println(string(responseAliasesJSON))

	responseAllocation, err := client.CatAllocation()
	assert.Nil(t, err)
	responseAllocationJSON, err := json.Marshal(responseAllocation)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "CatAllocation"))
	fmt.Println(string(responseAllocationJSON))

	responseCount, err := client.CatCount()
	assert.Nil(t, err)
	responseCountJSON, err := json.Marshal(responseCount)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "CatCount"))
	fmt.Println(string(responseCountJSON))

	responseIndices, err := client.CatIndices()
	assert.Nil(t, err)
	responseIndicesJSON, err := json.Marshal(responseIndices)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "CatIndices"))
	fmt.Println(string(responseIndicesJSON))
}

func TestClient_CreateIndex(t *testing.T) {
	client := NewClient(url)
	b, err := client.CreateIndex("demo")
	assert.Nil(t, err)
	assert.EqualValues(t, b, true)
}

func TestClient_IndexExists(t *testing.T) {
	client := NewClient(url)
	b, err := client.IndexExists("demo")
	assert.Nil(t, err)
	assert.EqualValues(t, b, true)
}

func TestClient_DeleteIndex(t *testing.T) {
	client := NewClient(url)
	b, err := client.DeleteIndex("user")
	assert.Nil(t, err)
	assert.EqualValues(t, b, true)
}

func TestClient_GetMapping(t *testing.T) {
	client := NewClient(url)
	resp, err := client.GetMapping("user")
	assert.Nil(t, err)
	respJSON, err := json.Marshal(resp)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Mapping"))
	fmt.Println(string(respJSON))
}

func TestClient_BulkProduct(t *testing.T) {
	type Product struct {
		Id               string    `json:"id"`
		Name             string    `json:"name"`
		Status           string    `json:"status"`
		Price            int64     `json:"price"`
		Recommend        float64   `json:"recommend"`
		InStock          bool      `json:"in_stock"`
		Inventory        int       `json:"inventory"`
		CreatedTime      time.Time `json:"created_time"`
		LastModifiedTime time.Time `json:"last_modified_time"`
	}
	var productArray []Product
	count := 5562148
	for i := 700000; i < count; i++ {
		productArray = append(productArray, Product{
			Id:               getUUID(),
			Name:             fmt.Sprintf("%s-%d", randomdata.LastName(), i),
			Status:           "ENABLE",
			Price:            int64(rand.Int63n(800000-123) + 123),
			Recommend:        rand.Float64(),
			InStock:          randomdata.Boolean(),
			Inventory:        randomdata.Number(80),
			CreatedTime:      time.Now(),
			LastModifiedTime: time.Now(),
		})
	}
	client := NewClient(url)
	var bulkArray []BulkDTO
	for _, product := range productArray {
		bulkArray = append(bulkArray, BulkDTO{
			Index:    "product",
			Type:     "tags",
			Data:     product,
			IsDelete: false,
		})
	}
	pageSize := 50000
	pageSum := (count + pageSize - 1) / pageSize
	for i := 0; i < pageSum; i++ {
		resp, err := client.Bulk(bulkArray[i*pageSize : (i+1)*pageSize])
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("---------- %s ----------", "Failed"))
		for _, res := range resp.Failed() {
			resJSON, err := json.Marshal(res)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(resJSON))
		}
	}

}

func TestClient_BulkUser(t *testing.T) {
	type Department struct {
		Id               int64     `json:"id"`
		Name             string    `json:"name"`
		Status           string    `json:"status"`
		CreatedTime      time.Time `json:"created_time"`
		LastModifiedTime time.Time `json:"last_modified_time"`
	}
	type User struct {
		Id               int64       `json:"id"`
		UserName         string      `json:"user_name"`
		Status           string      `json:"status"`
		Dept             *Department `json:"dept"`
		Age              int         `json:"age"`
		IsDelete         bool        `json:"is_delete"`
		Address          string      `json:"address"`
		Country          string      `json:"country"`
		Email            string      `json:"email"`
		UserType         int         `json:"user_type"`
		CreatedTime      time.Time   `json:"created_time"`
		LastModifiedTime time.Time   `json:"last_modified_time"`
	}
	var deptArray []Department
	for i := 1; i <= 10; i++ {
		dept := Department{Id: int64(i), Name: fmt.Sprintf("DEPT-%d", i), Status: "ENABLE", CreatedTime: time.Now(), LastModifiedTime: time.Now()}
		if i == 3 || i == 5 {
			dept.Status = "DISABLE"
		}
		deptArray = append(deptArray, dept)
	}
	var userArray []User
	for i := 1; i <= 50000; i++ {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		user := User{
			Id:               int64(i),
			UserName:         randomdata.SillyName(),
			Status:           "ENABLE",
			Dept:             &deptArray[r.Intn(len(deptArray))],
			Age:              int(rand.Int63n(45-15) + 15),
			IsDelete:         randomdata.Boolean(),
			Country:          randomdata.Country(randomdata.FullCountry),
			Email:            randomdata.Email(),
			Address:          randomdata.Address(),
			UserType:         r.Intn(4),
			CreatedTime:      time.Now(),
			LastModifiedTime: time.Now(),
		}
		if i%7 == 0 {
			user.Status = "DISABLE"
		}
		if i == 10 || i == 20 || i == 30 || i == 40 || i == 50 || i == 60 || i == 70 || i == 80 {
			user.IsDelete = true
		}
		userArray = append(userArray, user)
	}
	client := NewClient(url)
	var bulkArray []BulkDTO
	for _, user := range userArray {
		bulkArray = append(bulkArray, BulkDTO{
			Index:    "user",
			Type:     "tags",
			Data:     user,
			IsDelete: false,
		})
	}
	resp, err := client.Bulk(bulkArray)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Created"))
	for _, res := range resp.Created() {
		resJSON, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resJSON))
	}
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Updated"))
	for _, res := range resp.Updated() {
		resJSON, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resJSON))
	}
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Deleted"))
	for _, res := range resp.Deleted() {
		resJSON, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resJSON))
	}
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Succeeded"))
	for _, res := range resp.Succeeded() {
		resJSON, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resJSON))
	}
	fmt.Println(fmt.Sprintf("---------- %s ----------", "Failed"))
	for _, res := range resp.Failed() {
		resJSON, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resJSON))
	}
}

func TestClient_CatHealth(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		fmt.Println(r.Int63n(2))
		time.Sleep(500)
	}
}
