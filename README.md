# SuperSet

![Build](https://github.com/RyazMax/SuperSet/workflows/Build/badge.svg)
[![made-with-latex](https://img.shields.io/badge/Made%20with-LaTeX-1f425f.svg)](https://www.latex-project.org/)
[![GitHub issues](https://img.shields.io/github/issues/RyazMax/SuperSet.svg)
](https://github.com/RyazMax/SuperSet/issues/)

### Course work in Bauman Moscow State Technical University

## Terms of reference
See `./tor` folder 

## Start project

Too start project use next commands:

```bash
make start_tarantool
```

after that open another terminal and run:
```bash
go run src/main.go
```

To run tests start tarantool as described above and run `make tests`

>>> Current version of tests is not good enough, so do `make clean` every time you run it
