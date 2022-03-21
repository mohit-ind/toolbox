package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/toolboxservices/sftp"
)

func main() {
	fmt.Println((sftp.SSH{
		Host:     `test.rebex.net`,
		Port:     22,
		User:     "demo",
		Password: "password",
	}).UploadToRemote(logrus.WithField("Testing", "upload to remote"), sftp.FileDetails{
		DestinationAddress:    "./Desktop",
		FileNameAtDestination: "appventurez_test_forSFTP.xlsx",
	}))
}
