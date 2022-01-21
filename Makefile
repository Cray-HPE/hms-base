
all:  unittest
.PHONY:  unittest

unittest:
	go test ./... -cover
