//go:build !darwin

package main

// configureOpenMode is a no-op outside macOS, where hidapi already opens
// devices non-exclusively and each USB interface has its own device path.
func configureOpenMode() {}
