# File Maps Backend

[![Build Status](https://travis-ci.org/filemaps/filemaps-backend.svg?branch=master)](https://travis-ci.org/filemaps/filemaps-backend)
[![MPLv2 License](https://img.shields.io/badge/license-MPLv2-blue.svg?style=flat-square)](https://www.mozilla.org/MPL/2.0/)

## Building

First, install [glide][2], vendor package management.

    glide install
    go run build.go

### Windows

In Windows you need [mingw-w64 gcc][3] for compiling required [sqlite3][4].

## License

Backend component is licensed under the [MPLv2 License][1].

[1]: https://github.com/filemaps/filemaps-backend/blob/master/LICENSE
[2]: https://github.com/Masterminds/glide
[3]: https://mingw-w64.org/
[4]: https://github.com/mattn/go-sqlite3
