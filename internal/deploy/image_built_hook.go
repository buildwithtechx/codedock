package deploy

import "context"

type imageBuiltKey struct{}

// WithImageBuiltCallback attaches an image built notification callback to the context.
func WithImageBuiltCallback(ctx context.Context, cb func(tag string)) context.Context {
	return context.WithValue(ctx, imageBuiltKey{}, cb)
}

// NotifyImageBuilt triggers the image built notification callback if registered.
func NotifyImageBuilt(ctx context.Context, tag string) {
	if cb, ok := ctx.Value(imageBuiltKey{}).(func(string)); ok && cb != nil {
		cb(tag)
	}
}
