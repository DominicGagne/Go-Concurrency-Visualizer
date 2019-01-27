# Go-Concurrency-Visualizer
A Go runtime concurrency visualizer

![](assets/fanInOne.gif)

This tool is built on top of the `gotrace` tool created by Ivan Danyliuk.  The tool developed by Ivan is wonderful.  However, the original `gotrace` suffered from many usability issues, lack of useful error message outpout, and appears to now be unmaintained.
Original tool: https://github.com/divan/gotrace

This wrapper on top of `gotrace` aims to make the tool more accessible to those seeking to visualize concurrency in Go, and is usable out of the box.



![](assets/pongOne.gif)



Usage:

Copy a go source file into the root directory of this project, such as `hello.go`.
Type `./visualize hello.go`

The first time you run this tool, a docker image will be built from source, which may take a few minutes. Subsequent visualizations will be significantly faster.

