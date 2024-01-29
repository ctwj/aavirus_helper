package upload

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/ctwj/aavirus_helper/internal/lib"
	"github.com/ctwj/aavirus_helper/internal/pkg/config"
)

type Web struct {
	server *http.Server
	wg     sync.WaitGroup
	port   int
	init   bool
}

func NewWeb() *Web {
	return &Web{}
}

// 处理请求
func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello!")
}

// 设置静态目录为网站
func (t *Web) StartServer(port int) interface{} {

	t.port = port
	addr := fmt.Sprintf(":%d", port)
	if !t.init { // 不能重复注册路由
		fs := http.FileServer(http.Dir(config.OutputDir))
		http.Handle("/", fs)
		http.Handle("/ping", http.HandlerFunc(handleRequest))
		log.Printf("Starting server on http://localhost%s", addr)
	}
	t.init = true

	t.server = &http.Server{Addr: addr}

	resultChan := make(chan interface{}, 1)
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			resultChan <- map[string]interface{}{"status": false, "err": err.Error()}
			// return map[string]interface{}{"status": false, "err": err.Error()}
		}
		resultChan <- map[string]interface{}{"status": true}
		// return map[string]interface{}{"status": true}
	}()

	// return map[string]interface{}{"status": true}
	return <-resultChan
}

// 停止服务器
func (t *Web) StopServer() interface{} {
	err := t.server.Shutdown(nil)
	if err != nil {
		log.Fatalf("Server shutdown error: %v", err)
		return map[string]bool{"status": false}
	}
	log.Println("Server stopped")
	return map[string]bool{"status": true}
}

// 获取网站状态
func (t *Web) GetStatus(port int) interface{} {

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/ping", port), nil)
	if err != nil {
		fmt.Printf("Check server running error: %v \r\n", err)
		return map[string]bool{"status": false}
	}

	resp, err := client.Do(req)
	if err != nil {
		// 检查错误类型
		if err, ok := err.(*url.Error); ok && err.Timeout() {
			fmt.Printf("Check server running error: request timeout \r\n")
			return map[string]bool{"status": false}
		}

		fmt.Printf("Check server running error: %v  \r\n", err)
		return map[string]bool{"status": false}
	}
	defer resp.Body.Close()

	return map[string]bool{"status": true}
}

func (t *Web) GetIps() interface{} {
	ips, _ := lib.GetLocalIPs()
	return ips
}
