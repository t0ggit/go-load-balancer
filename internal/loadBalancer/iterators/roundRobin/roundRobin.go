package roundRobin

import (
    "go-load-balancer/internal/loadBalancer"
    "log/slog"
    "net/http"
    "net/url"
    "sync"
)

type Backend interface {
    SetURL(log *slog.Logger, url *url.URL)

    HandlerFunc(log *slog.Logger) http.HandlerFunc
}

type RoundRobin struct {
    pool  []Backend
    index int
    mux   *sync.Mutex
}

// New создает новый итератор со стратегией Round Robin для балансировщика
func New() *RoundRobin {
    return &RoundRobin{
        pool:  make([]Backend, 0),
        index: 0,
        mux:   &sync.Mutex{},
    }
}

// AddBackendToPool добавляет бекенд в пул бекендов итератора
func (r *RoundRobin) AddBackendToPool(log *slog.Logger, url *url.URL) {
    var newBackend Backend
    newBackend = &loadBalancer.Backend{}
    newBackend.SetURL(log, url)

    log.Debug("trying to lock mutex for adding backend to pool")
    defer log.Debug("backend added to pool, mutex unlocked")

    r.mux.Lock()
    defer r.mux.Unlock()
    r.pool = append(r.pool, newBackend)
}

func (r *RoundRobin) Next(log *slog.Logger) (*http.HandlerFunc, error) {
    log.Debug("trying to lock mutex for getting next backend")
    defer log.Debug("got next backend, mutex unlocked")

    r.mux.Lock()
    defer r.mux.Unlock()

    if len(r.pool) == 0 {
        return nil, loadBalancer.ErrNoBackends
    }

    r.index = (r.index + 1) % len(r.pool)
    hf := r.pool[r.index].HandlerFunc(log)
    return &hf, nil
}
