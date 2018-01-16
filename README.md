# OpenID Connect Token Ferry

Performs an OpenID Connect Authentication Flow from the command line using an
"out of band" redirect URL. The OpenID Connect Issuer will provide a "code"
after the user has been authenticated. That code needs to be fed into this CLI.

This little project has been inspired by [k8s-oidc-helper][koh], which also
solves this problem, but specifically for Google as an Issuer.

[koh]: https://github.com/micahhausler/k8s-oidc-helper

## Usage

    Usage:
    oidc-token-ferry [OPTIONS]

    OpenID Connect:
        --oidc.issuer-url=    Issuer URL (default: https://accounts.google.com)
        --oidc.client-id=     Client ID to be used
        --oidc.client-secret= Client Secret to be used

    Rendering:
        --render.go-template= Go Template used to render credentials

    Help Options:
    -h, --help                Show this help message

## Building

There's a `Makefile` that'll build a statically linked linux amd64 binary
using Docker.

## License

    MIT License

    Copyright (c) 2018 Tom Wieczorek

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
