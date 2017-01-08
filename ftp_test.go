package goftp

import "testing"
import "strings"
import "strconv"

func TestPassiveDataConnection(t *testing.T) {
	// connect to remote server
	ftpClient, err := NewFtp(getConnectionString())
	if err != nil {
		t.Error("connection failed!", err)
		return
	}

	// don't forget to close the connection
	defer ftpClient.Close()

	// send invalid user name
	err = ftpClient.Login("", "test")
	if err == nil {
		t.Error("logon should fail, because the username is missing!", err)
		return
	}

	// send invalid logon data
	err = ftpClient.Login("asdf", "test")
	if err == nil {
		t.Error("logon should fail, because the username is invalid!", err)
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
		return
	}

	// try to open a remote directory, that doesn't exist
	if err = ftpClient.OpenDirectory("/home/test/testfolder/"); err == nil {
		t.Error("change directory should failed!")
		return
	}

	// create the directory
	if err = ftpClient.CreateDirectory("/home/test/testfolder/"); err != nil {
		t.Error("unable to create direcory!", err)
		return
	}

	// try to create the same directory again
	if err = ftpClient.CreateDirectory("/home/test/testfolder/"); err == nil {
		t.Error("directory already exists! This should fail!")
		return
	}

	// open directory
	if err = ftpClient.OpenDirectory("/home/test/testfolder/"); err != nil {
		t.Error("unable to open directory!", err)
		return
	}

	// upload a local file to the remote FTP server
	if err = ftpClient.Upload("README.md", "README.md"); err != nil {
		t.Error("file upload failed!", err)
		return
	}
}

func TestActiveDataConnection(t *testing.T) {
	// connect to remote server
	ftpClient, err := NewFtp(getConnectionString())
	if err != nil {
		t.Error("connection failed!", err)
		return
	}

	// don't forget to close the connection
	defer ftpClient.Close()

	// send user credentials
	err = ftpClient.Login("test", "test")
	if err != nil {
		t.Error("logon failed", err)
		return
	}

	// create the directory
	if err = ftpClient.CreateDirectory("/home/test/testfolder_active/"); err != nil {
		t.Error("unable to create direcory!", err)
		return
	}

	// open directory
	if err = ftpClient.OpenDirectory("/home/test/testfolder_active/"); err != nil {
		t.Error("unable to open directory!", err)
		return
	}

	// set active mode with an invalid IP
	ftpClient.ActiveMode = true
	ftpClient.ActiveModeIPv4 = "1.2.3.4"

	// upload a local file to the remote FTP server
	if err = ftpClient.Upload("README.md", "README.md"); err == nil {
		t.Error("file upload should fail!")
		return
	}

	// set active mode
	ftpClient.ActiveMode = true
	ftpClient.ActiveModeIPv4 = "127.0.0.1"

	// upload a local file to the remote FTP server
	if err = ftpClient.Upload("README.md", "README.md"); err != nil {
		t.Error("file upload failed!", err)
		return
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
