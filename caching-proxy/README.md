https://roadmap.sh/projects/caching-server

# Caching Proxy CLI

A command-line tool written in Go that starts a proxy server with caching capabilities. It forwards requests to an origin server and stores responses locally to serve cached data for repeated requests.

---

## 📦 Requirements

- Go 1.20+
- An accessible origin server (API to be proxied and cached)

---

## 🚀 How to Run the Project

Run the following command in your terminal:

```bash
go run main.go caching-proxy --port {port} --origin {originURL}
