package spider

// 采集任务
type FetchJob struct {
	Tag   string `json:"tag"`
	Title string `json:"title"`
	Url   string `json:"url"`
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
	Content  *[]byte   `json:"content"`
	Parser   Parser    `json:"parser"`
}

// 存储器
type Saver interface {
	// 存储
	Save(*SaveJob) error
}

// 存储任务
type SaveJob struct {
	FetchJob *FetchJob   `json:"fetch_job"`
	Data     interface{} `json:"data"`
	Saver    Saver       `json:"saver"`
}
