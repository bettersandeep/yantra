# Report: review round 1, ticket #13

## 1. What changed

Only `src/select.rs`. No behaviour change in the engine: both findings were about missing evidence, and the measurements confirmed the existing code.

- `Node` doc comment: records the SQLite 3.46.1 results for `<=`, `>`, `>=`, `==`, `!=` on both sides, plus the right side of `<` and `<>`.
- `MAX_DEPTH` comment: the parenthesis claim now cites the measured query and result (1001 and 1002 pairs).
- New test `height_of_each_comparison_operator`: all eight comparison operators, deep operand on each side, height 999 works and height 1000 is an error.

`yt check` passed: fmt and clippy clean, 48 unit tests and 10 harness tests ok, `slt: random/expr/slt_good_0.test 9644/10012`. The baseline is unchanged.

## 2. Acceptance criteria

Only the height-rule criterion was touched this round. The others were not re-measured individually; the existing tests for them ran inside `yt check` and passed.

| Criterion | Query or check | Result seen |
|---|---|---|
| 2000 columns work, 2001 is an error | `column_limit` in `yt check` | ok |
| 1000 arguments work, 1001 is an error | `function_argument_limit` in `yt check` | ok |
| Height computed as SQLite does, above 1000 is an error | `height_of_each_expression`, `height_of_each_comparison_operator` in `yt check` | ok |
| Sum of 1000 / 1001 terms, `COALESCE(NULL, <sum>)` with 999 / 1000 terms, 999 unary minus | `height_limit_of_the_ticket` in `yt check` | ok. The 999 unary minus signs are still stopped by the engine's own limit of 100 |
| Height rule for comparisons found with `yt sqlite` | `yt sqlite .yt/q.sql`, table below | Highest operand plus 1 for every comparison operator, both sides |
| Each limit names its source | Comments at `MAX_COLUMNS`, `MAX_FUNCTION_ARGS`, `MAX_EXPR_HEIGHT`, `MAX_DEPTH` | Present; `MAX_DEPTH` now cites its measurement |
| Baseline is exact and does not go down | `baseline_must_match_exactly` in `yt check` | ok, 9644 |
| No record that passed before fails now | `slt` in `yt check` | 9644/10012, same as before |

Measurements on SQLite 3.46.1. `H` is 999 `~` before `1` (height 1000), `G` is 998 `~` before `1` (height 999). I checked the tilde and parenthesis counts in the file with a regex before running it.

| Query | Result |
|---|---|
| `SELECT H` | -2 |
| `SELECT ~ H` | Expression tree is too large (maximum depth 1000) |
| `H <= 1`, `1 <= H` | error, error |
| `G <= 1`, `1 <= G` | 1, 1 |
| `H > 1`, `1 > H` | error, error |
| `G > 1`, `1 > G` | 0, 0 |
| `H >= 1`, `1 >= H` | error, error |
| `G >= 1`, `1 >= G` | 1, 1 |
| `H == 1`, `1 == H` | error, error |
| `G == 1`, `1 == G` | 1, 1 |
| `H != 1`, `1 != H` | error, error |
| `G != 1`, `1 != G` | 0, 0 |
| `1 < H`, `1 < G` | error, 0 |
| `1 <> H`, `1 <> G` | error, 0 |
| `SELECT (((7)))` with 1001 pairs of parentheses | 7 |
| The same with 1002 pairs | 7 |

"error" is the same "Expression tree is too large" text in every row.

## 3. Known differences from SQLite

Unchanged from the first round:

- The engine's nesting limit of 100 stops parentheses and unary operators that SQLite accepts (1001 pairs of parentheses, 999 unary minus signs). The engine gives an error, not a wrong value.
- The argument limit is 1000 from the document; SQLite 3.46.1 stops at 127.
- The engine computes every `COALESCE` argument and every `CASE` part, so an error in a part SQLite skips is an error here.

## 4. Text for new tickets

- "Raise or remove the engine's own nesting limit (`MAX_DEPTH = 100`) so that 999 unary operators and more than 100 nested parentheses work as in SQLite." This needs a parser that does not use the call stack per level, or a measured safe stack depth.

## 5. Findings

1. **Comparison height rule for some operators is unverified: fixed.** The finding was real as a rule 5 gap: `<=`, `>`, `>=`, `==`, `!=` had no measurement. I measured them on both sides; the engine already matched, so no logic changed. The comment records the results and the new test covers them. The new test passes without any code change, because there was no behaviour bug to fail on. It is a regression guard, not a failing-first test.
2. **Parenthesis-depth claim about SQLite has no record: fixed.** The claim was true but unsourced. I measured 1001 and 1002 pairs on SQLite 3.46.1 (both give 7) and the comment now cites that. No test accompanies it: the engine rejects this input at its own limit of 100, which `deep_nesting_is_an_error` already covers.
