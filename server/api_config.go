package server

import 	(
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform string
}

func (cfg *apiConfig) middleWareMetricInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(1)
	return next
}

func (cfg *apiConfig) getHits() int {
	return int(cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetHits() {
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) setPlatform(platform string) {
	cfg.platform = platform
}

func (cfg *apiConfig) getPlatform() string{
	return cfg.platform
}