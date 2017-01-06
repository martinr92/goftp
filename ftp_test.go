package goftp

import "testing"

func TestSimpleUpload(t *testing.T) {
	// connect to remote server
	ftp, err := NewFtp("localhost:21")
	if err != nil {
		t.Error("connection failed!", err)
	}

	// send user credentials
	err = ftp.Login("test", "test")
	if err != nil {
		t.Error("logon failed", err)
	}

	// open remote directory
	// TODO: we should bring this step back, when we're able to create directories
	/*err = ftp.OpenDirectory("/test/")
	if err != nil {
		t.Error("change directory failed!", err)
	}*/

	// upload a local file to the remote FTP server
	err = ftp.Upload("README.md", "README.md")
	if err != nil {
		t.Error("file upload failed!", err)
	}
}
