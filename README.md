# goftp
[![Build Status](https://travis-ci.org/martinr92/goftp.svg?branch=master)](https://travis-ci.org/martinr92/goftp)
[![GoDoc](https://godoc.org/github.com/martinr92/goftp?status.svg)](https://godoc.org/github.com/martinr92/goftp)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinr92/goftp)](https://goreportcard.com/report/github.com/martinr92/goftp)
[![codecov](https://codecov.io/gh/martinr92/goftp/branch/master/graph/badge.svg)](https://codecov.io/gh/martinr92/goftp)

Goftp is a simple FTP library written in golang.
The implementation is based on the [RFC 959 - FILE TRANSFER PROTOCOL (FTP)](https://tools.ietf.org/html/rfc959)

## Features
* active and passive data connection mode

## Usage
Download the package and import it into your project.
```golang
import ftp "github.com/martinr92/goftp"
```

Connect to the remote server.
```golang
ftpClient, err := ftp.NewFtp("host.local:51000")
if err != nil {
    panic(err)
}
defer ftpClient.Close()
```

By default, the client uses a passive data connection for file transfer. If you want to use a active connection, just set the following:
```golang
ftpClient.ActiveMode = true
ftpClient.ActiveModeIPv4 = "1.2.3.4"
```

Send user credentials.
```golang
if err = ftpClient.Login("username", "password"); err != nil {
    panic(err)
}
```

Change the working directory.
```golang
if err = ftpClient.OpenDirectory("/some/folder/"); err != nil {
    panic(err)
}
```

Upload a file.
```golang
if err = ftpClient.Upload("/local/path/file.txt", "file.txt"); err != nil {
    panic(err)
}
```
