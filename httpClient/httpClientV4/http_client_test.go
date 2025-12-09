package httpClientV4

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClientBasic(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","data":"test"}`))
	}))
	defer server.Close()

	// 测试 GET 请求
	client := GET(
		URL(server.URL),
		AppendHeader("User-Agent", "TestClient/1.0"),
		Timeout(5_000_000_000), // 5 seconds
	)
	defer client.Release() // 归还到池中

	client.Send()

	if err := client.Error(); err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if client.GetStatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", client.GetStatusCode())
	}

	body := client.ToBytes()
	if len(body) == 0 {
		t.Error("Response body is empty")
	}
}

func TestHTTPClientPOST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()

	type RequestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	client := POST(
		URL(server.URL),
		JSON(RequestData{Name: "John", Age: 30}),
	)
	defer client.Release()

	client.Send()

	if err := client.Error(); err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if client.GetStatusCode() != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", client.GetStatusCode())
	}
}

func TestHTTPClientQueries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("key1") != "value1" {
			t.Errorf("Expected key1=value1, got %s", query.Get("key1"))
		}
		if query.Get("key2") != "value2" {
			t.Errorf("Expected key2=value2, got %s", query.Get("key2"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	queries := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	client := GET(
		URL(server.URL),
		Queries(queries),
	)
	defer client.Release()

	client.Send()

	if err := client.Error(); err != nil {
		t.Fatalf("Request failed: %v", err)
	}
}

func TestHTTPClientPool(t *testing.T) {
	// 测试对象池功能
	client1 := Acquire()
	client2 := Acquire()

	if client1 == client2 {
		t.Error("Should get different instances from pool")
	}

	client1.Release()
	client3 := Acquire()

	// client3 可能是重用了 client1
	client2.Release()
	client3.Release()
}

// 基准测试 - 对比 V2 和 V4 的性能
func BenchmarkHTTPClientV4_GET(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client := GET(URL(server.URL))
			client.Send()
			_ = client.ToBytes()
			client.Release()
		}
	})
}

func BenchmarkHTTPClientV4_POST(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	type Data struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client := POST(
				URL(server.URL),
				JSON(Data{Name: "test", Value: 123}),
			)
			client.Send()
			client.Release()
		}
	})
}

func BenchmarkHTTPClientV4_WithPool(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := Acquire()
		client.init(http.MethodGet, URL(server.URL))
		client.Send()
		client.Release()
	}
}

func BenchmarkHTTPClientV4_WithoutPool(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := &HTTPClient{
			queries: make(map[string]string),
			headers: make(http.Header),
		}
		client.init(http.MethodGet, URL(server.URL))
		client.Send()
	}
}
