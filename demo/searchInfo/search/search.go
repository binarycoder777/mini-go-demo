package search

import (
	"log"
	"sync"
)

// 注册用于搜索的匹配器的映射
var matchers = make(map[string]Matcher)

// Run 执行搜索
func Run(searchTerm string) {
	// 获取需要搜索的数据源列表
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个无缓冲的通道，接受匹配后的结果
	results := make(chan *Result)

	// 构造一个waitGroup，处理所有的数据源
	var waitGroup sync.WaitGroup

	// 设置需要等待处理
	// 每个数据源的goroutine数量
	waitGroup.Add(len(feeds))

	// 为每个数据源启动goroutine并行查找
	for _, feed := range feeds {
		// 获取数据源的匹配器用于查找
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		// 启动一个goroutine查询
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			defer waitGroup.Done()
		}(matcher, feed)
	}

	// 启动一个goroutine来监控是否所以得工作都完成了
	go func() {
		// 等候所有任务完成
		waitGroup.Wait()
		// 关闭通道，通知Display函数
		close(results)
	}()

	// 显示返回结果
	Display(results)
}

// Register 调用时，会注册一个匹配器，提供给后面的程序使用
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}
	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
