package system

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	projectservices "codedock/internal/services/projects"
)

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

type EdgeAuthHandler struct {
	service *projectservices.RouteProtectionService
	limiter sync.Mutex
	buckets map[string]*rateLimitEntry
}

func NewEdgeAuthHandler(service *projectservices.RouteProtectionService) *EdgeAuthHandler {
	return &EdgeAuthHandler{
		service: service,
		buckets: make(map[string]*rateLimitEntry),
	}
}

func (h *EdgeAuthHandler) checkRateLimit(key string, limit int) bool {
	h.limiter.Lock()
	defer h.limiter.Unlock()

	now := time.Now()
	entry, exists := h.buckets[key]
	if !exists || now.After(entry.resetTime) {
		h.buckets[key] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(time.Minute),
		}
		return true
	}

	if entry.count >= limit {
		return false
	}
	entry.count++
	return true
}

func extractHostAndURI(c echo.Context) (string, string, string) {
	host := c.Request().Header.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Request().Host
	}
	if colon := strings.Index(host, ":"); colon != -1 {
		host = host[:colon]
	}

	uri := c.Request().Header.Get("X-Forwarded-Uri")
	if uri == "" {
		uri = c.Request().RequestURI
	}
	if uri == "" {
		uri = "/"
	}

	method := c.Request().Header.Get("X-Forwarded-Method")
	if method == "" {
		method = c.Request().Method
	}
	return strings.ToLower(host), uri, method
}

func (h *EdgeAuthHandler) Verify(c echo.Context) error {
	host, uri, method := extractHostAndURI(c)

	if strings.HasPrefix(uri, "/.well-known/acme-challenge/") || strings.HasPrefix(uri, "/_codedock/access") {
		return c.NoContent(http.StatusOK)
	}

	_, policy, err := h.service.GetPolicyByHostname(c.Request().Context(), host)
	if err != nil || policy == nil {
		return c.NoContent(http.StatusOK)
	}

	if cookie, err := c.Cookie(projectservices.RouteAuthCookieName); err == nil && cookie != nil {
		if projectservices.VerifySessionToken(policy, host, cookie.Value) {
			c.Response().Header().Set("Cache-Control", "private, no-store")
			c.Response().Header().Set("X-Codedock-Authorized", policy.Version)
			return c.NoContent(http.StatusOK)
		}
	}

	accept := c.Request().Header.Get("Accept")
	if (method == "GET" || method == "HEAD") && strings.Contains(accept, "text/html") {
		return h.renderLoginForm(c, policy, host, uri, http.StatusUnauthorized, "")
	}

	return c.String(http.StatusUnauthorized, "This site requires a password. Open it in a browser to sign in.\n")
}

func (h *EdgeAuthHandler) LoginForm(c echo.Context) error {
	host, uri, _ := extractHostAndURI(c)
	_, policy, err := h.service.GetPolicyByHostname(c.Request().Context(), host)
	if err != nil || policy == nil {
		return c.Redirect(http.StatusFound, "/")
	}
	return h.renderLoginForm(c, policy, host, uri, http.StatusOK, "")
}

func (h *EdgeAuthHandler) Login(c echo.Context) error {
	host, _, _ := extractHostAndURI(c)
	_, policy, err := h.service.GetPolicyByHostname(c.Request().Context(), host)
	if err != nil || policy == nil {
		return c.Redirect(http.StatusFound, "/")
	}

	ip := c.RealIP()
	if !h.checkRateLimit(fmt.Sprintf("ip:%s", ip), 5) || !h.checkRateLimit(fmt.Sprintf("host:%s", host), 30) {
		c.Response().Header().Set("Retry-After", "60")
		return c.String(http.StatusTooManyRequests, "Too many sign-in attempts. Try again in a minute.\n")
	}

	returnPath := c.FormValue("return")
	if returnPath == "" || !strings.HasPrefix(returnPath, "/") || strings.Contains(returnPath, "\\") {
		returnPath = "/"
	}
	csrf := c.FormValue("csrf")
	password := c.FormValue("password")

	if !projectservices.VerifyCSRFToken(policy, host, returnPath, csrf) {
		return c.String(http.StatusBadRequest, "The sign-in form expired. Please reload the page.\n")
	}

	if !projectservices.VerifyPassword(password, policy) {
		return h.renderLoginForm(c, policy, host, returnPath, http.StatusUnauthorized, "That password is incorrect. Please try again.")
	}

	expires := time.Now().Unix() + int64(policy.SessionSeconds)
	token := projectservices.SignSessionToken(policy, host, expires)

	isSecure := c.Scheme() == "https" || c.Request().Header.Get("X-Forwarded-Proto") == "https"
	cookie := &http.Cookie{
		Name:     projectservices.RouteAuthCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   policy.SessionSeconds,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, returnPath)
}

func (h *EdgeAuthHandler) renderLoginForm(c echo.Context, policy *models.RouteProtectionPolicy, host, path string, status int, errMsg string) error {
	if path == "" || !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "/_codedock/access") {
		path = "/"
	}
	expires := time.Now().Unix() + 600
	csrf := projectservices.SignCSRFToken(policy, host, path, expires)

	errorHTML := ""
	if errMsg != "" {
		errorHTML = fmt.Sprintf(`<p class="error" role="alert">%s</p>`, html.EscapeString(errMsg))
	}

	doc := fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Private site</title>
<style>
*{box-sizing:border-box}
body{margin:0;min-height:100vh;display:grid;place-items:center;background:#09090b;color:#f4f4f5;font:15px -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif}
main{width:min(400px,calc(100%% - 32px));padding:32px;border-radius:20px;background:#18181b;border:1px solid #27272a;box-shadow:0 12px 32px rgba(0,0,0,0.4)}
.icon-box{display:inline-flex;align-items:center;justify-content:center;width:44px;height:44px;border-radius:12px;background:#27272a;margin-bottom:16px}
svg{width:22px;height:22px;stroke:#a1a1aa}
h1{font-size:20px;font-weight:600;margin:0 0 6px;letter-spacing:-0.02em}
p{color:#a1a1aa;margin:0 0 20px;font-size:13px;line-height:1.5}
label{display:block;margin:0 0 8px;font-size:13px;font-weight:500;color:#e4e4e7}
input[type="password"]{width:100%%;font:inherit;padding:12px 14px;border-radius:12px;background:#09090b;border:1px solid #3f3f46;color:#fff;outline:none;transition:border-color 0.15s}
input[type="password"]:focus{border-color:#38bdf8}
button{width:100%%;font:inherit;font-weight:500;margin-top:16px;padding:12px;border-radius:12px;background:#fff;color:#09090b;border:0;cursor:pointer;transition:opacity 0.15s}
button:hover{opacity:0.9}
.error{color:#f87171;background:rgba(239,68,68,0.1);padding:10px 12px;border-radius:8px;font-size:12px;margin:0 0 16px}
footer{display:flex;align-items:center;justify-content:center;gap:6px;margin-top:24px;font-size:11px;color:#71717a}
</style>
</head>
<body>
<main>
<div class="icon-box">
<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
</svg>
</div>
<h1>This site is private</h1>
<p>Enter the password shared by the site owner to continue.</p>
%s
<form method="post" action="/_codedock/access">
<input type="hidden" name="return" value="%s">
<input type="hidden" name="csrf" value="%s">
<label for="password">Password</label>
<input id="password" name="password" type="password" required maxlength="256" autofocus autocomplete="current-password" placeholder="Enter password">
<button type="submit">Continue</button>
</form>
<footer>Protected by Codedock</footer>
</main>
</body>
</html>`, errorHTML, html.EscapeString(path), html.EscapeString(csrf))

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response().Header().Set("X-Frame-Options", "DENY")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	c.Response().Header().Set("Referrer-Policy", "no-referrer")
	return c.HTML(status, doc)
}
