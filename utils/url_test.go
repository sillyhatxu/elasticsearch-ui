package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatURL(t *testing.T) {
	url := FormatURL("http", "127.0.0.1:8080")
	fmt.Println(url)
	assert.EqualValues(t, url, "http://127.0.0.1:8080")

	url = FormatURL("https", "127.0.0.1:80")
	fmt.Println(url)
	assert.EqualValues(t, url, "https://127.0.0.1:80")
}
