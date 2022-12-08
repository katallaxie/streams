# :partying_face: Template Go

## Get Started

```bash
go mod init github.com/ionos-cloud/REPO && go mod tidy
```

Features

* [Developer Containers](https://code.visualstudio.com/docs/remote/containers)
* [Editorconfig](https://editorconfig.org)
* [GoReleaser](https://goreleaser.com)
* [Git Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks#_git_hooks)
* [minikube](https://minikube.sigs.k8s.io/docs/start/)

> You can `sh scripts/postCreateCommand.sh` if you are not running in a remote container or on [Codespaces](https://github.com/features/codespaces).
> The template uses [`run`](https://github.com/katallaxie/run) as an alterantive to Makefiles

## Development

This template supports `Makefile` and `run` build tooling.

```bash
# build `main.go`
make 
```

Other available targets are

* `build`
* `fmt`
* `lint`
* `vet`
* `generate`
* `clean`

The convention is to use `make` to run the build.

Happy coding!
