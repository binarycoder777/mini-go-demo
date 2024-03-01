package search

import (
	"encoding/json"
	"os"
)

const dataFile = "data/data.json"

// Feed 处理的数据源信息
type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

// RetrieveFeeds 读取并反序列化数据源文件
func RetrieveFeeds() ([]*Feed, error) {
	// open file
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	// close file
	defer file.Close()

	// 将文件解码到一个切片
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	return feeds, err
}
