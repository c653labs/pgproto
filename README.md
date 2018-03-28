# pgproto

[![GoDoc](https://godoc.org/github.com/c653labs/pgproto?status.svg)](https://godoc.org/github.com/c653labs/pgproto)

Package `pgproto` is a pure Go protocol library for [PostgreSQL](https://www.postgresql.org/).
It provides the necessary structures and functions to parse and encode PostgreSQL messages.

The scope of `pgproto`` is only for parsing/encoding messages and does not handle connections between PostgreSQL client and server.


**NOTE:** `pgproto` is still under active development, while it will work for some basic use cases not everything is implemented yet.
