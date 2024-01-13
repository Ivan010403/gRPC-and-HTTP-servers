package handlers

import (
	"fmt"
	"io"

	proto "github.com/Ivan010403/proto/protoc/go"
)

type StreamHandler struct {
	proto.UnimplementedCloudServer
}

func (s *StreamHandler) UploadFile(stream proto.Cloud_UploadFileServer) error {
	for {
		file, err := stream.Recv()
		if err == io.EOF {
			//TODO: error when request was cancelled before executing
			return stream.SendAndClose(&proto.UploadFileResponce{NameFile: file.NameFile})
		}
		if err != nil {
			return err
		}
		fmt.Println(file.File)
	}
}
