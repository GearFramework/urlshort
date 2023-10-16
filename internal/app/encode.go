package app

import (
	"fmt"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"math/rand"
	"runtime"
	"time"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lenAlpha    = len(alphabet)
	defShortLen = 8
)

func (app *ShortApp) EncodeURL(url string) string {
	app.Store.Lock()
	defer app.Store.Unlock()
	code, exists := app.Store.GetCode(url)
	if !exists {
		code = app.getRandomString(defShortLen)
		if err := app.Store.Insert(url, code); err != nil {
			logger.Log.Error(err.Error())
		}
	}
	return fmt.Sprintf("%s/%s", app.Conf.ShortURLHost, code)
}

func (app *ShortApp) getRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = alphabet[rnd.Intn(lenAlpha)]
	}
	return string(b)
}

func (app *ShortApp) BatchEncodeURL(batch []pkg.BatchURLs) []pkg.ResultBatchShort {
	app.Store.Lock()
	defer app.Store.Unlock()
	res := []pkg.ResultBatchShort{}
	trc, urls := transformBatchByCorrelation(batch)
	for _, chunkURLs := range chunkingURLs(urls) {
		existCodes := app.Store.GetCodeBatch(chunkURLs)
		for existUrl, existCode := range existCodes {
			res = append(res, pkg.ResultBatchShort{
				CorrelationId: trc[existUrl],
				ShortURL:      fmt.Sprintf("%s/%s", app.Conf.ShortURLHost, existCode),
			})
		}
		if len(chunkURLs) == len(existCodes) {
			continue
		}
		notExistCodes := getNotExists(chunkURLs, existCodes)
		newShortUrls, pack := app.prepareNotExistsShortURLs(trc, notExistCodes)
		if err := app.Store.InsertBatch(pack); err != nil {
			logger.Log.Error(err.Error())
			continue
		}
		res = append(res, newShortUrls...)
	}
	return res
}

func (app *ShortApp) prepareNotExistsShortURLs(
	trc map[string]string,
	notExistCodes []string,
) ([]pkg.ResultBatchShort, [][]string) {
	res := []pkg.ResultBatchShort{}
	pack := [][]string{}
	for _, url := range notExistCodes {
		code := app.getRandomString(defShortLen)
		res = append(res, pkg.ResultBatchShort{
			CorrelationId: trc[url],
			ShortURL:      fmt.Sprintf("%s/%s", app.Conf.ShortURLHost, code),
		})
		pack = append(pack, []string{url, code})
	}
	return res, pack
}

func transformBatchByCorrelation(batch []pkg.BatchURLs) (map[string]string, []string) {
	trc := map[string]string{}
	urls := []string{}
	for _, packet := range batch {
		trc[packet.OriginalURL] = packet.CorrelationId
		urls = append(urls, packet.OriginalURL)
	}
	return trc, urls
}

func chunkingURLs(urls []string) [][]string {
	var chunks [][]string
	count := len(urls)
	numCPU := runtime.NumCPU()
	var chunkSize int
	if count <= numCPU {
		chunkSize = count
	} else {
		chunkSize = (count + numCPU - 1) / numCPU
	}
	for i := 0; i < count; i += chunkSize {
		end := i + chunkSize
		if end > count {
			end = count
		}
		chunks = append(chunks, urls[i:end])
	}
	return chunks
}

func getNotExists(chunkURLs []string, exists map[string]string) []string {
	if len(chunkURLs) == len(exists) {
		return []string{}
	}
	notExists := []string{}
	for _, url := range chunkURLs {
		if _, ok := exists[url]; ok {
			continue
		}
		notExists = append(notExists, url)
	}
	return notExists
}
