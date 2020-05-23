# Examples of Dear ImGui for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/inkyblackness/imgui-go-examples)](https://goreportcard.com/report/github.com/inkyblackness/imgui-go-examples)
[![Lint Status](https://github.com/inkyblackness/imgui-go-examples/workflows/golangci-lint/badge.svg)](https://github.com/inkyblackness/imgui-go-examples/actions)

This project contains a set of [Go](https://www.golang.org) examples for [imgui-go](https://github.com/inkyblackness/imgui-go), which is a wrapper for [**Dear ImGui**](https://github.com/ocornut/imgui).

It provides reference implementations on how to use and integrate **Dear ImGui** in Go.

![Screenshot](assets/screenshot.png)

## Layout
The project follows the basic concept of the examples of **Dear ImGui** by separating platform and renderer bindings from the example applications that wire them together in compatible constellations.

* `cmd` contains the main functions of the example applications. They typically combine a platform with a renderer.
* `internal` contains the reusable library components
  * `platforms` contains code for mouse/keyboard/gamepad inputs, cursor shape, timing, windowing. For example based on: [GLFW3](https://github.com/go-gl/glfw) and [SDL2](https://github.com/veandco/go-sdl2). 
  * `renderers` contains code for creating the main font texture, rendering imgui draw data. For example using: [OpenGL](https://github.com/go-gl/gl) (both v2.1 (fixed pipe) and v3.2 (shaders)) 
  * `demo` contains the common example code.

## Running examples

Some platforms and/or renderers may require dedicated build flags to be passed.
Each example comes with a separate `README.md` that describes how to run them.

> Build flags are used in order to avoid compiling all the libraries at once.

## License

The project is available under the terms of the **New BSD License** (see LICENSE file).
