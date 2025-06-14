package spiders

import (
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/XXueTu/wise/pkg/spiders/wechat"
)

type Pattern struct {
	patternMap map[string]PatternInterface
}

type PatternInterface interface {
	Identification(url string) bool
	GetData(url string) (string, string, error)
}

func NewPattern() *Pattern {
	// 初始化
	p := &Pattern{}
	p.patternMap = make(map[string]PatternInterface)
	p.patternMap["wechat"] = wechat.Init()
	return p
}

func (p *Pattern) GetPatternTypes(url string) string {
	for k := range p.patternMap {
		if p.patternMap[k].Identification(url) {
			return k
		}
	}
	return "unknown"
}

// title content error
func (p *Pattern) GetPattern(url string) (string, string, error) {
	for k, v := range p.patternMap {
		logx.Infof("spider name:%s,url:%s", k, url)
		if v.Identification(url) {
			return v.GetData(url)
		}
	}
	// title content error
	return "unknown", "unknown", errors.New("not supported spider")
}
