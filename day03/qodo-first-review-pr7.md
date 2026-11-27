<h3>Code Review by Qodo</h3>

<br/>

<img src="https://img.shields.io/badge/High-634FD1?style=flat-square" height="20px" alt="Action required">

<details>
<summary>  1.  Unreadable issue lists pass as ready <code>🐞 Bug</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
><b><i>runTicketList</i></b> checks JSON syntax but does not verify that the result is an issue array containing
>valid issue objects. If <b><i>gh</i></b> returns <b><i>null</i></b>, the command prints <b><i>[]</i></b> and exits 0; a partial object
>with a ready body can likewise exit 0 while reporting a fabricated issue number and empty metadata.
></pre>
></details>

> <details>
><summary>Agent prompt</summary>
><br/>
>
>```
>The issue below was found during a code review. Follow the provided context and guidance below and implement a solution
>
>## Issue description
>Valid JSON with the wrong shape or missing issue fields can make `ticket list` report success for input it could not read.
>## Fix Focus Areas
>- cmd/yantra/main.go[199-211]
>- cmd/yantra/main_test.go[113-126]
>## Recommended Fix
>Reject a null top-level value, null entries, nonpositive issue numbers, and missing required issue fields with exit code 2 before constructing tickets. Add tests for `null` and an incomplete issue object.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/4ac7d56c-efc1-463a-b65a-9240ea1a3a47?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/4ac7d56c-efc1-463a-b65a-9240ea1a3a47)

</details>


<details>
<summary>  2.  Empty scope passes the readiness gate <code>🐞 Bug</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
><b><i>scopePaths</i></b> appends a backticked match even when trimming leaves an empty string. An issue
>containing an acceptance criterion and <b><i>Scope: </i></b> followed by backticks around a space receives
><b><i>scope: [&quot;&quot;]</i></b> and is marked ready despite naming no path.
></pre>
></details>

> <details>
><summary>Agent prompt</summary>
><br/>
>
>```
>The issue below was found during a code review. Follow the provided context and guidance below and implement a solution
>
>## Issue description
>A whitespace-only backticked scope creates an empty path entry that satisfies the readiness check.
>## Fix Focus Areas
>- internal/ticket/ticket.go[63-76]
>- internal/ticket/ticket_test.go[8-53]
>## Recommended Fix
>Trim each extracted backticked path and append it only when nonempty, as the unquoted path branch already does. Test that a whitespace-only quoted scope produces `no scope paths` and is not ready.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/740a7530-887a-4b4c-bf8f-b8a824f133b9?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/740a7530-887a-4b4c-bf8f-b8a824f133b9)

</details>


<br/>

<details><summary><ins><strong>View medium (3)</strong></ins></summary><br/>
<img src="https://img.shields.io/badge/Medium-634FD1?style=flat-square" height="20px" alt="Remediation recommended">

<details>
<summary>  3.  More than 100 tickets disappear from lists <code>🐞 Bug</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
><b><i>runTicketList</i></b> requests at most 100 issues from <b><i>gh</i></b> and never fetches additional results. In a
>repository with more than 100 open <b><i>ticket</i></b> issues, the remaining issues are absent from the JSON
>and cannot affect the readiness exit code.
></pre>
></details>

> <details>
><summary>Agent prompt</summary>
><br/>
>
>```
>The issue below was found during a code review. Follow the provided context and guidance below and implement a solution
>
>## Issue description
>The fixed `gh issue list --limit 100` result silently omits matching tickets beyond the first 100.
>## Fix Focus Areas
>- cmd/yantra/main.go[191-217]
>- cmd/yantra/main_test.go[113-126]
>## Recommended Fix
>Retrieve all matching open issues rather than treating a capped response as complete. Ensure the output and aggregate readiness include results beyond the first batch, and add a test covering more than 100 issues.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/45479008-1c87-49c5-908b-f42a98186eb6?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/45479008-1c87-49c5-908b-f42a98186eb6)

</details>


<details>
<summary>  4.  Code samples can make tickets look ready <code>🐞 Bug</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
><b><i>ticket.New</i></b> applies the <b><i>Scope:</i></b> matcher to every body line without excluding fenced code blocks. A
>ticket with a real acceptance criterion but no scope declaration can be marked ready when a code
>sample contains a line such as <b><i>Scope: /etc/passwd</i></b>.
></pre>
></details>

> <details>
><summary>Agent prompt</summary>
><br/>
>
>```
>The issue below was found during a code review. Follow the provided context and guidance below and implement a solution
>
>## Issue description
>A `Scope:` example inside fenced Markdown is treated as ticket metadata and can satisfy readiness.
>## Fix Focus Areas
>- internal/ticket/ticket.go[29-49]
>- internal/ticket/ticket_test.go[8-53]
>## Recommended Fix
>Track fenced-code regions while scanning the issue body and skip metadata extraction within them, including acceptance headings and list items. Add a test with a criterion outside the fence and the only `Scope:` line inside it.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/f332aa4c-f0e4-40a0-945a-0bbf3dc6f6f1?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/f332aa4c-f0e4-40a0-945a-0bbf3dc6f6f1)

</details>


<details>
<summary>  5.  Valid headings make tickets look unready <code>🐞 Bug</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
><b><i>criteriaHead</i></b> rejects an ATX heading with closing hashes, while the list scanner treats a Setext
>underline as the end of the criteria section. An issue using <b><i>## Acceptance Criteria ##</i></b> or an
>underlined <b><i>Acceptance Criteria</i></b> heading can contain a list and a scope path yet be reported as
>having no acceptance criteria.
></pre>
></details>

> <details>
><summary>Agent prompt</summary>
><br/>
>
>```
>The issue below was found during a code review. Follow the provided context and guidance below and implement a solution
>
>## Issue description
>Valid closing-hash and Setext acceptance headings prevent their following list items from being collected.
>## Fix Focus Areas
>- internal/ticket/ticket.go[20-46]
>- internal/ticket/ticket_test.go[8-53]
>## Recommended Fix
>Accept optional closing hashes on an ATX acceptance heading, and recognize or skip the underline immediately following a Setext acceptance heading before collecting list items. Add tests for both heading forms with criteria and a scope path.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/dff80938-acb0-47a3-89c7-3f289cb0e74f?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/dff80938-acb0-47a3-89c7-3f289cb0e74f)

</details>


</details>

<br/>

<!-- qodo-self-review-slot -->
<!-- https://github.com/bettersandeep/yantra/commit/69047928281db8cd60f305f12a0908e7e5781ddf -->
