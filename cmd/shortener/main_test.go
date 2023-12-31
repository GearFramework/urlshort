package main

import (
	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/server"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct {
	name string
	t    *testing.T
	enc  *TestEncode
	dec  *TestDecode
}

type TestEncode struct {
	requestEncode    Req
	responseExpected RespExpected
	responseActual   RespActualEncode
	testEnc          func(t *testing.T, test *TestEncode)
}

type TestDecode struct {
	requestDecode    Req
	responseExpected RespExpected
	responseActual   RespActualDecode
	testDec          func(t *testing.T, test *TestDecode)
}

type Req struct {
	Method string
	URL    string
}

type RespExpected struct {
	ResponseURL string
	StatusCode  int
}

type RespActualEncode struct {
	r           *http.Response
	ResponseURL string
}

type RespActualDecode struct {
	r          *http.Response
	StatusCode int
}

func (test *Test) test(t *testing.T, api pkg.APIShortener) {
	test.t = t
	if test.enc != nil {
		test.testEncode(api)
		return
	}
	if test.dec != nil {
		test.testDecode(api)
	}
}

func (test *Test) testEncode(api pkg.APIShortener) {
	request := httptest.NewRequest(test.enc.requestEncode.Method, "/", strings.NewReader(test.enc.requestEncode.URL))
	w := httptest.NewRecorder()
	s := server.NewServer(&server.Config{Addr: "localhost:8080"}, api)
	s.InitRoutes()
	s.Router.ServeHTTP(w, request)
	response := w.Result()
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	assert.NoError(test.t, err)
	assert.Equal(test.t, test.enc.responseExpected.StatusCode, response.StatusCode)
	test.enc.responseActual = RespActualEncode{response, string(body)}
	if test.enc.testEnc != nil {
		test.enc.testEnc(test.t, test.enc)
	}
	if response.StatusCode == http.StatusCreated && test.dec != nil {
		test.dec.requestDecode.URL = string(body)
		test.testDecode(api)
	}
}

func (test *Test) testDecode(api pkg.APIShortener) {
	request := httptest.NewRequest(test.dec.requestDecode.Method, test.dec.requestDecode.URL, nil)
	w := httptest.NewRecorder()
	s := server.NewServer(&server.Config{Addr: "localhost:8080"}, api)
	s.InitRoutes()
	s.Router.ServeHTTP(w, request)
	response := w.Result()
	defer response.Body.Close()
	assert.Equal(test.t, test.dec.responseExpected.StatusCode, response.StatusCode)
	if test.dec.testDec != nil {
		test.dec.responseActual = RespActualDecode{response, response.StatusCode}
		test.dec.testDec(test.t, test.dec)
	}
}

func getTests() []Test {
	return []Test{
		{
			name: "valid url encode valid method decode",
			enc: &TestEncode{
				requestEncode:    Req{http.MethodPost, "https://ya.ru"},
				responseExpected: RespExpected{StatusCode: http.StatusCreated},
				testEnc: func(t *testing.T, test *TestEncode) {
					assert.Regexp(t, "^http://localhost:8080/[a-zA-Z0-9]{8}$", test.responseActual.ResponseURL)
					assert.Equal(t, "text/plain", test.responseActual.r.Header.Get("Content-Type"))
				},
			},
			dec: &TestDecode{
				requestDecode:    Req{Method: http.MethodGet},
				responseExpected: RespExpected{"https://ya.ru", http.StatusTemporaryRedirect},
				testDec: func(t *testing.T, test *TestDecode) {
					assert.Equal(t, test.responseExpected.ResponseURL, test.responseActual.r.Header.Get("Location"))
				},
			},
		}, {
			name: "valid url encode invalid method decode",
			enc: &TestEncode{
				requestEncode:    Req{http.MethodPost, "https://yandex.ru"},
				responseExpected: RespExpected{StatusCode: http.StatusCreated},
				testEnc: func(t *testing.T, test *TestEncode) {
					assert.Regexp(t, "^http://localhost:8080/[a-zA-Z0-9]{8}$", test.responseActual.ResponseURL)
					assert.Equal(t, "text/plain", test.responseActual.r.Header.Get("Content-Type"))
				},
			},
			dec: &TestDecode{
				requestDecode:    Req{Method: http.MethodPut},
				responseExpected: RespExpected{StatusCode: http.StatusBadRequest},
			},
		}, {
			name: "invalid url",
			enc: &TestEncode{
				requestEncode:    Req{http.MethodPost, "https//ya.ru"},
				responseExpected: RespExpected{StatusCode: http.StatusBadRequest},
			},
		}, {
			name: "invalid request method",
			enc: &TestEncode{
				requestEncode:    Req{http.MethodDelete, "https://ya.ru"},
				responseExpected: RespExpected{StatusCode: http.StatusBadRequest},
			},
		}, {
			name: "invalid short url",
			dec: &TestDecode{
				requestDecode:    Req{Method: http.MethodGet, URL: "http://localhost:8080/8tbujofj"},
				responseExpected: RespExpected{StatusCode: http.StatusBadRequest},
			},
		},
	}
}

func TestHandleServiceEncode(t *testing.T) {
	a := app.NewShortener(config.GetConfig())
	for _, test := range getTests() {
		t.Run(test.name, func(t *testing.T) {
			test.test(t, a)
		})
	}
}
