package ingress

import "errors"

var (
	ErrConnectionNotFound = errors.New("ingress connection not found")
	ErrHostnameNotFound   = errors.New("ingress hostname not found")
	ErrCloudflareAPI      = errors.New("cloudflare ingress api error")
	ErrInvalidConfig      = errors.New("invalid ingress configuration")
)
