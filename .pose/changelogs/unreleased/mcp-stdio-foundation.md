---
spec: mcp-stdio-foundation
category: added
breaking: false
refs:
---

`codinho serve` now starts a local MCP server over stdio, exposing an initial
vertical slice for agents to browse the curriculum and run a minimal
learning session: `catalog_search`, `catalog_get`, `session_start`,
`session_get` and `instruction_get`.
