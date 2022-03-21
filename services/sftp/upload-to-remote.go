package sftp

import (
	"io"
	"path/filepath"

	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
)

type FileDetails struct {
	SourceFile            io.Reader
	DestinationAddress    string
	FileNameAtDestination string
}

func (s SSH) UploadToRemote(log *logrus.Entry, fileDetails FileDetails) error {
	if s.Port <= 0 {
		s.Port = 22
	}
	sshClient, dialErr := s.dial()
	if dialErr != nil {
		if log != nil {
			log.WithError(dialErr).Error("error while dial for ssh connection")
		}
		return dialErr
	}
	defer func() {
		if err := sshClient.Close(); err != nil {
			// close instance error
			if log == nil {
				log = logrus.WithField("function_name", "UploadToRemote")
			}
			log.WithError(err).Error("error while closing ssh connection")
		}
	}()

	// create new SFTP client
	client, sftpErr := sftp.NewClient(sshClient)
	if sftpErr != nil {
		if log != nil {
			log.WithError(sftpErr).Error("error while closing ssh connection")
		}
		return sftpErr
	}
	defer func() {
		if closeClientErr := client.Close(); closeClientErr != nil {
			// close instance error
			if log == nil {
				log = logrus.WithField("function_name", "UploadToRemote")
			}
			log.WithError(closeClientErr).Error("error while closing sftp connection")
		}
	}()
	fileFullPath := filepath.Join(fileDetails.DestinationAddress, fileDetails.FileNameAtDestination)
	// create destination file
	dstFile, err := client.Create(fileFullPath)
	if err != nil {
		if log != nil {
			log.WithError(err).Errorf("error while creating file(%v) on remote", fileFullPath)
		}
		return err
	}
	defer func() {
		if closeDstFileErr := dstFile.Close(); closeDstFileErr != nil {
			// close instance error
			if log == nil {
				log = logrus.WithField("function_name", "UploadToRemote")
			}
			log.WithError(closeDstFileErr).Error("error while closing remote file instance")
		}
	}()

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, fileDetails.SourceFile)
	if err != nil {
		if log != nil {
			log.WithError(err).Error("error while copying file on remote")
		}
		return err
	}
	if log != nil {
		log.WithField("fileName", fileFullPath).Infof("%d bytes copied", bytes)
	}
	return nil
}

// func getHostKey(host, sshPathForKnownHostsFile string) (ssh.PublicKey, error) {
// 	// parse OpenSSH known_hosts file
// 	// ssh or use ssh-keyscan to get initial key
// 	if sshPathForKnownHostsFile == "" {
// 		sshPathForKnownHostsFile = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
// 	}
// 	file, err := os.Open(sshPathForKnownHostsFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	var hostKey ssh.PublicKey
// 	for scanner.Scan() {
// 		fields := strings.Split(scanner.Text(), " ")
// 		if len(fields) != 3 {
// 			continue
// 		}
// 		if strings.Contains(fields[0], host) {
// 			var err error
// 			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
// 			if err != nil {
// 				return nil, errors.Errorf("error parsing %q: %v", fields[2], err)
// 			}
// 			break
// 		}
// 	}

// 	if hostKey == nil {
// 		return nil, errors.Errorf("no hostkey found for %s", host)
// 	}
// 	return hostKey, err
// }
