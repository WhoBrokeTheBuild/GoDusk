[![Build Status](https://travis-ci.org/WhoBrokeTheBuild/GoDusk.svg?branch=master)](https://travis-ci.org/WhoBrokeTheBuild/GoDusk)

# Building

GoDusk relies on `go-bindata`, which can be automatically run with `go generate`.
After pulling, run the following:

```
go generate ./...
```

Then you can build and use the library normally.

```
cd example
make run
```

Using the Makefile is only required if you want the Git Short linked in.
