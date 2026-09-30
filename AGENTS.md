# data-slice — Notes for AI Agents

- **`ApplyOptions` mutates its input and returns a sub-slice of the same backing array.** Sorting reorders `value` in place, and `MaxRows`/`FirstRow` reslice it (`value[:n]`). A caller that needs the original order or an independent copy must clone before calling, or copy the result — the returned slice still aliases the caller's array.

- **Element types must implement `Comparer[T]`.** `ApplyOptions[T Comparer[T]]` requires each element to order itself against a peer by field name (`Compare(fieldName string, other T) int`). Without it the call won't compile.

- **Only `Sort`, `MaxRows`, and `FirstRow` options are honored.** Any other `data` option is ignored. `FirstRow` and `MaxRows` are distinct: `FirstRow` always trims to one row; `MaxRows` trims to its configured count.

## `MaxRows(0)` means NO LIMIT, and the three backends must agree

Zero is "unbounded" ecosystem-wide, and negatives are clamped to zero by both the constructor and the accessor. This package originally read `0 <= maxRows` and truncated to an empty result, while `data-mongo` already used `if opt > 0`; the test that asserted the old behavior was inverted along with the fix. Sort direction had the same split: an unrecognized `Direction` defaulted to **descending** here and **ascending** in mongo and mock. Both now route through `option.SortOption.IsDescending()` so all three answer identically. When adding an option, check it against the other two backends before assuming this one is right.
