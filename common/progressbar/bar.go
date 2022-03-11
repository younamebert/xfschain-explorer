package progressbar

import (
	"fmt"
	"xfschainbrowser/global"
)

type Bar struct {
	percent       int64  //百分比
	cur           int64  //当前进度位置
	total         int64  //总进度
	rate          string //进度条
	graph         string //显示符号
	lastHeight    int64  //最高的区块
	currentHeight int64  //当前高度
}

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph //初始化进度条位置
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

func (bar *Bar) Play(cur, lastHeight, currentHeight int64) {
	bar.cur = cur
	bar.currentHeight = currentHeight
	bar.lastHeight = lastHeight
	last := bar.percent
	bar.percent = bar.getPercent()

	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	global.GVA_LOG.Info(fmt.Sprintf("sync mode:desc order synchronous progress: %3d%%  %8d/%d currentHeight=%v lastHeight=%v", bar.percent, bar.cur, bar.total, currentHeight, lastHeight))
	// global.GVA_LOG.Info(fmt.Sprintf("\r[%-50s]%3d%%  %8d/%d currentHeight=%v lastHeight=%v sync mode:desc order synchronous", bar.rate, bar.percent, bar.cur, bar.total, currentHeight, lastHeight))

	// fmt.Printf("\r%3d%%  %8d/%d currentHeight=%v lastHeight=%v sync mode:desc order synchronous", bar.percent, bar.cur, bar.total, currentHeight, lastHeight)
	// fmt.Printf("lastHeight=%v currentHeight=%v")
}

func (bar *Bar) Finish() {
	// fmt.Println()
}
