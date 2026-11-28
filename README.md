# Evidence for the Qodo Software Factory Hackathon

This branch holds the evidence that the daily posts link to.
It has no code.

## How each file was made

- `*.png`: a screenshot of the public GitHub page, taken with a headless browser that was not logged in, then cropped. Nothing in the image was edited.
- `day01/qodo-first-review-pr2.md`: the body of Qodo's first review comment on pull request 2, as the GitHub API returned it on 2026-11-21 before the re-review replaced it. The same text is in the edit history of that comment.
- `day01/repro-high-findings.sh`: the script that reproduced the three High findings on commit `a364ba0`.
- `day01/demo-output.txt`: the output of `yantra referee diff` on a temp repo with one weakened test, run on commit `fd63959`.
- `day02/qodo-standards-after-import.json`: the output of `qodo rules list` on 2026-11-27 after the import from `AGENTS.md`, with internal IDs removed.
- `day03/ticket5-acceptance-run.txt`, `day03/ticket6-acceptance-run.txt`: the output of the new commands against the real issues and commits of this repo.
- `day03/qodo-first-review-pr7.md`, `day03/qodo-first-review-pr8.md`: the body of the first Qodo review comment on each pull request, as the GitHub API returned it before the re-review replaced it.
- `day01-database/`: evidence for the Day 1 post after it moved to the database repo `yantra-sql`. The two Qodo files are comment bodies from the GitHub API. `edge-probe.rs` is a throwaway test that was run against the merged engine, and `edge-probe-output.txt` is its output.
- `day02-database/`: evidence for the Day 2 post on the database repo. `guard-bypass-trial.sh` tries ways past the checks, and `guard-bypass-trial.txt` is its output for the first CI steps and for the guard script that replaced them.
- `day03-database/`: evidence for the Day 3 post on the database repo. The read-back files are the output of `gh issue view`. `unseen-files-run.txt` is the result of a run on two test files that are not in the repo.
- `day04/`: evidence for the Day 4 post. `tick.log` and `run-ticket20.log` are the logs of the scheduled job and of its run on ticket 20, as the runner wrote them. `agent-report-ticket20.md` is the final message of the headless agent. `probe-queries.txt` and `probe-lekha-vs-sqlite.tsv` are a check of the result against SQLite 3 that was made after the run.
- `day05/`: evidence for the Day 5 post. The two review files are comment bodies from the GitHub API. The JSON file is the output of `qodo rules list` after the second import, with internal IDs removed.
