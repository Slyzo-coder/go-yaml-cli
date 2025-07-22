go-yaml-cli
===========

A command line tool for the `github.com/yaml/go-yaml` library


## Description

This repo creates a command line tool called `go-yaml`.

It read YAML from stdin and outputs it in many formats, including the internal
states of `go-yaml`.

This program is extremely useful for understanding / debugging how
[go-yaml](https://github.com/yaml/go-yaml) loads YAML documents.

You might want to use its output to report bugs for behaviors of go-yaml or its
downstream consumers.


## Installation

Run:

```bash
$ make install PREFIX=<prefix>
```

Where `<prefix>` is the directory that contains a `bin` subdirectory where the
`go-yaml` binary will be installed: `<prefix>/bin/go-yaml`.

If you have GOROOT set in your environment, you could use:

```bash
$ make install PREFIX="$GOROOT"
```

> Note: `go install github.com/yaml/go-yaml-cli@latest` doesn't work for
> this repo.


## Usage

This program reads YAML from stdin and prints output according to the options
you give it (to stdout).

Example commands:

```
$ go-yaml --help

$ <file.yaml go-yaml -y
$ <file.yaml go-yaml -Y
$ <file.yaml go-yaml -J
$ <file.yaml go-yaml -t -c
$ <file.yaml go-yaml -e -p -c
$ <file.yaml go-yaml -n
```


## Testing

```
make test
```

No dependencies required.
On first run, this will install all dependencies including a local installation
of Go.


## License

Copyright 2025 - Ingy döt Net

This project is licensed under the Apache License, Version 2.0.

See the [LICENSE](LICENSE) file for details.
