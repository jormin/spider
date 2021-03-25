package spider

// 任务
type FetchJob struct {
	Tag   string `json:"tag" remark:"标签"`
	Title string `json:"title" remark:"标题"`
	Url   string `json:"url" remark:"采集Url"`
}

// 解析器
type Parser interface {
	// 解析
	Parse(*ParseJob) (*ParseResult, error)
}

// 解析结果
type ParseResult struct {
	FetchJobs []*FetchJob
	SaveJobs  []*SaveJob
}

// 解析任务
type ParseJob struct {
	FetchJob *FetchJob `json:"fetch_job"`
	Content  *[]byte   `json:"content" remark:"采集内容"`
}

// 存储器
type Saver interface {
	// 存储
	Save(*SaveJob) error
}

// 存储任务
type SaveJob struct {
	FetchJob *FetchJob   `json:"fetch_job"`
	Data     interface{} `json:"data" remark:"解析数据"`
}
