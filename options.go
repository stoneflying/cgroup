package cgroup

import (
	"log"
	"os"
)

var (
	// defaultLogger is a default Logger instance that logs to standard error with date and time information.
	defaultLogger = Logger(log.New(os.Stderr, "", log.LstdFlags))
)

// Option is a function type used to modify the configuration options of a CGroup.
type Option func(opts *Options)

type Options struct {
	Logger       Logger            // Logger is used to log information.
	PanicHandler func(interface{}) // PanicHandler is used to handle panic errors.
}

// loadOptions load the configuration options and set the default options.
func loadOptions(options ...Option) *Options {
	opts := loadCustomOptions(options...)
	loadDefaultOptions(opts)
	return opts
}

// loadCustomOptions load custom options.
func loadCustomOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// loadDefaultOptions load default options.
func loadDefaultOptions(opts *Options) {
	// If Logger is not included in the options, the default Logger instance is set as the Logger option.
	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}
}

// WithPanicHandler set the PanicHandler option.
func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

// WithLogger set the Logger option.
func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}
