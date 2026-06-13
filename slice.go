// Package dataslice applies query options (sorting and row limits) to in-memory slices.
package dataslice

import (
	"slices"

	"github.com/benpate/data/option"
)

// Comparer is implemented by types that can order themselves against a peer by field name.
type Comparer[T any] interface {
	Compare(fieldName string, other T) int
}

// ApplyOptions sorts and trims value according to the supplied options. It mutates value in
// place and returns a sub-slice of the same backing array, so callers needing an independent
// copy must make one.
func ApplyOptions[T Comparer[T]](value []T, options ...option.Option) []T {

	// Try to sort the items in the slice based on the sort criteria
	for _, opt := range options {
		if typed, ok := opt.(option.SortOption); ok {

			if typed.Direction == option.SortDirectionAscending {
				slices.SortFunc(value, func(a T, b T) int {
					return a.Compare(typed.FieldName, b)
				})
				break
			}

			slices.SortFunc(value, func(a T, b T) int {
				return b.Compare(typed.FieldName, a)
			})
			break
		}
	}

	// Apply MaxRows option (if present)
	for _, opt := range options {
		if typed, ok := opt.(option.MaxRowsOption); ok {
			if maxRows := typed.MaxRows(); 0 <= maxRows && maxRows < int64(len(value)) {
				value = value[:maxRows]
			}
			break
		}
	}

	// Apply FirstRow option (if present)
	for _, opt := range options {
		if _, ok := opt.(option.FirstRowOption); ok {
			if len(value) > 0 {
				value = value[:1]
			}
			break
		}
	}

	return value
}
