package shortlinkgeneratorgo

import (
	"solid/paste"
)

func init() {
	paste.Register("sha256", func(baseUrl string) paste.ShortLinkGenerator {
		return &sha256Generator{baseUrl}
	})
}

type sha256Generator struct {
	baseUrl string
}

func (g *sha256Generator) GenerateShortLink() string {
	return "sha256_short_link"
}
