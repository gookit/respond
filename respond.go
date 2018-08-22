package respond

// create a default instance
var Er = New()

// Initialize the default instance
func Initialize(config func(opts *Options)) {
	// config options
	config(Er.opts)

	// init
	Er.Initialize()
}
