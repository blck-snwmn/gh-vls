A gh extension to display a list of alerts that are `OPEN` in the dependabot security alerts.
## Create
```bash
$ gh extension create --precompiled=go gh-vls
```

## Build
```bash
go build
```

## Install
local
```bash
$ gh extension install .
```

from repository
```bash
$ gh extension install https://github.com/blck-snwmn/gh-vls
```

## Run
```bash
$ gh vls
```

## Update
```bash
$ git tag -a v1.1.0 -m 'update'
```