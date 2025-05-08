package loadBalancer

import (
    "log/slog"
    "net/http"
    "net/http/httputil"
    "net/url"
)

type Backend struct {
    Url          *url.URL
    ReverseProxy *httputil.ReverseProxy
}

// SetURL привязывает бекенд к URL
func (b *Backend) SetURL(log *slog.Logger, url *url.URL) {
    b.Url = url

    proxy := httputil.NewSingleHostReverseProxy(b.Url)
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        log.Error("backend unavailable",
            slog.String("backend_url", b.Url.String()),
            slog.String("error", err.Error()))

        http.Error(w, "please try again", http.StatusBadGateway)
    }

    b.ReverseProxy = proxy

    log.Debug("backend url set", slog.String("url", b.Url.String()))
}

// HandlerFunc возвращает обработчик запроса для Reverse Proxy к соответствующему бекенду
func (b *Backend) HandlerFunc(log *slog.Logger) http.HandlerFunc {
    log.Info("balanced to backend", slog.String("backend_url", b.Url.String()))
    return b.ReverseProxy.ServeHTTP
}
