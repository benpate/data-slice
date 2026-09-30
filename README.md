# 🔪 data-slice

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/data-slice.svg)](https://pkg.go.dev/github.com/benpate/data-slice)
[![Version](https://img.shields.io/github/v/release/benpate/data-slice?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/data-slice/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/data-slice/go.yml?branch=main)](https://github.com/benpate/data-slice/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/data-slice?style=flat-square)](https://goreportcard.com/report/github.com/benpate/data-slice)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/data-slice.svg?style=flat-square)](https://codecov.io/gh/benpate/data-slice)

## Apply Query Options to In-Memory Slices

`data-slice` applies [`benpate/data`](https://github.com/benpate/data) query options — sorting and row limits — to plain Go slices held in memory. Its single generic function, `ApplyOptions`, lets in-memory data honor the same `option.SortAsc`, `option.MaxRows`, and `option.FirstRow` arguments that the database adapters use.

### Using data-slice

```go
// Your slice element type implements Comparer, so the package knows how to order it by field name.
type Person struct {
    Name string
}

func (p Person) Compare(fieldName string, other Person) int {
    return strings.Compare(p.Name, other.Name)
}

// Sort by name ascending, then keep only the first 10 rows.
people = dataslice.ApplyOptions(people, option.SortAsc("name"), option.MaxRows(10))
```

## Pull Requests Welcome

I'm trying to make data-slice the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🔪
