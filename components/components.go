// Code generated by gtml; DO NOT EDIT.
// +build ignore
// v0.1.0 | you may see errors with types, you'll need to manage your own imports
// type support coming soon!

package components

import "strings"

func gtmlFor[T any](slice []T, callback func(i int, item T) string) string {
	var builder strings.Builder
	for i, item := range slice {
		builder.WriteString(callback(i, item))
	}
	return builder.String()
}

func gtmlIf(condition bool, fn func() string) string {
	if condition {
		return fn()
	}
	return ""
}

func gtmlElse(condition bool, fn func() string) string {
	if !condition {
		return fn()
	}
	return ""
}

func gtmlSlot(contentFunc func() string) string {
	return contentFunc()
}

