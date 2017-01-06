package goftp

import "testing"
import "strings"
import "strconv"

func TestSimpleUpload(t *testing.T) {
	// connect to remote server
	ftpClient, err := NewFtp(getConnectionString())
	if err != nil {
		t.Error("connection failed!", err)
	}

	// don't forget to close the connection
	defer ftpClient.Close()

	// send invalid user name
	err = ftpClient.Login("asdf", "test")
	if err == nil {
		t.Error("logon should fail!", err)
		return
	}

	// check response code
	if ftpError, ok := err.(*FtpError); ok {
		if !strings.HasPrefix(ftpError.ServerResponse, strconv.Itoa(int(FtpStatusLoginIncorrect))) {
			t.Error("unexpected response!", ftpError)
			return
		}

		// yea, the server said, that the login was incorrect
		if !strings.HasPrefix(ftpError.Error(), "invalid server response!") {
			t.Error("invalid error text response!")
			return
		}
	} else {
		t.Error("invalid error returned!")
		return
	}

	// send user credentials
	err = ftpClient.Login("test", "test")
	if err != nil {
		t.Error("logon failed", err)
	}

	// open remote directory
	// TODO: we should bring this step back, when we're able to create directories
	/*err = ftpClient.OpenDirectory("/test/")
	if err != nil {
		t.Error("change directory failed!", err)
	}*/

	// upload a local file to the remote FTP server
	err = ftpClient.Upload("README.md", "README.md")
	if err != nil {
		t.Error("file upload failed!", err)
	}
}

func TestLogonError(t *testing.T) {
	// connect to remote server
	ftpClient, err := NewFtp("host.not.exists:21")
	if err == nil {
		// lose the connection
		ftpClient.Close()

		t.Error("Error expected! invalid host name")
		return
	}
}

func getConnectionString() string {
	return "localhost:21"
}
