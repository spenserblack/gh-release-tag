# gh-release-tag

Generate release notes in your tag message

## How it works

This calls the GitHub API to generate release notes. However, instead
of creating a release, this creates a *tag* using the generated release
notes. This uses `git tag --cleanup=verbatim` to preserve all text
in the release notes, instead of treating some of it as commentary to
be trimmed.
