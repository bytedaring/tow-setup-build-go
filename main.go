package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	targetURL := "https://ark.cn-beijing.volces.com"
	fmt.Printf("Attempting to connect to %s ...\n", targetURL)

	// 1. 克隆默认的 http.Transport 以继承其所有优化设置
	//    (如连接池、Keep-Alive 等)，这是比从零创建更好的做法。
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// 2. 自定义 DialContext 函数，这是核心所在。
	//    这个函数会在每次需要建立新的 TCP 连接时被调用。
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 使用标准的 net.Dialer 来建立实际的连接
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		conn, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		// 关键步骤：连接成功后，从连接对象(conn)中获取远程地址并打印
		fmt.Printf(">>> Successfully connected to IP: %s\n", conn.RemoteAddr().String())
		return conn, nil
	}

	// 3. 创建一个使用我们自定义 Transport 的 http.Client
	client := &http.Client{
		Transport: transport,
	}

	// 4. 使用这个自定义的 client 发起请求，而不是全局的 http.Get()
	resp, err := client.Get(targetURL)
	if err != nil {
		fmt.Printf("Error: Request failed: %v\n", err)
		os.Exit(1) // 以失败状态码退出
	}
	defer resp.Body.Close()

	fmt.Printf("Success! Status Code: %d\n", resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response Body: %s\n", string(body))
}
