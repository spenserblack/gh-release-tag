# gh-release-tag

[![CI](https://github.com/spenserblack/gh-release-tag/actions/workflows/ci.yml/badge.svg)](https://github.com/spenserblack/gh-release-tag/actions/workflows/ci.yml)

Generate release notes in your tag message

## How it works

This calls the GitHub API to generate release notes. However, instead
of creating a release, this creates a *tag* using the generated release
notes. This uses `git tag --cleanup=verbatim` to preserve all text
in the release notes, instead of treating some of it as commentary to
be trimmed.
