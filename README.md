# hcledit
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/minamijoyo/hcledit.svg)](https://github.com/minamijoyo/hcledit/releases/latest)
[![GoDoc](https://godoc.org/github.com/minamijoyo/hcledit/hcledit?status.svg)](https://godoc.org/github.com/minamijoyo/hcledit)

## Features

- Edit HCL with command line: attribute get/set, block get/list/mv
- CLI-friendly: easily pipe and combine other commands
- Schemaless: independent of specific HCL applications (such as Terraform)
- Keep comments: very important when editing in automation scripts
- HCL2 support (not HCL1)

The hcledit focuses on editing HCL with command line, is not aiming for generic query tools. It was originally born for the refactoring of Terraform configurations, but is not limited to specific applications.
The HCL specification is somewhat generic, so usability takes precedence over strictness if there is room for interpreting meanings in a schemaless approach.

## Install

### Homebrew

If you are macOS user:

```
$ brew install minamijoyo/hcledit/hcledit
```

### Download

Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/minamijoyo/hcledit/releases

### Source

If you have Go 1.13+ development environment:

```
$ git clone https://github.com/minamijoyo/hcledit
$ cd hcledit/
$ make install
$ hcledit --version
```

## Usage

```
$ hcledit --help
A stream editor for HCL

Usage:
  hcledit [command]

Available Commands:
  attribute   Edit attribute
  block       Edit block
  help        Help about any command
  version     Print version

Flags:
  -h, --help   help for hcledit

Use "hcledit [command] --help" for more information about a command.
```

### attribute

```
$ hcledit attribute --help
Edit attribute

Usage:
  hcledit attribute [flags]
  hcledit attribute [command]

Available Commands:
  get         Get attribute
  set         Set attribute

Flags:
  -h, --help   help for attribute

Use "hcledit attribute [command] --help" for more information about a command.
```

### block

```
$ hcledit block --help
Edit block

Usage:
  hcledit block [flags]
  hcledit block [command]

Available Commands:
  get         Get block
  list        List block
  mv          Move block (Rename block type and labels)

Flags:
  -h, --help   help for block

Use "hcledit block [command] --help" for more information about a command.
```

## License

MIT
