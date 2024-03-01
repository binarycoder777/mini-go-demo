package search

import (
	"fmt"
	"log"
)

// Result 搜索结果
type Result struct {
	Field   string
	Content string
}

// Matcher 搜索类型的行为
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

// Match 匹配函数，由每个goroutine并发执行
func Match(match Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	searchResults, err := match.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}
	for _, result := range searchResults {
		results <- result
	}
}

// Display 从每个单独的 goroutine 接收到结果后在终端输出
func Display(results chan *Result) {
	for result := range results {
		fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
	}
}
