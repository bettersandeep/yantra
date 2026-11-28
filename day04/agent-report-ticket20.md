# Report: ticket 20, comparison and boolean operators

## 1. What changed

- `src/lexer.rs`: six new tokens for `=`/`==`, `<>`/`!=`, `<`, `<=`, `>`, `>=`. `<<` and `>>` are explicit lexer errors, and `!` alone, `|`, `&`, `~` stay errors.
- `src/select.rs`: four new parser levels above the arithmetic (`OR`, `AND`, the `=`/`IS` row, the `<` row) and a `compare` function.
- `NOT` is read in `unary`, with an operand at the `=` row. SQLite accepts `NOT` as the operand of any operator and lets it take the rest of the expression up to `AND`/`OR`; this reproduces that.
- `IS NULL` and `IS NOT NULL` accept only the literal `NULL` on the right. `1 IS NULL + 1` is an explicit error, because SQLite reads it as `1 IS (NULL + 1)`.
- Tests: 7 new tests in `select.rs` and 1 in `lexer.rs`. The old "must be an error" entries for `1 = 1`, `1 < 2`, `1 AND 1`, `NOT 1` and `1 IS NULL` are removed.
- `README.md`: the Status section lists the new operators and their rules.
- `tests/slt/baseline.txt`: not changed. The measured count is still 8995.

Sources (rule 5): the precedence table in section 2 of https://www.sqlite.org/lang_expr.html, and `yt sqlite` (SQLite 3.46.1) for each expected value. The code comments quote the statements and answers.

## 2. Acceptance criteria

| Criterion | Check that I ran | Result |
|---|---|---|
| Comparisons give 1 or 0, NULL operand gives NULL | `yt sqlite`: `SELECT 1 = 1, 1 == 2, 1 <> 2, 1 != 1, 1 < 2, 2 <= 2, 1 > 2, 2 >= 3` and the same with NULL operands. Unit test `comparison_gives_one_zero_or_null` | SQLite: `1,0,1,0,1,1,0,0` and eight NULLs. Test passes with the same values |
| `NOT`, `AND`, `OR` three-valued logic | `yt sqlite`: all 10 operand pairs of `AND` and of `OR`, and `NOT 0, NOT 1, NOT 5, NOT - 3, NOT NULL, NOT NOT 7`. Unit test `not_and_or_follow_three_valued_logic` | SQLite: `NULL AND 0` = 0, `NULL OR 1` = 1, `NOT NULL` = NULL, `2 AND - 3` = 1. Test passes with the same values |
| `IS NULL` / `IS NOT NULL` give 1 or 0 | `yt sqlite`: `SELECT 1 IS NULL, NULL IS NULL, 1 IS NOT NULL, NULL IS NOT NULL, 1 / 0 IS NULL, 0 IS NULL`. Unit test `is_null_and_is_not_null_give_one_or_zero` | SQLite: `0,1,1,0,1,0`. Test passes |
| Precedence of the table, left to right in one row | `yt sqlite`: 50 expressions across 7 statements (for example `1 = 2 < 3`, `3 > 2 > 1`, `NOT 1 = 2`, `1 OR 0 AND 0`, `NULL IS NULL = 1`, `- NOT 0 = 5`). Unit tests `precedence_of_comparison_and_boolean_operators` and `not_as_an_operand_takes_the_rest_of_the_expression` | Tests pass with SQLite's values |
| Result is an integer, usable in arithmetic, functions, aggregates | `yt sqlite`: `SELECT (1 < 2) + (2 < 3), (1 = 1) * 5, - (1 = 1), COALESCE(NULL = 1, 7), NULLIF(1 = 1, 1), COUNT(NULL = 1), SUM(2 > 1), MAX(1 AND NULL), MIN(NOT 0), CAST(1 < 2 AS INTEGER)`. Unit test `comparison_result_is_an_integer_operand` | SQLite: `2,5,-1,7,NULL,0,1,NULL,1,1`. Test passes |
| Each other operator is an error | Unit test `other_operators_are_errors` (62 statements) and the lexer test `rejects_unsupported_text` | Tests pass: each one is an error |
| Baseline has the exact pass count, not lower | `yt check` | `slt: random/expr/slt_good_0.test 8995/10012`, same as before, baseline test passes |
| No record that passed before fails now | `yt check` | The count is equal, 8995. I did not compare the set of passing records one by one. The change only accepts input that was an error before |

`yt check` was green (fmt, clippy, 37 unit tests, 10 harness tests) before the README edit. After the README edit I ran `yt failures 3`, which compiled and gave 8995 again. I did not run the full `yt check` again after the README edit.

## 3. Known differences from SQLite

- **The score did not go up.** The two records I saw that use these operators (lines 92 and 2444) have them inside `CASE`, so they still fail until the `CASE` ticket.
- `AND` and `OR` compute both operands. `SELECT 0 AND 9223372036854775807 + 1` gives 0 in SQLite 3.46.1 and an integer overflow error here.
- Values with SQLite, errors here: `IS`/`IS NOT` with a right side that is not the literal `NULL`, including `1 IS (NULL)`, `1 IS NULL + 1` (SQLite: 0), `1 IS NOT NULL - 1` (SQLite: 1) and `1 IS NOT NOT NULL` (SQLite: 1).
- Values with SQLite, errors here, as the ticket asks: `~`, `||`, `&`, `|`, `<<`, `>>`, `BETWEEN`, `IN`, `LIKE`, `GLOB`, `ISNULL`, `NOTNULL`, postfix `NOT NULL`, `IS [NOT] DISTINCT FROM`.

## 4. Text for new tickets

- `CASE` expressions, both forms. This is the ticket that turns the new operators into passing records.
- `AND`/`OR` should not fail on an error in an operand that SQLite does not need. Needs a source for SQLite's exact evaluation rule first; `yt sqlite` line 43 shows `abs(-9223372036854775807 - 1) AND 0` gives 0 there.
- `IS` and `IS NOT` with a general right side, and `IS [NOT] DISTINCT FROM`.
- `ISNULL`, `NOTNULL` and postfix `NOT NULL`.
- `CAST(<expr> AS REAL)` and `AVG`: both appear in the first 40 failing records.
- `CREATE TABLE` and `INSERT`: the first 12 failing records.

## 5. Review findings

Not a review round.
