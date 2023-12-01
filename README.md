# github-fork-update

[![Go Version][go_version_img]][go_dev_url]
[![Go Report Card][go_report_img]][go_report_url]
[![Code Coverage][go_code_coverage_img]][repo_url]
[![License][repo_license_img]][repo_license_url]

## Description
github-fork-update is a tool to sync all forks in the owner's account from upstream.

## Contributing
Please see our [Contributing](./CONTRIBUTING.md) for how to contribute to the project.

## Setting up for development...
git clone https://github.com/mjdusa/github-fork-update.git

## Pre-commit Hooks
When you clone this repository to your workstation, make sure to install the [pre-commit](https://pre-commit.com/) hooks. [GitHub Repo](https://github.com/pre-commit/pre-commit)

### Installing tools
```
brew install pre-commit
```

### Check installed versions.
```
pre-commit --version
pre-commit 3.3.2
```

### Update configured pre-commit plugins.  Updates repo versions in .pre-commit-config.yaml to the latest.
```
pre-commit autoupdate
```

### Install pre-commit into the local git.
```
pre-commit install --install-hooks
```

### Run pre-commit checks manually.
```
pre-commit run --all-files
```

## Running...
```
make release
...
./dist/github-fork-update -auth=[github-auth-token]
```


<!-- Go -->

[go_download_url]: https://golang.org/dl/
[go_install_url]: https://golang.org/cmd/go/#hdr-Compile_and_install_packages_and_dependencies
[go_version_img]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[go_report_url]: https://goreportcard.com/report/github.com/mjdusa/github-fork-update
[go_code_coverage_img]: https://img.shields.io/badge/code_coverage-88%25-success?style=for-the-badge&logo=none
[go_dev_url]: https://pkg.go.dev/github.com/mjdusa/github-fork-update

<!-- Repository -->

[repo_url]: https://github.com/mjdusa/github-fork-update
[repo_logo_url]: https://github.com/mjdusa/github-fork-update/wiki/Logo
[repo_logo_img]: https://github.com/mjdusa/github-fork-update/assets/11155743/95024afc-5e3b-4d6f-8c9c-5daaa51d080d
[repo_license_url]: https://github.com/mjdusa/github-fork-update/blob/main/LICENSE
[repo_license_img]: http://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge&logo=none
[repo_cc_url]: https://creativecommons.org/licenses/by-sa/4.0/
[repo_v2_url]: https://github.com/mjdusa/github-fork-update/tree/v2
[repo_v3_url]: https://github.com/mjdusa/github-fork-update/tree/v3
[repo_issues_url]: https://github.com/mjdusa/github-fork-update/issues
[repo_pull_request_url]: https://github.com/mjdusa/github-fork-update/pulls
[repo_discussions_url]: https://github.com/mjdusa/github-fork-update/discussions
[repo_releases_url]: https://github.com/mjdusa/github-fork-update/releases
[repo_wiki_url]: https://github.com/mjdusa/github-fork-update/wiki
[repo_wiki_img]: https://img.shields.io/badge/docs-wiki_page-blue?style=for-the-badge&logo=none
[repo_wiki_faq_url]: https://github.com/mjdusa/github-fork-update/wiki/FAQ

<!-- Project -->

<!-- Author -->

[author]: https://github.com/mjdusa

<!-- Readme links -->

[dev_to_url]: https://dev.to/
[redis_url]: https://redis.io/
[postgresql_url]: https://postgresql.org/
[nginx_url]: https://nginx.org/
[traefik_url]: https://traefik.io/traefik/
[vitejs_url]: https://vitejs.dev/
[vuejs_url]: https://vuejs.org/
[react_url]: https://reactjs.org/
[preact_url]: https://preactjs.com/
[nextjs_url]: https://nextjs.org/
[nuxt3_url]: https://v3.nuxtjs.org/
[svelte_url]: https://svelte.dev/
[lit_url]: https://lit.dev/
[chi_url]: https://github.com/go-chi/chi
[fiber_url]: https://github.com/gofiber/fiber
[net_http_url]: https://golang.org/pkg/net/http/
[docker_url]: https://hub.docker.com/r/koddr/cgapp
[python_url]: https://www.python.org/downloads/
[ansible_url]: https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#installing-ansible-on-specific-operating-systems
[brew_url]: https://brew.sh/
