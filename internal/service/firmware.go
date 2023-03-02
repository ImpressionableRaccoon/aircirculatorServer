package service

import (
	"io"
	"log"
	"os"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

func (s *Service) GetFirmware(md5Hash string, sha256Hash string) (firmware []byte, err error) {
	var file *os.File
	file, err = os.Open(s.cfg.FirmwarePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("file close failed: %s", err)
		}
	}(file)

	firmware, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if md5Hash == utils.GetMD5(firmware) {
		return nil, ErrFirmwareNoUpdates
	}
	if sha256Hash == utils.GetSHA256(firmware) {
		return nil, ErrFirmwareNoUpdates
	}

	return firmware, nil
}
