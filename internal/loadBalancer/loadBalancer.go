package loadBalancer

import (
    "errors"
    "log/slog"
    "net/http"
    "net/url"
)

// Iterator определяет алгоритм распределения запросов по бекендам и способ хранения бекендов
type Iterator interface {
    AddBackendToPool(log *slog.Logger, url *url.URL)
    Next(log *slog.Logger) (*http.HandlerFunc, error)
}

var (
    ErrNoBackends = errors.New("no backends in pool")
)

type LoadBalancer struct {
    iterator Iterator
    log      *slog.Logger
}

// New создает новый балансировщик с указанным итератором
func New(log *slog.Logger, iterator Iterator) *LoadBalancer {
    log.Debug("load balancer created")
    return &LoadBalancer{
        iterator: iterator,
        log:      log,
    }
}

// TryToRegisterNewBackend пробует зарегистрировать бекенд в балансировщике.
// Если указанный URL не валиден, то бекенд не зарегистрируется.
func (lb *LoadBalancer) TryToRegisterNewBackend(rawUrl string) {
    parsedUrl, err := url.Parse(rawUrl)
    if err != nil {
        lb.log.Error("backend invalid url: %s", err.Error())
        return
    }
    lb.log.Debug("backend url parsed", slog.String("url", parsedUrl.String()))

    lb.iterator.AddBackendToPool(lb.log, parsedUrl)

    lb.log.Info("new backend registered", slog.String("url", parsedUrl.String()))
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log := lb.log.With(slog.String("remote_addr", r.RemoteAddr))

    hf, err := lb.iterator.Next(log)
    if err != nil {
        log.Error("cannot get next backend", slog.String("error", err.Error()))
        http.Error(w, "please try again later", http.StatusServiceUnavailable)
        return
    }

    (*hf).ServeHTTP(w, r)
}
