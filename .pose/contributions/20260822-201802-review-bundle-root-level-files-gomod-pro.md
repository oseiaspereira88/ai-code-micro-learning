---
title: "review bundle: root-level files (go.mod, PROJECT.md) are unclassifiable blockers, not excluded"
type: "bug"
module: ""
created_at: "2026-08-22T20:18:02Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: review bundle: root-level files (go.mod, PROJECT.md) are unclassifiable blockers, not excluded

POSE 1.7.6 schema-v1 reproduction: in a single-module Go repo where the module root is '.' (go.mod at repo root, source under cmd/ and internal/ subdirectories — the normal layout for a small Go CLI), attribute a change set to a spec that modifies go.mod plus files under cmd/ and internal/. Run 'pose review bundle spec:<slug> --explain' or '--seal'. Observed: cmd/*.go and internal/**/*.go classify correctly as class:implementation (matched against an auto-discovered component whose Path is the subdirectory), but go.mod at the repo root produces BOTH 'exclude.=go.mod reason:attributed path is outside the semantic review subject' AND '[ERROR] unclassified review subject path go.mod', and the bundle refuses to seal. The same happens for any other root-level project file that isn't one of a small hardcoded set (only POSE.md/AGENTS.md/README.md are recognized as documentation; PROJECT.md, go.mod, go.sum, package.json, Cargo.toml, tsconfig.json, .gitignore, CI config at repo root are not covered by any rule) and isn't a governance path under .pose/. Root cause (inferred from behavior, since the classification list is clearly tuned for this project's own multi-directory layout: pose-mcp/, docs-site/, mcp-enforce/, scripts/, tests/, locales/, .github/): the per-component matching loop skips a component whose Path resolves to '.' (the repo root) specifically to avoid it swallowing everything recursively -- but that means a project whose only module IS the root can never get any of its bare root-level files classified, and an unclassified path is treated as BOTH excluded from the subject AND a hard blocker (mutually exclusive outcomes are set on the same file), so the bundle can never seal while such a file is part of the attributed change set. This makes review bundles unreachable for the common case of a small single-module project (Go, Node, Rust, ...) whose manifest lives at the repo root rather than in a named subpackage directory. Expected: either (a) an unclassified attributed path should be excluded from the subject (as the exclude.= line already claims) without also being appended to blockers, and/or (b) common root-level project manifests (go.mod, go.sum, package.json, package-lock.json, Cargo.toml, Cargo.lock, pyproject.toml, tsconfig.json) should classify as 'governance' or a new 'manifest' class out of the box, and/or (c) a component whose discovered/declared Path is '.' should still classify its own direct (non-recursive) root-level files instead of being skipped outright. No project code or private paths are required to reproduce: a fresh 'pose init' plus 'go mod init' at the repo root, one file under a subdirectory, and one spec whose attributed change set includes both is sufficient.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/38
