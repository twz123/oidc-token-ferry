# OpenID Connect Token Ferry

[![Build](https://github.com/twz123/oidc-token-ferry/workflows/Build/badge.svg)](https://github.com/twz123/oidc-token-ferry/actions?query=workflow%3ABuild)

Performs an OpenID Connect Authentication Flow from the command line using an
"out of band" redirect URL. The OpenID Connect Issuer will provide a "code"
after the user has been authenticated. That code needs to be fed into this CLI.

This little project has been inspired by [k8s-oidc-helper][koh], which also
solves this problem, but specifically for Google as Identity Provider.

[koh]: https://github.com/micahhausler/k8s-oidc-helper

## Usage

General usage:

    Usage:
    oidc-token-ferry [OPTIONS] <command>

    Help Options:
    -h, --help  Show this help message

    Available commands:
    patch-kubeconfig    patches Kubernetes kubeconfig files
    render-go-template  renders credentials using Go Templates
    render-json         renders credentials as JSON
    version             Show oidc-token-ferry version information

How to patch a kubeconfig:

    Usage:
    oidc-token-ferry [OPTIONS] patch-kubeconfig [patch-kubeconfig-OPTIONS] [KUBECONFIG_FILE] [OUTPUT_FILE]

    Help Options:
    -h, --help                 Show this help message

    [patch-kubeconfig command options]
            --user-name=       User name to use when generating client configuration. Either user-name or user-claim-name may be specified.
            --user-claim-name= Claim that defines the user name to use when generating client configuration. Either user-name or user-claim-name may be specified.
            --no-open-url      Don't open the redirect URL in a browser automatically

        OpenID Connect Options:
        -u, --issuer-url=      IdP Issuer URL to be contacted (default: https://accounts.google.com)
        -i, --client-id=       Client ID to be used
        -s, --client-secret=   Client Secret to be used
        -r, --redirect-url=    Redirect URL to be communicated to the IdP (needs to indicate "out of band") (default: urn:ietf:wg:oauth:2.0:oob)
        -c, --claim=           Additional claims to be requested

    [patch-kubeconfig command arguments]
    KUBECONFIG_FILE:           Path to the kubeconfig file to be patched. Uses the default discovery mechanism if omitted/empty. Special value '-' (hyphen) means read from STDIN.
    OUTPUT_FILE:               Path to the patched kubeconfig file to be written. Overwrites kubeconfig if omitted/empty. Special value '-' (hyphen) means write to STDOUT.

How to render credentials via go-template:

    Usage:
    oidc-token-ferry [OPTIONS] render-go-template [render-go-template-OPTIONS] TEMPLATE_STRING

    Help Options:
    -h, --help                 Show this help message

    [render-go-template command options]
        -o, --output-file=     Output file to write (defaults to STDOUT if omitted)
            --no-open-url      Don't open the redirect URL in a browser automatically

        OpenID Connect Options:
        -u, --issuer-url=      IdP Issuer URL to be contacted (default: https://accounts.google.com)
        -i, --client-id=       Client ID to be used
        -s, --client-secret=   Client Secret to be used
        -r, --redirect-url=    Redirect URL to be communicated to the IdP (needs to indicate "out of band") (default: urn:ietf:wg:oauth:2.0:oob)
        -c, --claim=           Additional claims to be requested

    [render-go-template command arguments]
    TEMPLATE_STRING:           Go Template to be rendered. An empty template indicates that the template is to be read from STDIN.

## Building

    make

This will build a statically linked binary.

## License

    MIT License

    Copyright (c) 2018-2022 Tom Wieczorek

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
