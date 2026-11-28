SQLite enforces three limits that the engine does not have. Above each limit the engine returns a value where SQLite gives an error. This breaks rule 1 of AGENTS.md. The coding agent found these during ticket #10.

This ticket was revised after an investigation with Qodo Codebase Wisdom. The first wording is in the edit history of this issue.

## Sources

From https://www.sqlite.org/limits.html (the words are SQLite's):

> The default setting for SQLITE_MAX_COLUMN is 2000. [It bounds, among others,] the number of columns in the result set of a SELECT statement.

> The SQLITE_MAX_EXPR_DEPTH parameter determines the maximum expression tree depth. [...] The current implementation has a default value of 1000.

> The SQLITE_MAX_FUNCTION_ARG parameter determines the maximum number of parameters that can be passed to an SQL function. For many years, the default value was about 100 but the default was raised to 1000 beginning with SQLite version 3.48.0 (2025-01-14).

## What the investigation found

- The engine has its own limit, `MAX_DEPTH = 100` in `src/select.rs`. It counts each recursion of the parser through `unary()`. This is not the depth that SQLite counts.
- The engine has no syntax tree. Each parse function returns a value (`Result<Option<i64>, Error>`). The binary operators are loops, so a chain of 1001 terms does not go deeper in the parser, and the engine cannot see it today.
- SQLite counts the height of the expression tree. Measured on SQLite 3.46.1 and 3.49.1: a sum of 1000 terms works and 1001 is an error; 999 unary minus signs work and 1000 are an error; 1001 pairs of parentheses work.
- The argument limit depends on the version. SQLite 3.46.1, which `yt sqlite` runs, stops at 127 arguments. SQLite 3.48 and later stop at 1000. This ticket takes the documented default of today, 1000.
- Pull request #12 once added a limit of 1000 tokens for a whole statement and removed it again, because it was not in its ticket and it rejected statements that SQLite accepts. Do not bring that back.
- All six pull requests that built the engine changed `src/select.rs`. No other repo depends on this one.

Acceptance criteria:

- A SELECT with 2000 result columns works. With 2001 it is an error.
- A function call with 1000 arguments works. With 1001 it is an error. Do not use `yt sqlite` to check this criterion, because its SQLite version has a different limit. The expected results are the two in this line.
- The engine computes the height of each expression as SQLite does: an integer literal or NULL has height 1, a unary operator and a binary operator have the height of their highest operand plus 1, a function call has the height of its highest argument plus 1, and parentheses add nothing. An expression with a height above 1000 is an error.
- These cases hold: a sum of 1000 terms works and 1001 is an error. 999 unary minus signs before a literal need the engine's own nesting limit to allow them; if that limit stops them first, the error stays and the report says so. `COALESCE(NULL, <sum of 999 terms>)` works and the same with 1000 terms is an error.
- For `CASE`, `CAST`, comparisons and the boolean operators, find the height rule with `yt sqlite` (the depth limit is the same in its version), write the queries and results in the report, and name the rule in a code comment.
- Each limit names its source in a code comment.
- `tests/slt/baseline.txt` has the exact pass count. The number does not go down.
- No record that passed before fails now.

Scope: `src/`, `tests/slt/baseline.txt`, `README.md`.
