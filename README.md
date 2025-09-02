# note

Quick CLI note-taking tool.

## Install

### Go Install
```bash
go install github.com/ali-chapman/note@latest
```

### Binary Download
Download from [releases](https://github.com/ali-chapman/note/releases) and add to PATH.

## Usage

```bash
note                    # List/search existing notes
note "meeting notes"    # Create or edit a note
```

## Requirements

- `fzf` for note selection
- `$EDITOR` environment variable set
- Optional: `$NOTES_DIRECTORY` (defaults to `~/.notes`)

## Notes

Notes are stored as markdown files with date prefixes: `02-Jan-2006 note-name.md`