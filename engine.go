package spider

import (
	"errors"
	"regexp"
	"time"

	"gitlab.wcxst.com/jormin/gofetcher/pkg/fetcher"
	"gitlab.wcxst.com/jormin/gohelper"
)

// 引擎
type Engine struct {
	// 采集器配置
	FetcherCfg FetcherConfig
	// 解析结果
	parseResults chan *ParseResult
	// 采集器客户端
	fetcherClients chan *fetcher.Client
	// 采集任务
	fetchJobs chan *FetchJob
	// 解析任务
	parseJobs chan *ParseJob
	// 存储任务
	saveJobs chan *SaveJob
	// 解析器
	Parsers map[string]Parser
	// 存储器
	Savers map[string]Saver
}

// 采集器配置
type FetcherConfig struct {
	PerClientNum int      `json:"per_client_num"`
	Addr         []string `json:"addr"`
}

// 获取引擎
func NewEngine(FetcherCfg FetcherConfig, parsers map[string]Parser, savers map[string]Saver) *Engine {
	e := &Engine{
		FetcherCfg:     FetcherCfg,
		parseResults:   make(chan *ParseResult, 1000),
		fetcherClients: nil,
		fetchJobs:      make(chan *FetchJob, 1000),
		parseJobs:      make(chan *ParseJob, 1000),
		saveJobs:       make(chan *SaveJob, 1000),
		Parsers:        parsers,
		Savers:         savers,
	}
	return e
}

// 运行
func (e *Engine) Run() {
	headers := map[string]string{
		"Referer":                   "https://xa.fang.anjuke.com/loupan/?from=navigation",
		"sec-ch-ua":                 `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`,
		"sec-ch-ua-mobile":          "?0",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
	}

	// 初始化采集器
	gohelper.Must(e.FetcherInit())

	// 采集任务
	var fetchJob *FetchJob
	// 解析任务
	var parseJob *ParseJob
	// 存储任务
	var saveJob *SaveJob

	for {
		select {
		// 处理采集任务
		case fetchJob = <-e.fetchJobs:
			go func(fetchJob *FetchJob) {
				// 获取客户端
				fetcherClient := e.getFetcherClient()
				// 归还客户端
				defer e.backFetcherClient(fetcherClient)
				content, err := fetcherClient.Fetch(fetchJob.Title, fetchJob.Url, 20*time.Millisecond, headers)
				if err != nil {
					return
				}
				// 提交解析任务
				e.SubmitParseJob(fetchJob, &content)
			}(fetchJob)
		// 处理解析任务
		case parseJob = <-e.parseJobs:
			go func(parseJob *ParseJob) {
				// 检测有没有对应的解析器
				parser, ok := e.Parsers[fetchJob.Tag]
				if !ok {
					return
				}
				result, err := parser.Parse(parseJob)
				if err != nil {
					return
				}
				e.SubmitParseResult(result)
			}(parseJob)
		// 处理存储任务
		case saveJob = <-e.saveJobs:
			go func(saveJob *SaveJob) {
				// 检测有没有对应的存储器
				saver, ok := e.Savers[fetchJob.Tag]
				if !ok {
					return
				}
				err := saver.Save(saveJob)
				if err != nil {
					return
				}
			}(saveJob)
		// 处理解析结果
		case result := <-e.parseResults:
			go func(result *ParseResult) {
				e.DealParseResult(result)
			}(result)
		default:
			// os.Exit(1)
		}
	}
}

// 初始化采集器
func (e *Engine) FetcherInit() error {
	if len(e.FetcherCfg.Addr) == 0 {
		return errors.New("invalid fetcher addresses")
	}
	totalClientNum := e.FetcherCfg.PerClientNum * len(e.FetcherCfg.Addr)
	e.fetcherClients = make(chan *fetcher.Client, totalClientNum)
	for index, addr := range e.FetcherCfg.Addr {
		for len(e.fetcherClients) < e.FetcherCfg.PerClientNum*(index+1) {
			fetcherClient, err := fetcher.NewClient(addr)
			if err != nil {
				continue
			}
			e.fetcherClients <- fetcherClient
		}
	}
	return nil
}

// 获取采集器客户端
func (e *Engine) getFetcherClient() *fetcher.Client {
	return <-e.fetcherClients
}

// 归还采集器客户端
func (e *Engine) backFetcherClient(fetcherClient *fetcher.Client) {
	e.fetcherClients <- fetcherClient
}

// 提交采集任务
func (e *Engine) SubmitFetchJob(job *FetchJob) {
	e.fetchJobs <- job
}

// 提交解析任务
func (e *Engine) SubmitParseJob(job *FetchJob, content *[]byte) {
	e.parseJobs <- &ParseJob{
		FetchJob: job,
		Content:  content,
	}
}

// 提交存储任务
func (e *Engine) SubmitSaveJob(job *SaveJob) {
	e.saveJobs <- job
}

// 提交解析结果
func (e *Engine) SubmitParseResult(result *ParseResult) {
	if result != nil {
		e.parseResults <- result
	}
}

// 处理解析结果
func (e *Engine) DealParseResult(result *ParseResult) {
	// 处理采集任务
	if len(result.FetchJobs) > 0 {
		for _, fetchJob := range result.FetchJobs {
			e.SubmitFetchJob(fetchJob)
		}
	}
	// 处理存储任务
	if len(result.SaveJobs) > 0 {
		for _, saveJob := range result.SaveJobs {
			e.SubmitSaveJob(saveJob)
		}
	}
}

// 从Url中获取Code
func GetCodeFromUrl(url string, index int) string {
	return regexp.MustCompile("[0-9]+").FindAllString(url, -1)[index]
}
