package action

import (
	"bytes"
	"crypto/subtle"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// The basic proxy type. Implements http.Handler.
type ProxyServer struct {
	Config *Config
}

func newProxyServer(config *Config) *ProxyServer {
	proxy := ProxyServer{
		Config: config,
	}

	return &proxy
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	heartbeatURL := p.getHeartbeatUrl(r.URL.Path)
	if heartbeatURL == nil {
		matchingURL := p.getMatchingURL(r.URL.Path)
		if matchingURL == nil {
			blockRequest(w, r)
			return
		} else {
			if isMethodAllowed(r.Method, matchingURL.Methods) {
				if !isBasicAuthenticationAllowed(matchingURL.BasicAuthentication, r) {
					authBlockRequest(w, r)
					return
				} else {
					if err := validateRules(matchingURL.Rules, r); err != nil {
						blockRequest(w, r)
						return
					} else {
						if err := validateHeaders(matchingURL.Headers, r); err != nil {
							if matchingURL.Debug != nil && matchingURL.Debug.Headers != nil && *matchingURL.Debug.Headers == true {
								log.Println(err)
							}
							blockRequest(w, r)
							return
						} else {
							serveReverseProxy(fmt.Sprintf("%s%s", matchingURL.BackendHost, matchingURL.Path), w, r, matchingURL.Debug)
							return
						}
					}
				}
			} else {
				blockRequest(w, r)
				return
			}
		}
	} else {
		heartbeatRequest(w, r, heartbeatURL)
		return
	}
}

func validateRules(rules *ConfigRouteRules, req *http.Request) error {
	if rules == nil {
		return nil
	} else {
		if rules.HasBody != nil {
			if *rules.HasBody == true && req.Body == http.NoBody {
				return ErrRuleHasNoBody
			} else if *rules.HasBody == false && req.Body != http.NoBody {
				return ErrRuleHasBody
			}
		}
		if rules.HasQueryString != nil {
			if *rules.HasQueryString == true && len(req.URL.Query()) == 0 {
				return ErrRuleHasNoQueryStrings
			} else if *rules.HasQueryString == false && len(req.URL.Query()) > 0 {
				return ErrRuleHasQueryStrings
			}
		}
		return nil
	}
}

func validateHeaders(headers *[]string, req *http.Request) error {
	if headers == nil {
		return nil
	} else {
		for _, header := range *headers {
			if _, ok := req.Header[header]; !ok {
				return errors.New(fmt.Sprintf("missing header: %s", header))
			}
		}
		return nil
	}
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request, debug *ConfigRouteDebug) {
	logRequest(LogStatusPassed, req, target)

	if debug != nil {
		debugRequest(req, debug)
	}

	targetURL, _ := url.Parse(target)

	targetURL.Path = ""

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	req.URL.Host = targetURL.Host
	req.URL.Scheme = targetURL.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Header.Set("X-Api-Gateway", "Mekong")
	req.Host = targetURL.Host

	proxy.ServeHTTP(res, req)
}

func debugRequest(r *http.Request, debug *ConfigRouteDebug) {
	if debug != nil {
		if debug.URL != nil && *debug.URL == true {
			log.Println(r.URL)
		}
		if debug.Headers != nil && *debug.Headers == true {
			log.Println(r.Header)
		}
		if debug.Body != nil && *debug.Body == true {
			buf, err := ioutil.ReadAll(r.Body)
			if err == nil {
				log.Println(ioutil.NopCloser(bytes.NewBuffer(buf)))
			}
		}
	}
}

type LogStatus string

const (
	LogStatusAuthRequired LogStatus = "AUTH_REQUIRED"
	LogStatusPassed       LogStatus = "PASSED"
	LogStatusBlocked      LogStatus = "BLOCKED"
	LogStatusHeartbeat    LogStatus = "HEARTBEAT"
)

func logRequest(statusCode LogStatus, r *http.Request, b string) {
	switch statusCode {
	case LogStatusPassed:
		log.Println(fmt.Sprintf("%s |   \u001b[42;30m   %s  \u001B[0m   | %s | \u001B[42;44m %s \u001B[0m | %s -> %s", "[MEKONG]", statusCode+" 🚀", r.Host, r.Method, r.URL, b))
	case LogStatusAuthRequired:
		log.Println(fmt.Sprintf("%s | \u001b[43m %s \u001B[0m | %s | \u001B[42;44m %s \u001B[0m | %s", "[MEKONG]", statusCode+" 🛑", r.Host, r.Method, r.URL))
	case LogStatusHeartbeat:
		log.Println(fmt.Sprintf("%s |  \u001b[47m  %s  \u001B[0m  | %s | \u001B[42;44m %s \u001B[0m | %s", "[MEKONG]", statusCode+" 💓", r.Host, r.Method, r.URL))
	default:
		log.Println(fmt.Sprintf("%s |   \u001b[42;41m  %s  \u001B[0m   | %s | \u001B[42;44m %s \u001B[0m | %s", "[MEKONG]", statusCode+" 🚫", r.Host, r.Method, r.URL))
	}
}

func (p *ProxyServer) getMatchingURL(path string) *ConfigRoutes {
	for _, elem := range p.Config.Routes {
		if matchSplitURL(strings.Split(path, "/"), strings.Split(elem.Path, "/")) {
			return &elem
		}
	}
	return nil
}

func (p *ProxyServer) getHeartbeatUrl(path string) *ConfigHeartbeat {
	if p.Config.Heartbeat == nil {
		return nil
	}
	for _, elem := range *p.Config.Heartbeat {
		if matchSplitURL(strings.Split(path, "/"), strings.Split(elem.Path, "/")) {
			return &elem
		}
	}
	return nil
}

func matchSplitURL(incomingURLSplit []string, targetURLSplit []string) bool {
	matching := 0
	if len(incomingURLSplit) == len(targetURLSplit) {
		for index, elem := range targetURLSplit {
			// Check if the element is a wildcard indication aka *
			if elem == "*" {
				matching++
			} else if incomingURLSplit[index] == elem {
				matching++
			}
		}
		return len(targetURLSplit) == matching
	} else {
		return false
	}
}

func isMethodAllowed(method string, allowedMethods []HTTPMethod) bool {
	for _, elem := range allowedMethods {
		if elem.String() == method {
			return true
		}
	}
	return false
}

func isBasicAuthenticationAllowed(basicAuth *ConfigRouteBasicAuth, r *http.Request) bool {
	if basicAuth != nil {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(basicAuth.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(basicAuth.Password)) != 1 {
			return false
		}
		return true
	}
	return true
}

func blockRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(LogStatusBlocked, r, "")
	w.WriteHeader(500)
	w.Write([]byte("Access denied"))
}

func authBlockRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(LogStatusAuthRequired, r, "")
	w.WriteHeader(401)
	w.Write([]byte("Unauthorised"))
	return
}

func heartbeatRequest(w http.ResponseWriter, r *http.Request, ch *ConfigHeartbeat) {
	logRequest(LogStatusHeartbeat, r, "")
	if ch.ResponseCode == nil {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(*ch.ResponseCode)
	}
	if ch.Message != nil {
		w.Write([]byte(*ch.Message))
	}
	return
}
