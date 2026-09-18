# About hook-store

## What it owns

SQLite registry of bridge-managed harness hooks — each row a handler (harness, event, matcher, shell command) bound to a `global`, `instance` or `session` scope. The server queries it at session spawn to decide which hooks to wire into the harness's native mechanism. A pure data layer: it does not execute hooks, render settings JSON, or talk to harnesses.

## Where this prompt lives

These sections are stored in agent-store as a project prompt collection and rendered, with identical text, to `AGENTS.md` and `CLAUDE.md` at the root of this repo, so that whichever file a harness reads it gets the same thing. Edit them on dash `/files`, or edit either rendered file: the 15-minute scan carries the edit back into the sections and out to the other file. The host prompt keeps one row for this repo with only what an agent elsewhere needs.
