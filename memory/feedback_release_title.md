---
name: GitHub release title format
description: Use plain version number (e.g. "0.1.0") as the release title, not prefixed with "v"
type: feedback
---

Use the bare version number as the GitHub release title (e.g. "0.1.0"), not "v0.1.0".

**Why:** User corrected this explicitly when creating the first release.

**How to apply:** When running `gh release create`, set `--title` to the version without the `v` prefix. The tag itself should still use the `v` prefix (e.g. `v0.1.0`).
