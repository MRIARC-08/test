package main

import (
	"fmt"
	"net/http"

	"github.com/MRIARC-08/test/pkg/cache"
)

func main() {
	appCache := cache.NewShardedCache()
	appCache.Set("service_status", "operational")

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			val, found := appCache.Get(key)
			if !found {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%v", val)
			return
		}

		if r.Method == http.MethodPost {
			val := r.URL.Query().Get("val")
			appCache.Set(key, val)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "saved %s", key)
			return
		}
	})

	fmt.Println("Server running on :8080...")
	_ = http.ListenAndServe(":8080", nil)
}
