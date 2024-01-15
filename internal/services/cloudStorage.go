package service

import (
	"gRPCserver/internal/services/files"
	"gRPCserver/internal/storage/postgres"
	"log/slog"
)

// Add interface
// Close connection to database gracefully

type Cloud struct {
	log  *slog.Logger
	strg *postgres.Storage
}

func NewCloud(logger *slog.Logger, storage *postgres.Storage) *Cloud {
	return &Cloud{
		log:  logger,
		strg: storage,
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

	err = c.strg.SaveFile(name)
	if err != nil {
		c.log.Error("failed to save file in db", slog.Any("err", err))
	}

	c.log.Info("Write() successfully completed", slog.String("name_new_file", name))
	return nil
}

func (c *Cloud) Update(buff []byte, name, file_type string) error {
	c.log.Info("Update() started", slog.String("name", name), slog.String("type", file_type))

	cf := files.File{Name: name, Filetype: file_type}
	name, err := cf.UpdateFile(buff)
	if err != nil {
		c.log.Error("updating file error", slog.Any("err", err))
		return err
	}

	err = c.strg.UpdateFile(name)
	if err != nil {
		c.log.Error("failed to update date in db", slog.Any("err", err))
	}

	c.log.Info("Update() successfully completed", slog.String("name_new_file", name))
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

	err = c.strg.DeleteFile(name + "." + file_type)
	if err != nil {
		c.log.Error("failed to delete file in db", slog.Any("err", err))
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

func (c *Cloud) GetFullData() ([]postgres.File, error) {
	c.log.Info("GetFullData() started")

	data, err := c.strg.GetFullData()
	if err != nil {
		c.log.Error("getting full data error", slog.Any("err", err))
		return nil, err
	}

	c.log.Info("GetFullData() successfully completed")
	return data, nil
}
