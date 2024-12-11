package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type IPResponse struct {
	IPv4 string `json:"ipv4,omitempty"`
	IPv6 string `json:"ipv6,omitempty"`
}

func getIP(r *http.Request) (string, string) {
	var ipv4, ipv6 string

	// 获取请求头中的 X-Real-IP 或 X-Forwarded-For
	ipList := r.Header.Get("X-Forwarded-For")
	if ipList != "" {
		// 假设第一个是客户端的真实 IP
		ipv4 = ipList
	} else {
		ipv4 = r.RemoteAddr
	}

	// 确定是否在 IPv6 格式
	if len(ipv4) > 0 && ipv4[0] == '[' {
		ipv6 = ipv4[1 : len(ipv4)-1] // 移除方括号
		ipv4 = ""
	}

	return ipv4, ipv6
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ipv4, ipv6 := getIP(r)
	response := IPResponse{}

	if r.URL.Query().Get("format") == "v4" && ipv4 != "" {
		response.IPv4 = ipv4
	} else if r.URL.Query().Get("format") == "v6" && ipv6 != "" {
		response.IPv6 = ipv6
	} else {
		// 默认输出 IPv4 和 IPv6
		response.IPv4 = ipv4
		response.IPv6 = ipv6
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/get-ip", ipHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
