# Production TODO

## Cookie Security

- [ ] Set `Secure: true` in `setAnonCookie()` (`internal/shorturl/handler/create_short_url.go:119`)
- [ ] Use `SameSite: Lax` if frontend/backend on same domain
- [ ] Use `SameSite: None` if frontend/backend on different domains (requires `Secure: true`)

## CORS

- [ ] Update `AllowOrigins` in `pkg/httpserver/echo_server.go` with production frontend URL
- [ ] Remove localhost/dev origins before deploying

## HTTPS

- [ ] Set up HTTPS (required for `Secure: true` cookies)
- [ ] Redirect HTTP to HTTPS
