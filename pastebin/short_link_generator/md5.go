package shortlinkgeneratorgo

import (
	"solid/paste"
)

func init() {
	paste.Register("md5", func(baseUrl string) paste.ShortLinkGenerator {
		return &md5Generator{baseUrl}
	})
}

type md5Generator struct {
	baseUrl string
}

func (g *md5Generator) GenerateShortLink() string {
	// Lisvoc substitution principle
	// ここでは簡略化のため、常に同じ短縮URLを返すようにしています。
	// 実際の実装では、interfaceで定義した事前条件と事後条件を元に実装する
	return "md5_short_link"
}
