# goftp
[![Build Status](https://travis-ci.org/martinr92/goftp.svg?branch=master)](https://travis-ci.org/martinr92/goftp)

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
