package paste

/*
 * このservice.goはDIPを満たすためのコードである。
 * ShortLinkGenerator interfaceが下位実装と同じパッケージにあるべきではなく、DIPの原則では抽象は上位モジュールが所有すべき（もしくは独立した層に置く）。
 * Goの慣習では「interfaceは使う側が定義する」なので
 */

// 全体実装制約
// ShortLinkGenerator generates short URLs based on a base URL.
// All implementations must produce a non-empty URL that begins
// with the baseUrl provided at construction time.
type ShortLinkGenerator interface {
	// 事後条件
	// GenerateShortLink returns a short URL.
	// The returned string is always a valid URL starting with the baseUrl
	// passed to the constructor (e.g., "https://example.com/abc123").
	GenerateShortLink() string
}

// 事前条件:
// GeneratorConstructor creates a ShortLinkGenerator for the given baseUrl.
// baseUrl must be a valid URL ending with "/" (e.g., "https://example.com/").
type GeneratorConstructor func(baseUrl string) ShortLinkGenerator

var registry = map[string]GeneratorConstructor{}

func Register(name string, constructor GeneratorConstructor) {
	registry[name] = constructor
}

// lspを満たすため、ここの戻り値は error を返すべき、理由として、呼び出し元はnilを期待していないため、lsp違反になる(事後条件)
// ※ここはわざと直してない
func ShortLinkGeneratorFactory(algorithmType string, baseUrl string) ShortLinkGenerator {
	constructor, ok := registry[algorithmType]
	if !ok {
		return nil
	}
	return constructor(baseUrl)
}
