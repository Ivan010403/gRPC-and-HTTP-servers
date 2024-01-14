package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	proto "github.com/Ivan010403/proto/protoc/go"
)

type StreamHandler struct {
	ChanSave   chan struct{}
	ChanDelete chan struct{}
	ChanCheck  chan struct{}
	proto.UnimplementedCloudServer
}

func (s *StreamHandler) UploadFile(stream proto.Cloud_UploadFileServer) error {
	s.ChanSave <- struct{}{}
	defer func() {
		<-s.ChanSave
	}()

	file_bytes := bytes.Buffer{}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&proto.UploadFileResponce{NameFile: "hello!"})
			break
		}
		if err != nil {
			break
		}

		file_bytes.Write(r.GetFile())
	}

	f, err := os.Create("test2.jpg")
	if err != nil {
		fmt.Println("error creation", err)
	}

	file_bytes.WriteTo(f)
	return nil
}
