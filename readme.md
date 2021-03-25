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
        Title: "西安安居客",
        Url:   "https://xa.fang.anjuke.com/loupan/all/p1/",
    },
)

// 执行
engine.Run()
```

### 采集器

```go
// 解析器
type ListParser struct{}

func (p *ListParser) Parse(job *spider.ParseJob) (*spider.ParseResult, error) {
    return nil
}
```

### 存储器

```go
// 存储器
type ListSaver struct {}

// 存储
func (l ListSaver) Save(job *spider.SaveJob) error {
    return nil
}
```