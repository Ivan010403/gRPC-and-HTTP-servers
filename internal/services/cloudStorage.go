package cloudStorage

import (
	"gRPCserver/internal/services/files"
	"log/slog"
)

// Add interface

type Cloud struct {
	log *slog.Logger
}

func NewCloud(logger *slog.Logger) *Cloud {
	return &Cloud{
		log: logger,
	}
}

func (c *Cloud) Write(buff []byte, name, file_type string) error {
	c.log.Info("Write() started", slog.String("name", name), slog.String("type", file_type))

	cf := files.File{Name: name, Filetype: file_type}
	name, err := cf.WriteFile(buff)
	if err != nil {
		c.log.Error("writing file error", slog.Any("err", err))
		return err
	}

	c.log.Info("Write() successfully completed", slog.String("name_new_file", name))
	return nil
}

func (c *Cloud) Delete(name, file_type string) error {
	c.log.Info("Delete() started", slog.String("name", name), slog.String("type", file_type))

	cf := files.File{Name: name, Filetype: file_type}
	err := cf.DeleteFile()
	if err != nil {
		c.log.Error("deleting file error", slog.Any("err", err))
		return err
	}

	c.log.Info("Delete() successfully completed")

	return nil
}

func (c *Cloud) Get(name, file_type string) ([]byte, error) {
	c.log.Info("Get() started", slog.String("name", name), slog.String("type", file_type))

	sf := files.File{Name: name, Filetype: file_type}
	buff, err := sf.ReadFile()
	if err != nil {
		c.log.Error("getting file error", slog.Any("err", err))
		return nil, err
	}

	c.log.Info("Get() successfully completed")
	return buff, nil
}
