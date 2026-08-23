---
spec: session-orchestration-disclosure
category: added
breaking: false
refs:
---

`ailearn serve` now supports the full session lifecycle: `session_configure`,
`session_pause`, `session_resume`, `session_finish` and `granularity_adjust`
join `session_start`/`session_get`/`instruction_get`, letting an agent
pause, resume, finish and adjust the instructional depth of a running
session with optimistic-concurrency safety and idempotent retries.
