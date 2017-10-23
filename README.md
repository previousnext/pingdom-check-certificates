Pingdom: Check Certificate
==========================

[![CircleCI](https://circleci.com/gh/previousnext/pingdom-check-certificates.svg)](https://circleci.com/gh/previousnext/pingdom-check-certificates)

**Maintainer**: Nick Schuch

## Overview

Checks the certificate of each host registered as a check in Pingdom.

## Development

### Tools

* **Dependency management** - https://github.com/golang/dep
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint

### Workflow

**Installing a new dependency**

```bash
dep ensure -add github.com/foo/bar
```

**Running quality checks**

```bash
make lint test
```

**Building binaries**

```bash
make build
```
