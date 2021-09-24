package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const cookieName = "proxykey"

func UrlInspector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		key := GetKeyValue(req)

		if proxy, ok := siteKeys[key]; ok {

			cookie := http.Cookie{
				Name:     cookieName,
				Value:    key,
				HttpOnly: true,
				MaxAge:   30,
			}
			http.SetCookie(resp, &cookie)
			resp.Header().Add("proxykey", key)
			resp.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", 5))

			proxy.ServeHTTP(resp, req)
			return
		}

		fmt.Println("Key not found: ", key)

		for k := range siteKeys {
			fmt.Println("keys we got: ", k)
		}

		http.Error(resp, "invalid key", http.StatusNotFound)
	})
}

func GetKeyValue(req *http.Request) string {
	key := req.URL.Query().Get("key")
	if key == "" {
		if cookieKey, err := req.Cookie(cookieName); err == nil {
			key = cookieKey.Value
		}
	}

	return key
}

func CreateReverseProxy(urlKey string) *httputil.ReverseProxy {

	urlParsed, _ := url.Parse(urlKey)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: urlParsed.Scheme,
		Host:   urlParsed.Host,
	})

	proxy.Director = func(r *http.Request) {

		r.Host = urlParsed.Host
		r.URL.Scheme = urlParsed.Scheme
		r.URL.Host = urlParsed.Host + r.URL.Host
		r.URL.Path = urlParsed.Path + r.URL.Path
	}

	return proxy
}
