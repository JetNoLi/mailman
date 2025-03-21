package common

//TODO: add error

type OptFns[optType any] func(args *optType)

func CreateOptions[optType any](baseOptions *optType, optFns ...OptFns[optType]) (options *optType) {
	options = baseOptions

	if options == nil {
		options = new(optType)
	}

	for _, fn := range optFns {
		fn(options)
	}

	return options
}
