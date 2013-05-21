package benchmark

import (
	"net/http"
	"net/url"
	"sync"
)

type OperationHelper struct {
}

func NewOperationHelper() OperationHelper {
	return OperationHelper{}
}

func (o *OperationHelper) Prepare() {

}

func (o *OperationHelper) Cleanup() {

}

// logs in url and returns client that holds cookies
func (o *OperationHelper) Login(address, usernameField, username, passwordField, password string) http.Client {
	cookieJar := new(CookieJarImpl)
	client := http.Client{nil, nil, cookieJar}

	values := url.Values{}
	values.Set(usernameField, username)
	values.Set(passwordField, password)

	resp, _ := client.PostForm(address, values)
	resp.Body.Close()

	return client
}

func (o *OperationHelper) Logout() {

}

type CookieJarImpl struct {
	lock    sync.Mutex
	cookies map[string][]*http.Cookie
}

func (cookieJar *CookieJarImpl) SetCookies(u *url.URL, cookies []*http.Cookie) {
	cookieJar.lock.Lock()
	cookieJar.cookies[u.Host] = cookies
	cookieJar.lock.Unlock()
}

func (cookieJar *CookieJarImpl) Cookies(u *url.URL) []*http.Cookie {
	return cookieJar.cookies[u.Host]
}
