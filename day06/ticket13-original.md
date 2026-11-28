SQLite 3.49.1 enforces three limits that the engine does not have. For each one the engine returns a value where SQLite gives an error. This breaks rule 1 of AGENTS.md. The coding agent found these during ticket #10.

| Limit | Query | SQLite 3.49.1 | lekha |
|---|---|---|---|
| Argument count, max 1000 | `COALESCE` with 1001 arguments | error: too many arguments on function COALESCE | a value |
| Expression depth, max 1000 | `SELECT 1 + 1 + ... + 1` with 1001 terms | error: Expression tree is too large (maximum depth 1000) | value 1001 |
| Column count, max 2000 | `SELECT 1, 1, ... , 1` with 2001 items | error: too many columns in result set | a row of 2001 values |

Notes:
- SQLite counts the height of the expression tree. A chain of binary operators adds one level per term, a function call adds one level above its deepest argument, and parentheses add none.
- The three numbers are observed, not read from a document. Rule 5 needs https://www.sqlite.org/limits.html as the source before the numbers go into code.
- The engine has its own limit of 100 nested levels. That one gives an error where SQLite gives a value, which rule 1 allows.

Acceptance criteria:
- For each of the three limits, the engine gives an error for the first statement over the limit and a value for the last statement under it, as SQLite does.
- Each limit names its source in a code comment.
- The score does not go down.

Scope: `src/`, `README.md`.
