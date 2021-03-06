package teaproxy

import (
	"bytes"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

type testResponseWriter struct {
	a    *assert.Assertion
	data []byte
}

func testNewResponseWriter(a *assert.Assertion) *testResponseWriter {
	return &testResponseWriter{
		a: a,
	}
}

func (this *testResponseWriter) Header() http.Header {
	return http.Header{}
}

func (this *testResponseWriter) Write(data []byte) (int, error) {
	this.data = append(this.data, data ...)
	return len(data), nil
}

func (this *testResponseWriter) WriteHeader(statusCode int) {
}

func (this *testResponseWriter) Close() {
	this.a.Log(string(this.data))
}

func TestRequest_Call(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	request := NewRequest(nil)
	err := request.Call(writer)
	a.IsNotNil(err)
	if err != nil {
		a.Log(err.Error())
	}
}

func TestRequest_CallRoot(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	request := NewRequest(nil)
	request.root = Tea.ViewsDir() + "/@default"
	request.uri = "/layout.css"
	err := request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallBackend(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.backend = &teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1",
	}
	request.backend.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallProxy(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx", nil)
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"

	proxy := teaconfigs.NewServerConfig()
	proxy.AddBackend(&teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1:80",
	})
	/**proxy.AddBackend(&teaconfigs.ServerBackendConfig{
		Address: "127.0.0.1:81",
	})**/
	request.proxy = proxy

	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallFastcgi(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx/box/version", bytes.NewBuffer([]byte("hello=world")))
	//req, err := http.NewRequest("GET", "/index.php", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.serverAddr = "127.0.0.1:80"

	request.fastcgi = &teaconfigs.FastcgiConfig{
		Params: map[string]string{
			"SCRIPT_FILENAME": "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk/index.php",
			//"DOCUMENT_ROOT":   "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk",
		},
		Pass: "127.0.0.1:9000",
	}
	request.fastcgi.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_CallFastcgiPerformance(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()
	writer := testNewResponseWriter(a)

	req, err := http.NewRequest("GET", "/index.php?__ACTION__=/@wx/box/version", bytes.NewBuffer([]byte("hello=world")))
	//req, err := http.NewRequest("GET", "/index.php", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		a.Fatal(err)
	}
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	request := NewRequest(req)
	request.scheme = "http"
	request.host = "wx.balefm.cn"
	request.serverAddr = "127.0.0.1:80"

	request.fastcgi = &teaconfigs.FastcgiConfig{
		Params: map[string]string{
			"SCRIPT_FILENAME": "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk/index.php",
			//"DOCUMENT_ROOT":   "/Users/liuxiangchao/Documents/Projects/pp/apps/baleshop.ppk",
		},
		Pass: "127.0.0.1:9000",
	}
	request.fastcgi.Validate()
	err = request.Call(writer)
	a.IsNil(err)
	writer.Close()

	a.Log("status:", request.responseStatus, request.responseStatusMessage)
	a.Log("requestTime:", request.requestTime)
	a.Log("bytes send:", request.responseBytesSent, request.responseBodyBytesSent)
}

func TestRequest_Format(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}
	rawReq.RemoteAddr = "127.0.0.1:1234"
	rawReq.Header.Add("Content-Type", "text/plain")

	req := NewRequest(rawReq)
	req.uri = "/hello/world?name=Lu&age=20"
	req.method = "GET"
	req.filePath = "hello.go"
	req.scheme = "http"

	a.IsTrue(req.requestRemoteAddr() == "127.0.0.1:1234")
	a.IsTrue(req.requestRemotePort() == 1234)
	a.IsTrue(req.requestURI() == req.uri)
	a.IsTrue(req.requestPath() == "/hello/world")
	a.IsTrue(req.requestMethod() == "GET")
	a.IsTrue(req.requestLength() > 0)
	a.IsTrue(req.requestFilename() == req.filePath)
	a.IsTrue(req.requestProto() == "HTTP/1.1")
	a.IsTrue(req.requestQueryString() == "name=Lu&age=20")
	a.IsTrue(req.requestQueryParam("name") == "Lu")

	t.Log(req.format("hello ${teaVersion} remoteAddr:${remoteAddr} name:${arg.name} header:${header.Content-Type} test:${test}"))
}

func TestRequest_Index(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	req := NewRequest(rawReq)
	req.index = []string{}
	t.Log(req.findIndexFile(Tea.Root))

	req.index = []string{"main.go", "main2.go", "run.sh"}
	a.Equals(req.findIndexFile(Tea.Root), "main.go")

	req.index = []string{"main.*"}
	a.Equals(req.findIndexFile(Tea.Root), "main.go")
}

func TestRequest_LocationVariables(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	server := teaconfigs.NewServerConfig()
	server.Root = "/home"

	{
		location := teaconfigs.NewLocationConfig()
		location.On = true
		location.Pattern = "~ /hello/(\\w)(\\w+)"
		location.Root = "/hello/${1}/${host}"
		location.Index = []string{"hello_${1}${2}"}
		location.Charset = "${arg.charset}"
		location.SetHeader("hello", "${1}")
		err := location.Validate()
		a.IsNil(err)

		server.AddLocation(location)

		matches, ok := location.Match("/hello/world")
		if ok {
			t.Log(matches)
		}
	}

	err = server.Validate()
	a.IsNil(err)

	req := NewRequest(rawReq)
	req.uri = "/hello/world?charset=utf-8"
	req.host = "www.example.com"

	err = req.configure(server, 0)
	if err != nil {
		t.Log(err.Error())
	}
	a.IsNil(err)

	t.Log("request uri:", req.requestURI())
	t.Log("root:", req.root)
	t.Log("index:", req.index)
	t.Log("charset:", req.charset)

	for _, header := range req.headers {
		t.Log("headers:", header.Name, ":", header.Value)
	}
}

func TestRequest_RewriteVariables(t *testing.T) {
	a := assert.NewAssertion(t).Quiet()

	rawReq, err := http.NewRequest("GET", "http://www.example.com/hello/world?name=Lu&age=20", bytes.NewBuffer([]byte("hello=world")))
	if err != nil {
		t.Fatal(err)
	}

	server := teaconfigs.NewServerConfig()
	server.Root = "/home/${arg.charset}"
	server.Charset = "[${arg.charset}]"
	server.AddHeader(&teaconfigs.HeaderConfig{
		Name:  "Charset",
		Value: "${arg.charset}",
	})

	{
		location := teaconfigs.NewLocationConfig()
		location.On = true
		location.Pattern = "/"

		rewriteRule := teaconfigs.NewRewriteRule()
		rewriteRule.Pattern = "^/hello/(\\w+)$"
		rewriteRule.Replace = "/he/${1}${requestPath}?arg=${arg.charset}"
		location.AddRewriteRule(rewriteRule)

		err := location.Validate()
		a.IsNil(err)

		server.AddLocation(location)
	}

	err = server.Validate()
	a.IsNil(err)

	req := NewRequest(rawReq)
	req.uri = "/hello/world?charset=utf-8"
	req.host = "www.example.com"

	err = req.configure(server, 0)
	if err != nil {
		t.Log(err.Error())
	}
	a.IsNil(err)

	t.Log("request uri:", req.uri)
	t.Log("root:", req.root)
	t.Log("index:", req.index)
	t.Log("charset:", req.charset)

	for _, header := range req.headers {
		t.Log("headers:", header.Name, ":", header.Value)
	}
}

func TestPerformanceBackend(t *testing.T) {
	beforeTime := time.Now()

	countSuccess := 0
	countFail := 0

	locker := sync.Mutex{}
	wg := sync.WaitGroup{}
	threads := 1000
	connections := 100
	wg.Add(threads)

	for i := 0; i < threads; i ++ {
		go func() {
			for j := 0; j < connections; j ++ {
				req, err := http.NewRequest("GET", "http://127.0.0.1:9992/benchmark", nil)

				if err != nil {
					t.Fatal(err)
				}

				c := SharedClientPool.client("127.0.0.1:9992")
				resp, err := c.Do(req)

				if err != nil {
					locker.Lock()
					countFail ++
					locker.Unlock()
				} else {
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil || len(data) == 0 || strings.Index(string(data), "benchmark") == -1 {
						locker.Lock()
						countFail ++
						locker.Unlock()
					} else {
						locker.Lock()
						countSuccess ++
						locker.Unlock()
					}

					//io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	t.Log("success:", countSuccess, "fail:", countFail, "qps:", int(float64(countSuccess+countFail)/time.Since(beforeTime).Seconds()))
}

func TestPerformanceStatic(t *testing.T) {
	beforeTime := time.Now()

	countSuccess := 0
	countFail := 0

	locker := sync.Mutex{}
	wg := sync.WaitGroup{}
	threads := 100
	connections := 100
	wg.Add(threads)

	for i := 0; i < threads; i ++ {
		go func() {
			for j := 0; j < connections; j ++ {
				req, err := http.NewRequest("GET", "http://127.0.0.1:9993/css/semantic.min.css", nil)

				if err != nil {
					t.Fatal(err)
				}

				c := SharedClientPool.client("127.0.0.1:9993")
				resp, err := c.Do(req)

				if err != nil {
					locker.Lock()
					countFail ++
					locker.Unlock()
				} else {
					data, err := ioutil.ReadAll(resp.Body)
					if err != nil || len(data) == 0 || strings.Index(string(data), "Semantic") == -1 {
						locker.Lock()
						countFail ++
						locker.Unlock()
					} else {
						locker.Lock()
						countSuccess ++
						locker.Unlock()
					}

					//io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	t.Log("success:", countSuccess, "fail:", countFail, "qps:", int(float64(countSuccess+countFail)/time.Since(beforeTime).Seconds()))
}