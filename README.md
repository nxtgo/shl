# shl

experimental syntax highlight parser in the terminal.

# about

basically this program takes a .json file used to define syntax
highlight rules, as seen in [examples](examples).
this project's usage is intended for buffers.

# usage

```sh
go run main.go rules.json file.ext
```

or use pipes

```sh
cat file.ext | go run main.go rules.json
```
# guide

rules are defined in a json array where each rule has:

* `pattern`: regex pattern to match
* `color`: hex color like `#ff0000` or named color like `red`
* `capture`: (optional) specify which capture group to color (1-indexed)
* `nested`: (optional) array of patterns to apply inside this match

rules are applied in order. once text is matched by a rule, it cannot be matched again by later rules (non-overlapping).

basic example:
```json
[
  {
    "pattern": "\\bif\\b",
    "color": "#bb9af7"
  }
]
```

capturing groups example (colors only the function name, not the parenthesis):
```json
{
  "pattern": "\\b([a-z][a-zA-Z0-9_]*)\\s*\\(",
  "color": "#7aa2f7",
  "capture": 1
}
```

nested patterns example (colors format specifiers inside strings):
```json
{
  "pattern": "\"(?:\\\\.|[^\\\\\"])*\"",
  "color": "#9ece6a",
  "nested": ["%[#0\\-+ ]?\\d*\\.?\\d*[vTtbcdoOqxXUeEfFgGsp%]"]
}
```

# examples

* [go](examples/go) - tokyonight colors for go syntax
* [c](examples/c) - tokyonight colors for c syntax

# license

available under cc0 1.0 (public domain).
