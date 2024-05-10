# SSH Portfolio

A small fun project inspired by [terminal.shop](https://www.terminal.shop/) to create an SSH-based portfolio website to learn more about SSH and the Go programming language.

Run `ssh ssh-portfolio.fly.dev` to see this in action!

## Deployment

Todo - this is a work in progress and I still need to figure out how to deploy this application and hook it up to my domain [dannyisaac.com](dannyissac.com)

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
