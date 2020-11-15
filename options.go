package echomsg

type (
	Options struct {
		code    int
		message string
		field   string
	}

	Option func(*Options)
)

func MapError(field string, message string) Option {
	return func(args *Options) {
		args.field = field
		args.message = message
	}
}

func Error(message string) Option {
	return func(args *Options) {
		args.message = message
	}
}

func Code(code int) Option {
	return func(args *Options) {
		args.code = code
	}
}
