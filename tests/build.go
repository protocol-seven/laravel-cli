// go:build tests
//go:build tests
// +build tests

package tests

// This file enables the tests package to be built with the "tests" build tag.
// This allows separation of test code from main application code.
//
// To run tests with this build tag:
//   go test -tags tests ./tests/...
//
// To run without integration tests:
//   go test -tags tests ./tests/ -run "Test.*" -skip "Integration"
