package benchmark

import (
	"net/http"
	"net/url"
)

type CookieJarImpl struct {
	cookies map[string][]*http.Cookie
}

func (cookieJar *CookieJarImpl) SetCookies(u *url.URL, cookies []*http.Cookie) {
	cookieJar.cookies[u.Host] = cookies
}

func (cookieJar *CookieJarImpl) Cookies(u *url.URL) []*http.Cookie {
	return cookieJar.cookies[u.Host]
}
