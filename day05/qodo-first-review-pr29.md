<h3>Code Review by Qodo</h3>

<br/>

<img src="https://img.shields.io/badge/High-634FD1?style=flat-square" height="20px" alt="Action required">

<details>
<summary>  1.  A pull request can silence the guard <code>🐞 Bug</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
>The Guard step redirects the base-branch script into <b><i>ci/guard.sh</i></b>, but shell redirection follows a
>symlink supplied by the pull request. If the pull request replaces that file with a symlink to
><b><i>/dev/null</i></b>, the fetched script is discarded and <b><i>sh ci/guard.sh</i></b> exits successfully without running
>the checks that reject combined engine and measurement changes.
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
>A pull request can make `ci/guard.sh` a symlink to `/dev/null`, causing the base-branch guard to be discarded while the Guard step succeeds.
>## Fix Focus Areas
>- .github/workflows/ci.yml[19-23]
>## Recommended Fix
>Reject a symlinked `ci` directory, fetch the base guard into a temporary file outside the checkout, then remove the PR-controlled `ci/guard.sh` path before installing the fetched contents as a regular file there. Run that installed file so the guard's repository-relative paths continue to work. Add a CI regression test using a pull request that replaces the guard with a symlink and changes the engine.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/1b8e552f-1784-438f-b734-94dcea4f9cf9?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/1b8e552f-1784-438f-b734-94dcea4f9cf9)

</details>


<br/>

<details><summary><ins><strong>View medium (1)</strong></ins></summary><br/>
<img src="https://img.shields.io/badge/Medium-634FD1?style=flat-square" height="20px" alt="Remediation recommended">

<details>
<summary>  2.  A pull request can skip the guard in CI <code>📎 Requirement gap</code></summary>

<br/>

> <details>
><summary>Description</summary>
><br/>
>
><pre>
>The new <b><i>Guard</i></b> step in <b><i>.github/workflows/ci.yml</i></b> fetches the base branch&#x27;s <b><i>ci/guard.sh</i></b> only when
>the workflow runs, but the workflow file itself is read from the pull request&#x27;s head. A pull request
>that deletes this step and also changes engine code reaches CI with no guard and still passes, which
>is the failure case of rule 5040763; the PR description discloses this gap, but the failure criteria
>still hold until the workflow file is protected.
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
>The guard step in `.github/workflows/ci.yml` is part of the file that a pull request can change. A pull request can delete the step and change the engine, and CI still passes, which violates rule 5040763.
>
>## Fix Focus Areas
>- .github/workflows/ci.yml[16-22]
>
>## Recommended Fix
>Run the guard from a workflow that the pull request cannot modify, for example a `pull_request_target` job that checks out the base branch and runs its `ci/guard.sh`, or add a required status check whose workflow lives on the base branch. Alternatively, add a CODEOWNERS entry requiring review for `.github/` and document the residual risk in AGENTS.md rule 6.
>```
> <code>ⓘ Copy this prompt and use it to remediate the issue with your preferred AI generation tools</code>
></details>

[Dismiss ↗](https://app.qodo.ai/findings/5616ccdc-e8a1-4acb-9132-1d4857574403?action=dismiss&origin=git) | [View ↗](https://app.qodo.ai/findings/5616ccdc-e8a1-4acb-9132-1d4857574403)

</details>


</details>

<br/>

<!-- qodo-self-review-slot -->
<!-- https://github.com/bettersandeep/lekha/commit/dd011b938c2487698f0ae0ab9ef07e17a56df366 -->
