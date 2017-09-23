# File Maps

![File Maps Logo](https://codeboy.fi/filemaps_logo.png)

[![Build Status](https://travis-ci.org/filemaps/filemaps.svg?branch=master)](https://travis-ci.org/filemaps/filemaps)
[![MPLv2 License](https://img.shields.io/badge/license-MPLv2-blue.svg?style=flat-square)](https://www.mozilla.org/MPL/2.0/)

Tool for visually organizing project files.

## Building

First, install [glide][2], vendor package management.

    glide install
    go run build.go

## Development

### Development with Angular CLI

1. `ng serve --base-href /ui/`
2. `filemaps -no-browser -cors-allow-origin "*"`
3. Open your browser on http://localhost:4200/ui

## License

File Maps is licensed under the [MPLv2 License][1].

[1]: https://github.com/filemaps/filemaps/blob/master/LICENSE
[2]: https://github.com/Masterminds/glide
