package rtr

import "ewails/wrappers/common"

type HandlerOptions struct {
	middleware []MiddlewareFn
}

type HandlerOptsFn = common.OptFns[HandlerOptions]

func WithMiddleware(middlewareFns ...MiddlewareFn) HandlerOptsFn {
	return func(options *HandlerOptions) {
		options.middleware = middlewareFns
	}
}
