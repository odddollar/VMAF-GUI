# VMAF GUI

A clean and simple graphical user interface for [Netflix's VMAF](https://github.com/Netflix/vmaf) video comparison algorithm.

VMAF is extremely useful for comparing the perceptual quality of a compressed video to its source/reference, but the command to run it is often long and annoying to write. This program offers a reproducible way of running VMAF for multiple video files, with additional quality-of-life features.

## Features



## Building

This program is built using the [Go](https://go.dev/) programming language and the [Fyne](https://fyne.io/) UI framework. To setup a development environment, you'll need to install:

- [Go](https://go.dev/)
- A C compiler ([w64devkit](https://github.com/skeeto/w64devkit), [Cygwin](https://cygwin.com/), [MSYS2](https://www.msys2.org/), or similar)
- [Fyne's tooling](https://docs.fyne.io/started/packaging/)
  - Can be installed with `go install fyne.io/tools/cmd/fyne@latest`
- (Optional) [UPX](https://github.com/upx/upx)
  - Fyne projects can be large when compiled. Not necessary, but nice to have

Clone this repository:

```bash
git clone https://github.com/odddollar/VMAF-GUI.git
cd VMAF-GUI
```

Run for testing/development with:

```bash
go run .
```

Package for release with:

```bash
fyne package --release
```

(Optional) Use UPX to transparently compress the executable:

```bash
upx --ultra-brute "VMAF GUI.exe"
```

## Screenshots

<div align="center">
    <img src="./screenshots/Results.png" alt="Results"><br><br>
    <img src="./screenshots/Compare.png" alt="Compare">
</div>
