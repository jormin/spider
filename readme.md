## About

爬虫

## Use

### 调用方法
```go
const (
	TagList = "list"
)

// 初始化
engine := spider.NewEngine(
    spider.FetcherConfig{
        PerClientNum: 5,
        Addr: []string{
            "127.0.0.1:10001",
            "127.0.0.1:10002",
            "127.0.0.1:10003",
            "127.0.0.1:10004",
        },
    },
    map[string]spider.Parser{
        TagList: &ListParser{},
    },
    map[string]spider.Saver{
        TagList: &ListSaver{},
    },
)

// 提交采集任务
engine.SubmitFetchJob(
    &spider.FetchJob{
        Tag:   TagList,
        Title: "your title",
        Url:   "your url",
    },
)

// 执行
engine.Run()
```

### 采集器

```go
// 解析器
type ListParser struct{}

// 自定义的解析方法
func (p *ListParser) Parse(job *spider.ParseJob) (*spider.ParseResult, error) {
    return nil
}
```

### 存储器

```go
// 存储器
type ListSaver struct {}

// 自定义的存储方法
func (l ListSaver) Save(job *spider.SaveJob) error {
    return nil
}
```