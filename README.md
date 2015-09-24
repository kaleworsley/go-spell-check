# go-spell-check

Spell checking for comments in go packages.

## Setup

```
go get github.com/kaleworsley/go-spell-check
```

## Usage

Check the spelling of the package in the current directory.

```
cd /path/to/package/src
go-spell-check
```

Check the spelling of the package in a given directory.

```
go-spell-check /path/to/package/src
```

Un-camelcase words before checking.

```
go-spell-check -camel 
```

Use a different dictionary language. (Default is `en_US`).

```
go-spell-check -lang="en_GB"
```

## License

Copyright 2015, Kale Worsley.

go-spell-check is made available under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
