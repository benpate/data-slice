package dataslice

import (
	"strings"
	"testing"

	"github.com/benpate/data/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testRecord is a simple Comparer used throughout these tests. It can be
// sorted on either the "name" or "rank" field.
type testRecord struct {
	Name string
	Rank int
}

// Compare implements the Comparer interface. Unknown field names compare as
// equal (return 0), which mirrors the behavior of a query against a missing
// column.
func (record testRecord) Compare(fieldName string, other testRecord) int {
	switch fieldName {
	case "name":
		return strings.Compare(record.Name, other.Name)
	case "rank":
		return record.Rank - other.Rank
	default:
		return 0
	}
}

// names is a small helper that extracts the Name field from a slice of records
// so that ordering assertions stay readable.
func names(records []testRecord) []string {
	result := make([]string, len(records))
	for index, record := range records {
		result[index] = record.Name
	}
	return result
}

func TestApplyOptions_NoOptions(t *testing.T) {

	// With no options the slice should be returned untouched, in its original order.
	value := []testRecord{
		{Name: "charlie", Rank: 3},
		{Name: "alice", Rank: 1},
		{Name: "bob", Rank: 2},
	}

	result := ApplyOptions(value)

	require.Len(t, result, 3)
	assert.Equal(t, []string{"charlie", "alice", "bob"}, names(result))
}

func TestApplyOptions_EmptySlice(t *testing.T) {

	// Applying every kind of option to an empty slice must not panic and must
	// return an empty slice.
	value := []testRecord{}

	result := ApplyOptions(value,
		option.SortAsc("name"),
		option.MaxRows(10),
		option.FirstRow(),
	)

	assert.Empty(t, result)
}

func TestApplyOptions_SortAscending(t *testing.T) {

	value := []testRecord{
		{Name: "charlie", Rank: 3},
		{Name: "alice", Rank: 1},
		{Name: "bob", Rank: 2},
	}

	result := ApplyOptions(value, option.SortAsc("name"))

	require.Len(t, result, 3)
	assert.Equal(t, []string{"alice", "bob", "charlie"}, names(result))
}

func TestApplyOptions_SortDescending(t *testing.T) {

	value := []testRecord{
		{Name: "charlie", Rank: 3},
		{Name: "alice", Rank: 1},
		{Name: "bob", Rank: 2},
	}

	result := ApplyOptions(value, option.SortDesc("name"))

	require.Len(t, result, 3)
	assert.Equal(t, []string{"charlie", "bob", "alice"}, names(result))
}

func TestApplyOptions_SortByDifferentField(t *testing.T) {

	// Sorting on "rank" must order by the integer field, independent of Name.
	value := []testRecord{
		{Name: "alice", Rank: 3},
		{Name: "bob", Rank: 1},
		{Name: "charlie", Rank: 2},
	}

	result := ApplyOptions(value, option.SortAsc("rank"))

	require.Len(t, result, 3)
	assert.Equal(t, []string{"bob", "charlie", "alice"}, names(result))
}

func TestApplyOptions_OnlyFirstSortApplied(t *testing.T) {

	// When two sort options are present, only the first one should take effect
	// because ApplyOptions breaks out of the loop after the first match.
	value := []testRecord{
		{Name: "charlie", Rank: 1},
		{Name: "alice", Rank: 3},
		{Name: "bob", Rank: 2},
	}

	// First sort is by name ascending; the second (rank) should be ignored.
	result := ApplyOptions(value, option.SortAsc("name"), option.SortAsc("rank"))

	require.Len(t, result, 3)
	assert.Equal(t, []string{"alice", "bob", "charlie"}, names(result))
}

func TestApplyOptions_MaxRowsTruncates(t *testing.T) {

	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
		{Name: "charlie"},
	}

	result := ApplyOptions(value, option.MaxRows(2))

	require.Len(t, result, 2)
	assert.Equal(t, []string{"alice", "bob"}, names(result))
}

func TestApplyOptions_MaxRowsEqualToLength(t *testing.T) {

	// When MaxRows equals the slice length the slice should be returned in full
	// (the truncation guard is `maxRows < len`).
	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
	}

	result := ApplyOptions(value, option.MaxRows(2))

	require.Len(t, result, 2)
	assert.Equal(t, []string{"alice", "bob"}, names(result))
}

func TestApplyOptions_MaxRowsGreaterThanLength(t *testing.T) {

	// MaxRows larger than the slice is a no-op.
	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
	}

	result := ApplyOptions(value, option.MaxRows(100))

	require.Len(t, result, 2)
	assert.Equal(t, []string{"alice", "bob"}, names(result))
}

func TestApplyOptions_MaxRowsZero(t *testing.T) {

	// MaxRows of zero means "no limit", matching data-mongo and the zero value
	// of MaxRowsOption. Use FirstRow to ask for a single row instead.
	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
	}

	result := ApplyOptions(value, option.MaxRows(0))

	require.Len(t, result, 2)
	assert.Equal(t, []string{"alice", "bob"}, names(result))
}

func TestApplyOptions_MaxRowsNegative(t *testing.T) {

	// A negative MaxRows is clamped to zero by the option package, so it reads
	// as "no limit" rather than panicking on a `value[:negative]` reslice.
	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
	}

	result := ApplyOptions(value, option.MaxRows(-1))

	require.Len(t, result, 2)
	assert.Equal(t, []string{"alice", "bob"}, names(result))
}

func TestApplyOptions_FirstRow(t *testing.T) {

	// FirstRow keeps only the first record in the slice.
	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
		{Name: "charlie"},
	}

	result := ApplyOptions(value, option.FirstRow())

	require.Len(t, result, 1)
	assert.Equal(t, []string{"alice"}, names(result))
}

func TestApplyOptions_FirstRowOnEmptySlice(t *testing.T) {

	// FirstRow on an empty slice is guarded by a `len(value) > 0` check, so it
	// must not panic and must return an empty slice.
	value := []testRecord{}

	result := ApplyOptions(value, option.FirstRow())

	assert.Empty(t, result)
}

func TestApplyOptions_SortThenMaxRows(t *testing.T) {

	// Options are applied in order: sort first, then truncate. The result should
	// be the lowest-ranked records, in ascending order.
	value := []testRecord{
		{Name: "charlie", Rank: 3},
		{Name: "alice", Rank: 1},
		{Name: "delta", Rank: 4},
		{Name: "bob", Rank: 2},
	}

	result := ApplyOptions(value, option.SortAsc("rank"), option.MaxRows(2))

	require.Len(t, result, 2)
	assert.Equal(t, []string{"alice", "bob"}, names(result))
}

func TestApplyOptions_SortThenFirstRow(t *testing.T) {

	// Sorting descending and then taking the first row yields the single
	// highest-ranked record.
	value := []testRecord{
		{Name: "charlie", Rank: 3},
		{Name: "alice", Rank: 1},
		{Name: "delta", Rank: 4},
		{Name: "bob", Rank: 2},
	}

	result := ApplyOptions(value, option.SortDesc("rank"), option.FirstRow())

	require.Len(t, result, 1)
	assert.Equal(t, []string{"delta"}, names(result))
}

func TestApplyOptions_MaxRowsThenFirstRow(t *testing.T) {

	// MaxRows trims to two records, then FirstRow trims to one.
	value := []testRecord{
		{Name: "alice"},
		{Name: "bob"},
		{Name: "charlie"},
	}

	result := ApplyOptions(value, option.MaxRows(2), option.FirstRow())

	require.Len(t, result, 1)
	assert.Equal(t, []string{"alice"}, names(result))
}

func TestApplyOptions_UnknownFieldIsStable(t *testing.T) {

	// Sorting on an unknown field name compares everything as equal. SortFunc is
	// not guaranteed to be stable, so we only assert that the contents are
	// preserved, not their order.
	value := []testRecord{
		{Name: "charlie"},
		{Name: "alice"},
		{Name: "bob"},
	}

	result := ApplyOptions(value, option.SortAsc("missing"))

	require.Len(t, result, 3)
	assert.ElementsMatch(t, []string{"alice", "bob", "charlie"}, names(result))
}
