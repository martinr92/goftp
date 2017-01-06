# goftp
[![Build Status](https://travis-ci.org/martinr92/goftp.svg?branch=master)](https://travis-ci.org/martinr92/goftp)
[![GoDoc](https://godoc.org/github.com/martinr92/goftp?status.svg)](https://godoc.org/github.com/martinr92/goftp)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinr92/goftp)](https://goreportcard.com/report/github.com/martinr92/goftp)
[![codecov](https://codecov.io/gh/martinr92/goftp/branch/master/graph/badge.svg)](https://codecov.io/gh/martinr92/goftp)

Goftp is a simple FTP library written in golang.

## Usage
```golang
// connect to remote server
ftp, err := ftp.NewFtp("host.local:51000")
if err != nil {
    panic(err)
}

// send user credentials
err = ftp.Login("username", "password")
if err != nil {
    panic(err)
}

// open remote directory
err = ftp.OpenDirectory("/some/folder/")
if err != nil {
    panic(err)
}

// upload a local file to the remote FTP server
err = ftp.Upload("/local/file/to/path/file.txt", "file.txt")
if err != nil {
    panic(err)
}
```
