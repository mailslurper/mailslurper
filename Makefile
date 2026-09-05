# Test mail scripts. The SMTP listener must already be on :1025 (SMTPS)
# (go run ./cmd/mylslurper).
#
#   make send-one
#   make send-one TO=someone@example.com
#   make send-suite
#   make send-mime
#   make send-all

PYTHON ?= python3
TO     ?=

.PHONY: help test test-js coverage send-one send-suite send-mime send-all

help:
	@echo "test        run Go tests (with race detector)"
	@echo "test-js     run frontend unit tests"
	@echo "coverage    Go test coverage summary"
	@echo "send-one    one plain-text message  (TO=addr, optional)"
	@echo "send-suite  attachments, XSS, Date headers, HTML"
	@echo "send-mime   a more complex MIME message"
	@echo "send-all    run all three"

test:
	go test ./... -race -count=1

test-js:
	node --test $$(find web/static/js -name '*.test.js')

coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out | tail -1

send-one:
	$(PYTHON) bin/send-one.py $(TO)

send-suite:
	$(PYTHON) bin/send-suite.py

send-mime:
	$(PYTHON) bin/send-mime.py

send-all: send-one send-suite send-mime
