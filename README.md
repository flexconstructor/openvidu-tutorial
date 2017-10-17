# OpenViDu-tutorial


Is a GoLang port of  [openvidu-mvc-java](http://openvidu.io/docs/tutorials/openvidu-mvc-java/) application.

A secure [OpenViDu](http://openvidu.io/) sample app with a GoLang backend and a traditional MVC frontend.

## How to dev

You should ensure that your project is cloned into `$GOPATH/src/github.com/flexconstructor/openvidu-tutorial` directory, o
therwise Golang tools may work with code incorrectly.

The simple start is:
```bash
make deps
make run
make run.goconvey
```

Use `docker-compose` to boot up (or restart) [dockerized environment][2] for development:
```bash
make build
docker-compose up --build

# or in one command
make run
```

To resolve project dependencies use docker-wrapped commands from [`Makefile`][1]:
```bash
make deps

# or concrete type
make deps.glide cmd=update
make deps.tools
```

To run tests or lint project use docker-wrapped commands from [`Makefile`][1]:
```bash
make test
make lint
```

To format project sources use docker-wrapped command from [`Makefile`][1]:
```bash
make fmt
```

To run GoConvey Web UI for continuous testing use docker-wrapped command from [`Makefile`][1] and access it on `8080` port:
```bash
make run.goconvey
# available on http://localhost:8080/
```

Take a look at [`Makefile`][1] for command usage details.

## Toolchain overview

The following Golang tools are used: 
- [Glide][10] for dependencies management.
- [Go Meta Linter][11] for code linting.
- [GoConvey][12] as testing framework.

To have fully reproducible builds and runtime environment [Docker][13] is used.


[1]: Makefile
[2]: docker-compose.yml
[10]: https://github.com/Masterminds/glide
[11]: https://github.com/alecthomas/gometalinter
[12]: https://github.com/smartystreets/goconvey
[13]: https://www.docker.com
