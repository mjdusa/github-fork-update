# github-fork-update

[![Go Version][go_version_img]][go_dev_url]
[![Go Report Card][go_report_img]][go_report_url]
[![Code Coverage][go_code_coverage_img]][repo_url]
[![License][repo_license_img]][repo_license_url]

## Description
github-fork-update is a tool to sync all forks in the owner's account from upstream.

## Contributing
Please see our [Contributing](./CONTRIBUTING.md) for how to contribute to the project.

## Setting up for development
```bash
git clone https://github.com/mjdusa/github-fork-update.git
```

## Pre-commit Hooks
When you clone this repository to your workstation, make sure to install the [pre-commit](https://pre-commit.com/) hooks. [GitHub repository](https://github.com/pre-commit/pre-commit)

### Installing tools
```bash
brew install pre-commit
brew install gitleaks
- or -
git clone https://github.com/gitleaks/gitleaks.git
cd gitleaks
make build
```

### Check installed versions
```bash
pre-commit --version
pre-commit 3.3.2
```

### Update configured pre-commit plugins
Updates repository versions in .pre-commit-config.yaml to the latest
```bash
pre-commit autoupdate
```

### Install pre-commit into the local git
```bash
pre-commit install --install-hooks
```

### Run pre-commit checks manually
```bash
pre-commit run --all-files
```

## Running
```bash
make release
...
./dist/github-fork-update -auth=[github-auth-token]
```

## Profiling

### Prereq
```bash
brew install graphviz
```

### Creating PDF for CPU profiling
```bash
go tool pprof -pdf ./dist/github-fork-update cpu-profile.pprof > cpu-profile.pdf
```

### Creating PDF for memory profiling
```bash
go tool pprof -pdf ./dist/github-fork-update mem-profile.pprof > mem-profile.pdf
```

<!-- Go -->

[go_version_img]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go
[go_dev_url]: https://pkg.go.dev/github.com/mjdusa/github-fork-update
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[go_report_url]: https://goreportcard.com/report/github.com/mjdusa/github-fork-update
[go_code_coverage_img]: https://img.shields.io/badge/code_coverage-92.6%25-success?style=for-the-badge&logo=none

<!-- Repository -->

[repo_url]: https://github.com/mjdusa/github-fork-update
[repo_license_url]: https://github.com/mjdusa/github-fork-update/blob/main/LICENSE
[repo_license_img]: http://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge&logo=none
