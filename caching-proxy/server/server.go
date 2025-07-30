package server

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/freitasmatheusrn/caching-proxy/utils"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Address string
	Origin  string
	Cache   sync.Map
}

func NewServer(addr, origin string) *Server {
	return &Server{
		Address: addr,
		Origin:  origin,
	}
}

func (s *Server) cacheProxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		normalized := utils.NormalizeRequest(r)
		hash := utils.HashString(normalized)
		if cached, ok := s.Cache.Load(hash); ok {
			fmt.Println("✅ HIT cache:", normalized)
			w.Header().Set("X-Cache", "HIT")
			w.Write(cached.([]byte))
			return
		}
		targetURL := s.Origin + r.URL.Path
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}
		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, "Error reaching origin server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response", http.StatusInternalServerError)
			return
		}
		s.Cache.Store(hash, body)
		fmt.Println("🆕 MISS cache. Stored:", normalized)

		w.Header().Set("X-Cache", "MISS")
		w.Write(body)
	})
}

func (s *Server) ClearCache() {
	s.Cache.Clear()
}

func (s *Server) InitServer() error {
	r := chi.NewRouter()
	r.Use(s.cacheProxyMiddleware)
	r.Use(middleware.Logger)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not supported route", http.StatusNotFound)
	})
	addr := fmt.Sprintf(":%s", s.Address)
	return http.ListenAndServe(addr, r)
}
