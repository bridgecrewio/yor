# Security

## Reporting a Vulnerability

If you think you have found a potential security vulnerability in `yor`,
please email psirt@paloaltonetworks.com directly. Do not file a public issue. If
English is not your first language, please try to describe the problem
and its impact to the best of your ability. For greater detail, please
use your native language and we will try our best to translate it using
online services.

Please also include the code you used to find the problem and the
shortest amount of code necessary to reproduce it.

Please do not disclose this to anyone else. We will retrieve a CVE
identifier if necessary and give you full credit under whatever name or
alias you provide. We will only request an identifier when we have a fix
and can publish it in a release.

We will respect your privacy and will only publicize your involvement if
you grant us permission.

## Process

The following information discusses the process the `yor` project
follows in response to vulnerability disclosures. If you are disclosing
a vulnerability, this section of the documentation lets you know how we
will respond to your disclosure.

### Timeline

When you report an issue, one of the project members will respond to you
within a few days. This initial response will at the very least confirm
receipt of the report.

If we were able to rapidly reproduce the issue, the initial response
will also contain confirmation of the issue. If we are not, we will
often ask for more information about the reproduction scenario.

Our goal is to have a fix for any vulnerability released within two
weeks of the initial disclosure. This may potentially involve shipping
an interim release that simply disables function while a more mature fix
can be prepared, but will in the vast majority of cases mean shipping a
complete release as soon as possible.

Throughout the fix process, we will keep you up to speed with how the fix
is progressing. Once the fix is prepared, we will notify you that we
believe we have a fix. Often we will ask you to confirm the fix resolves
the problem in your environment, especially if we are not confident of
our reproduction scenario.

At this point, we will prepare for the release. We will obtain a CVE
number if one is required, providing you with full credit for the
discovery. We will also decide on a planned release date, and let you
know when it is.

On the release day, we will push the patch to our public repository, along
with an updated changelog that describes the issue.

At this point, we will publicize the release. This will involve
announcement on our Slack channel (https://codifiedsecurity.slack.com/)
and all other communication mechanisms available to the core team.

We will also explicitly mention which commits contain the fix to make it
easier for other distributors and users to easily patch their own
versions of `yor` if upgrading is not an option.

## Known Unfixed Vulnerabilities

The following CVEs are present in transitive or direct dependencies and have
**not** been patched in the current release because every available upstream
fix raises the minimum required Go language version above the `go 1.19`
directive declared in `go.mod`. Bumping the language version is considered a
breaking change for downstream consumers and CI pipelines, so the upgrades
have been intentionally deferred.

| CVE | Module | Current | Upstream fix | Reason deferred |
| --- | --- | --- | --- | --- |
| CVE-2025-21614 | `github.com/go-git/go-git/v5` | v5.11.0 | v5.13.0 | Requires `go >= 1.23` |
| CVE-2024-6257  | `github.com/hashicorp/go-getter` | v1.6.2 | v1.7.5 | Requires `go >= 1.23` (also pulls in ~40 new transitive deps) |
| CVE-2025-8959  | `github.com/hashicorp/go-getter` | v1.6.2 | v1.7.9 | Same as above |
| CVE-2026-4660  | `github.com/hashicorp/go-getter` | v1.6.2 | v1.8.6 | Same as above |
| CVE-2025-0377  | `github.com/hashicorp/go-slug`   | v0.5.0 | v0.16.3 | Requires `go >= 1.22` |

### Stdlib CVEs (mitigated)

The following Go stdlib CVEs are mitigated by the `toolchain go1.26.3`
directive in `go.mod`, which pins the build toolchain to a patched release
without changing the `go 1.19` language compatibility level:

- CVE-2026-39836 (`net`)
- CVE-2026-33814 (`net/http`)
- CVE-2026-33811 (`net`)

Builders using Go `< 1.26.3` should upgrade their toolchain (or use Go
`>= 1.25.10` on the 1.25.x line) to ensure the stdlib fixes are applied.
