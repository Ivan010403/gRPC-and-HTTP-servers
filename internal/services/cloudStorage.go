package cloudStorage

import (
	"fmt"
	"gRPCserver/internal/services/files"
	"log/slog"
)

type FileWriter interface {
	WriteFile([]byte) (string, error)
}

type Cloud struct {
	log *slog.Logger
	FileWriter
}

func NewCloud(logger *slog.Logger) *Cloud {
	return &Cloud{
		log: logger,
	}
}

func (c *Cloud) Write(buff []byte, name, file_type string) error {
	c.log.Info("request started", slog.String("name", name), slog.String("type", file_type))

	cf := files.File{Name: name, Filetype: file_type}
	name, err := cf.WriteFile(buff)
	if err != nil {
		return fmt.Errorf("can not write file: %w", err)
	}

	c.log.Info("write() successfully completed", slog.String("name_new_file", name))
	return nil
}
