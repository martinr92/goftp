package goftp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// FtpStatus status codes by RFC 959
// https://tools.ietf.org/html/rfc959
type FtpStatus int

const (
	// FtpStatusFileOK - File status okay; about to open data connection.
	FtpStatusFileOK FtpStatus = 150

	// FtpStatusReadyForNewUser - Service ready for new user.
	FtpStatusReadyForNewUser FtpStatus = 220

	// FtpStatusClosingDataConnection - Closing data connection.
	FtpStatusClosingDataConnection FtpStatus = 226

	// FtpStatusEnteringPassiveMode - Entering Passive Mode (h1,h2,h3,h4,p1,p2).
	FtpStatusEnteringPassiveMode FtpStatus = 227

	// FtpStatusLoginOK - User logged in, proceed.
	FtpStatusLoginOK FtpStatus = 230

	// FtpStatusFileActionOK - Requested file action okay, completed.
	FtpStatusFileActionOK FtpStatus = 250

	// FtpStatusPathCreated - Path created.
	FtpStatusPathCreated FtpStatus = 257

	// FtpStatusUserNameOK - User name okay, need password.
	FtpStatusUserNameOK FtpStatus = 331

	// FtpStatusLoginIncorrect - Not logged in.
	FtpStatusLoginIncorrect FtpStatus = 530

	// FtpStatusRequestedActionNotTaken - Requested action not taken.
	// File unavailable (e.g., file not found, no access).
	FtpStatusRequestedActionNotTaken FtpStatus = 550
)

// FtpError is a custom error struct for FTP communication errors.
type FtpError struct {
	ExpectedStatusCodes []FtpStatus
	ServerResponse      string
}

func (ftpError *FtpError) Error() string {
	errorString := "invalid server response!"
	errorString += "\n expected status codes: "
	for i, statusCode := range ftpError.ExpectedStatusCodes {
		if i != 0 {
			errorString += ","
		}
		errorString += strconv.Itoa(int(statusCode))
	}
	errorString += "\n server response: " + ftpError.ServerResponse
	return errorString
}

// Ftp object for remote connection.
type Ftp struct {
	remoteAddr string
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
}

// NewFtp creates a new FTP connection object.
//
// Examples:
//		NewFtp("my.host.de:1234")
func NewFtp(remote string) (*Ftp, error) {
	// try to connect on remote host
	conn, err := net.Dial("tcp", remote)
	if err != nil {
		return nil, err
	}

	// initialize reader and writer buffer for communication
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// create new ftp connection object
	ftp := &Ftp{remoteAddr: remote, connection: conn, reader: reader, writer: writer}

	// ignore welcome message
	ftp.read()

	return ftp, nil
}

// Login sends credentails to the FTP server and verifies the server login response status.
func (ftp *Ftp) Login(user string, password string) error {
	// send username
	_, _, err := ftp.writeCommand("USER "+user, []FtpStatus{FtpStatusUserNameOK})
	if err != nil {
		return err
	}

	// send password
	_, _, err = ftp.writeCommand("PASS "+password, []FtpStatus{FtpStatusLoginOK})
	if err != nil {
		return err
	}

	return err
}

// OpenDirectory changes the current working directory.
func (ftp *Ftp) OpenDirectory(directory string) error {
	// send new directory path
	_, _, err := ftp.writeCommand("CWD "+directory, []FtpStatus{FtpStatusFileActionOK})
	if err != nil {
		return err
	}

	// great!
	return nil
}

// CreateDirectory tells the server, to create a new named directory.
func (ftp *Ftp) CreateDirectory(directory string) error {
	// send new folder path
	_, _, err := ftp.writeCommand("MKD "+directory, []FtpStatus{FtpStatusPathCreated})
	if err != nil {
		return err
	}

	// great!
	return nil
}

// Upload a file to the remote server.
// Make shure, that the correct directory is already open!
// You can use OpenDirectory to change in the directory you want to use.
func (ftp *Ftp) Upload(localFilePath string, remoteFilePath string) error {
	// get passive connection data
	port, err := ftp.passiveConnection()
	if err != nil {
		return err
	}

	// open passive connection
	host := strings.Split(ftp.remoteAddr, ":")[0]
	passiveRemoteAddr := host + ":" + strconv.Itoa(port)
	fmt.Println("open passive connection:", passiveRemoteAddr)
	passiveConn, err := net.Dial("tcp", passiveRemoteAddr)
	if err != nil {
		return err
	}
	defer passiveConn.Close()

	// send store request
	_, _, err = ftp.writeCommand("STOR "+remoteFilePath, []FtpStatus{FtpStatusFileOK})
	if err != nil {
		return err
	}

	// open local file
	localFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}

	// send data to remote server
	// TODO: do some custom stuff instead for bandwidth limitation and progress
	_, err = io.Copy(passiveConn, localFile)
	if err != nil {
		return err
	}
	passiveConn.Close()

	// check master connectio nstatus
	_, _, err = ftp.readCommand([]FtpStatus{FtpStatusClosingDataConnection})
	if err != nil {
		return err
	}

	return nil
}

// Close quits the connection.
func (ftp *Ftp) Close() {
	ftp.connection.Close()
	fmt.Println("Connection closed")
}

func (ftp *Ftp) checkTextStatus(text string, statusCodes []FtpStatus) (matchedStatusCode *FtpStatus, err error) {
	for _, statusCode := range statusCodes {
		statusCodeString := strconv.Itoa(int(statusCode))
		found := strings.HasPrefix(text, statusCodeString)
		if found {
			return &statusCode, nil
		}
	}

	err = &FtpError{ServerResponse: text, ExpectedStatusCodes: statusCodes}
	return nil, err
}

func (ftp *Ftp) passiveConnection() (int, error) {
	responseText, _, err := ftp.writeCommand("PASV", []FtpStatus{FtpStatusEnteringPassiveMode})
	if err != nil {
		return 0, err
	}

	// parse connection port
	data := strings.Split(responseText, "(")[1]
	ipPortDataString := strings.Split(data, ")")[0]
	ipPortData := strings.Split(ipPortDataString, ",")

	// get port parts
	portPart1String := ipPortData[4]
	portPart2String := ipPortData[5]
	fmt.Println("Port Part 1:", portPart1String, "; Port Part 2:", portPart2String)

	// calculate port number
	portPart1, err := strconv.Atoi(portPart1String)
	portPart2, err := strconv.Atoi(portPart2String)
	port := portPart1*256 + portPart2
	fmt.Println("Calculated Port:", port)

	return port, nil
}

func (ftp *Ftp) read() (string, error) {
	fmt.Println("start reading...")
	text, err := ftp.reader.ReadString('\n')
	if err != nil {
		return text, err
	}

	fmt.Println("read: >>", text, "<<")
	return text, err
}

func (ftp *Ftp) readCommand(expectedStatusCodes []FtpStatus) (responseText string, matchedStatusCode *FtpStatus, err error) {
	// read server response
	responseText, err = ftp.read()
	if err != nil {
		return responseText, nil, err
	}

	// check response code
	matchedStatusCode, err = ftp.checkTextStatus(responseText, expectedStatusCodes)
	if err != nil {
		return responseText, matchedStatusCode, err
	}

	return responseText, matchedStatusCode, nil
}

func (ftp *Ftp) write(command string) error {
	fmt.Println("write: >>", command, "<<")
	_, err := ftp.writer.WriteString(command + "\n")
	ftp.writer.Flush()
	fmt.Println("write executed")
	return err
}

func (ftp *Ftp) writeCommand(command string, expectedStatusCodes []FtpStatus) (responseText string, matchedStatusCode *FtpStatus, err error) {
	// send command
	err = ftp.write(command)
	if err != nil {
		return "", nil, err
	}

	// check server response
	responseText, matchedStatusCode, err = ftp.readCommand(expectedStatusCodes)
	if err != nil {
		return responseText, matchedStatusCode, err
	}

	// everything was working great
	return responseText, matchedStatusCode, nil
}
