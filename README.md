MylSlurper
==========

MylSlurper is a local SMTP server with a web UI for capturing and inspecting email during development. Point your app at it instead of a real mail server, then browse, search, and prune the messages it slurps.

This project is inspired by [MailSlurper](https://github.com/mailslurper/mailslurper). It is an independent rewrite: Go 1.26, SQLite without CGO, and a vanilla HTML/CSS/JS UI embedded in the binary.

Compiling
---------

MylSlurper is a single Go module with no CGO dependencies and no separate
frontend build step — the web UI is plain HTML/CSS/JS, embedded into the
binary at compile time via `go:embed`.

* Go 1.26 or higher
* Git

```bash
git clone git@github.com:p0vidl0/mylslurper.git
cd mylslurper
go build ./cmd/mylslurper
go build ./cmd/createcredentials
```

Run it with `go run ./cmd/mylslurper` or `./mylslurper` (reads `config.json`
in the working directory if present; see `cmd/mylslurper/config.json` for a
sample). Every field can also be set via environment variables such as
`HTTP_PORT` / `SMTP_PORT` / `DB_FILE` / `AUTH_SCHEME`.

Quickstart With Docker
----------------------

```bash
docker run -it --rm --name mylslurper -p 4436:4436 -p 4437:4437 -p 1025:1025 ghcr.io/p0vidl0/mylslurper
```

The image is a drop-in replacement for `oryd/mailslurper:latest-smtps`:
port **4436** is the web UI, **4437** is the classic service REST API
(`GET /mailcount`, `GET /mail`), and **1025** is implicit TLS SMTP
(`smtps://localhost:1025/?skip_ssl_verify=true`). HTTP listeners stay
plain; only SMTP uses the baked-in self-signed certificate. Unset
`CERT_FILE` and `KEY_FILE` to disable SMTPS.

Scripts
-------

The server must be listening on `:1025` (SMTPS). Scripts live in `bin/`:

| Script | Purpose |
| --- | --- |
| `bin/send-one.py` | One plain-text message; optional To address |
| `bin/send-suite.py` | Batch of cases (attachments, XSS, Date headers, HTML) |
| `bin/send-mime.py` | A more complex MIME message |

```bash
go run ./cmd/mylslurper
python3 bin/send-one.py
python3 bin/send-one.py someone@example.com
python3 bin/send-suite.py
python3 bin/send-mime.py
```

Library and Framework Credits
-----------------------------

This application uses a handful of open source libraries.

* [bluemonday](https://github.com/microcosm-cc/bluemonday) - BSD 3-Clause
* [go-cache](https://github.com/patrickmn/go-cache) - MIT
* [golang-jwt](https://github.com/golang-jwt/jwt) - MIT
* [Logrus](https://github.com/sirupsen/logrus) - MIT
* [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - BSD 3-Clause (pure Go, no CGO)
* [open](https://github.com/skratchdot/open-golang) - MIT
* [Inter](https://github.com/rsms/inter) - SIL Open Font License 1.1
* [JetBrains Mono](https://github.com/JetBrains/JetBrainsMono) - SIL Open Font License 1.1

Everything else — the SMTP server, HTTP API, and web UI — is stdlib Go and
framework-free HTML/CSS/JS.

License
-------

MIT. Copyright (c) 2013-2018 Adam Presley. The full license text is in
[LICENSE](LICENSE).
