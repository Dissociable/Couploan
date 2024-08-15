package regex

type Options struct {
	// ManipulateStringResult allows you to manipulate the result of a string
	//
	// NOTE: result may be a nil in case of no match, so take care of that too.
	ManipulateStringResult func(result *string) *string
}

type Regex struct {
	options Options
}

func New(options *Options) *Regex {
	if options == nil {
		options = &Options{}
	}
	return &Regex{
		options: *options,
	}
}

// Must returns a string
func Must(r *string) string {
	if r == nil {
		return ""
	}
	return *r
}
