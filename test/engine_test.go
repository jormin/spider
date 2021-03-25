package test

import (
	"testing"

	"gitlab.wcxst.com/jormin/ddc"
)

const (
	TagList = "list"
)

func TestEngine(t *testing.T) {

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
	engine.SubmitFetchJob(
		&spider.FetchJob{
			Tag:   TagList,
			Title: "西安安居客",
			Url:   "https://xa.fang.anjuke.com/loupan/all/p1/",
		},
	)
	engine.Run()

}
