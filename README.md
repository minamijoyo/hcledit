# hcledit
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/minamijoyo/hcledit.svg)](https://github.com/minamijoyo/hcledit/releases/latest)
[![GoDoc](https://godoc.org/github.com/minamijoyo/hcledit/hcledit?status.svg)](https://godoc.org/github.com/minamijoyo/hcledit)

## Features

- CLI-friendly: Read HCL from stdin, edit and write to stdout, easily pipe and combine other commands
- Keep comments: You can update lots of existing HCL files with automation scripts
- Schemaless: independent of specific HCL applications
- HCL2 support (not HCL1)

The hcledit focuses on editing HCL with command line, doesn't aim for generic query tools. It was originally born for refactoring Terraform configurations, but it's not limited to specific applications.
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
A command line editor for HCL

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
  rm          Remove attribute
  set         Set attribute

Flags:
  -h, --help   help for attribute

Use "hcledit attribute [command] --help" for more information about a command.
```

Given the following file:

```attr.hcl
resource "foo" "bar" {
  attr1 = "val1"
  nested {
    attr2 = "val2"
  }
}
```

```
$ cat tmp/attr.hcl | hcledit attribute get resource.foo.bar.nested.attr2
"val2"
```

```
$ cat tmp/attr.hcl | hcledit attribute set resource.foo.bar.nested.attr2 '"val3"'
resource "foo" "bar" {
  attr1 = "val1"
  nested {
    attr2 = "val3"
  }
}
```

```
$ cat tmp/attr.hcl | hcledit attribute rm resource.foo.bar.attr1
resource "foo" "bar" {
  nested {
    attr2 = "val2"
  }
}
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
  rm          Remove block

Flags:
  -h, --help   help for block

Use "hcledit block [command] --help" for more information about a command.
```

Given the following file:

```block.hcl
resource "foo" "bar" {
  attr1 = "val1"
}

resource "foo" "baz" {
  attr1 = "val2"
}
```

```
$ cat tmp/block.hcl | hcledit block list
resource.foo.bar
resource.foo.baz
```

```
$ cat tmp/block.hcl | hcledit block get resource.foo.bar
resource "foo" "bar" {
  attr1 = "val1"
}
```

```
$ cat tmp/block.hcl | hcledit block mv resource.foo.bar resource.foo.qux
resource "foo" "qux" {
  attr1 = "val1"
}

resource "foo" "baz" {
  attr1 = "val2"
}
```

```
$ cat tmp/block.hcl | hcledit block rm resource.foo.baz
resource "foo" "bar" {
  attr1 = "val1"
}

```

## License

MIT
