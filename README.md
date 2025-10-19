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

# examples

* [go.json](examples/go.json) - tokyonight colors for go syntax

# license

available under CC0 1.0 (public domain).
