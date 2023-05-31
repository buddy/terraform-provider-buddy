<img alt="Terraform" src="https://www.datocms-assets.com/2885/1629941242-logo-terraform-main.svg" width="600px">

Terraform Provider for Buddy
=============================

- [Documentation](https://www.terraform.io/docs/providers/buddy/index.html)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0.11
- [Go](https://golang.org/doc/install) >= 1.19 (to build the provider plugin)

## Developing

```sh
$ go install
```

### Linter
```sh
$ make lint
```

### Running Tests

```sh
$ BUDDY_TOKEN=example123 BUDDY_BASE_URL=https://api.buddy.works make test
```
