# buildimg

[Docker](https://docker.com/) now has multi-platform builds via BuildKit (the `buildx` subcommand). For folks, like myself, working on a Macbook running Apple Silicon, this is both a blessing and a curse. It's a blessing because it makes it possible to build `linux/amd64` images even though I'm running on an ARM CPU. It's a curse because it's a complicated Inception-esque stack of virtualization. `buildimg` is a simple command I use to automate building and publishing multi-platform images. In particular, I have a number of projects where the containers run on `linux/amd64` but are developed on my laptop, which natively produces `linux/arm64`.

## Installation

`buildimg` is a plain ole [Go](https://golang.org/) and can be installed using `go install` kind of like this:

```
go install github.com/kellegous/buildimg@latest
```

## Examples

The best way to explain how `buildimg` works is to just look at some examples.

The following command will build and push a new `linux/amd64` image for `kellegous/example`. The tag for the build will be determined from the git directory in the current working directory and so this will push an image that looks similar to `kellegous/exmaple:2193a7f6`.

```
$ buildimg --target=linux/amd64 kellegous/example
```

The next example is similar except it will push both a `linux/amd64` image and a `linux/arm64` image. Like the previous command, the image will be tagged with the latest SHA from the git directory in the current working directory.

```
$ buildimg --target=linux/amd64 --target=linux/arm64 kellegous/example
```

Next, instead of relying on git for the tag, let's explicitly tag the image with `latest`.

```
$ buildimg --target=linux/amd64 --target=linux/arm64 --tag=latest kellegous/example
```

What if we do not wish to push the images at all? The following command will build a `linux/amd64` image and export it to the file `example-amd64.tar`. That file can be imported into docker via `docker import` as it follows the [OCI Image Layout](https://github.com/opencontainers/image-spec/blob/main/image-layout.md).

```
$ buildimg --target=linux/amd64:example-amd64.tar kellegous/example
```

Again, multiple `--target` flags can be specified, each with a destination file that allows one to create a multi-platform build where each platform is saved in an OCI image tarball. For instance, this command creates `linux/amd64` build and stores it in `example-amd64.tar` and also creates a `linux/arm64` build and stores it in the file `example-arm64.tar`.

```
$ buildimg --target=linux/amd64:example-amd64.tar --target=linux/arm64:example-arm64.tar kellegous/example
```

Generally, any target that specifies a destination path (i.e. `linux/arm64:example.tar`) will be exported to a local tarbar and any target that does not specify a destination path will be pushed. So for instance, the following command exports a `linux/arm64` build to the local file `example.tar` and also pushes a `linux/amd64` build to the [docker.io](http://hub.docker.com/) registry.

```
$ buildimg --target=linux/arm64:example.tar --target=linux/amd64 kellegous/example
```

Note that it is possible to push a particular platform and also export it locally with the same command. In this example, a `linux/amd64` image is both pushed to the [docker.io](https://hub.docker.com/) registry and exported locally to the file `example.tar`.

```
$ buildimg --target=linux/amd64 --target=linux/amd64:example.tar kellegous/example
```

## Authors
 - Kelly Norton [kellegous.com](https://kellegous.com/about)
