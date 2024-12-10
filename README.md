# sedr

`sedr` is a command-line tool for in-place file editing using regular expressions.

## Usage

```
sedr <regex> <original> <substitution> <file>
```

- `<regex>`: Regular expression pattern
- `<original>`: Original text to replace (default: $0)
- `<substitution>`: Replacement text
- `<file>`: File to process

## Example

```bash
sedr '(\d{4})-(\d{2})-(\d{2})' '$0' '$1-$3-$2' file.txt
```

This command will change date formats from YYYY-MM-DD to YYYY-DD-MM in file.txt.

## Features

- In-place file editing
- Regular expression support
- Capture group substitution
- Preserves file permissions

## Notes

- `$0` represents the entire matched string
- `$+` can be used to represent a literal `$`
- Use `$1`, `$2`, etc. in the substitution to refer to capture groups

## Installation

```
go install github.com/doraemonkeys/sedr@latest
```
