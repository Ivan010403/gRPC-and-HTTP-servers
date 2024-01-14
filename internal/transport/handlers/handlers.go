package handlers

import (
	"bytes"
	"fmt"
	"io"

	proto "github.com/Ivan010403/proto/protoc/go"
)

type FileWork interface {
	Write([]byte, string, string) error
}

type StreamHandler struct {
	ChanSave   chan struct{}
	ChanDelete chan struct{}
	ChanCheck  chan struct{}
	proto.UnimplementedCloudServer
	Worker FileWork
}

// TODO: add validation
func (s *StreamHandler) UploadFile(stream proto.Cloud_UploadFileServer) error {
	s.ChanSave <- struct{}{}
	defer func() {
		<-s.ChanSave
	}()

	fmt.Println("Hello from UploadFile")

	file_bytes := bytes.Buffer{}

	r, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("can not get NameFile from req")
	}
	name := r.GetNameFile()

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&proto.UploadFileResponce{NameFile: "hello!"})
			break
		}
		if err != nil {
			return err
		}

		_, err = file_bytes.Write(r.GetFile())
		if err != nil {
			return fmt.Errorf("failed to write in internal buff %w", err)
		}
	}

	fmt.Println(file_bytes.Len(), name)

	err = s.Worker.Write(file_bytes.Bytes(), name, "jpg")

	if err != nil {
		return fmt.Errorf("failed in calling Write() %w", err)
	}

	return nil
}
