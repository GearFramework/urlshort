package pkg

type ApiShortener interface {
	EncodeURL(url string) string
	DecodeURL(shortURL string) (string, error)
	AddShortly(url, code string)
}
