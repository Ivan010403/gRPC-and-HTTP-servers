package handlers

import (
	"bytes"
	"context"
	"fmt"
	"gRPCserver/internal/storage/postgres"
	"io"
	"os"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileWork interface {
	Write([]byte, string, string) error
	Update([]byte, string, string) error
	Delete(string, string) error
	Get(string, string) ([]byte, error)
	GetFullData() ([]postgres.File, error)
}

type CloudServer struct {
	ChanSave   chan struct{}
	ChanDelete chan struct{}
	ChanCheck  chan struct{}
	proto.UnimplementedCloudServer
	Worker FileWork
}

// TODO: add validation ON ALL FILE
func (s *CloudServer) UploadFile(stream proto.Cloud_UploadFileServer) error {
	s.ChanSave <- struct{}{}
	defer func() {
		<-s.ChanSave
	}()

	file_bytes := bytes.Buffer{}

	r, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("can not get NameFile and FileFormat from req")
	}
	name := r.GetNameFile()
	format := r.GetFileFormat()

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&proto.UploadFileResponce{FullName: r.NameFile + r.FileFormat})
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

	if _, err := os.Stat(name + "." + format); os.IsNotExist(err) {
		err = s.Worker.Write(file_bytes.Bytes(), name, format)
		if err != nil {
			return fmt.Errorf("failed in calling Write() %w", err)
		}
	} else {
		err = s.Worker.Update(file_bytes.Bytes(), name, format)
		if err != nil {
			return fmt.Errorf("failed in calling Update() %w", err)
		}
	}
	return nil
}

func (s *CloudServer) DeleteFile(ctx context.Context, request *proto.DeleteFileRequest) (*proto.DeleteFileResponce, error) {
	s.ChanDelete <- struct{}{}
	defer func() {
		<-s.ChanDelete
	}()
	//TODO: change that unconvenient
	full_name := request.GetNameFile() + "." + request.GetFileFormat()
	err := s.Worker.Delete(request.GetNameFile(), request.GetFileFormat())

	if err != nil {
		return nil, status.Error(codes.Internal, "can't delete file")
	}

	return &proto.DeleteFileResponce{FullName: full_name}, nil

}

// TODO: new chan and new operation
func (s *CloudServer) GetFile(request *proto.GetFileRequest, stream proto.Cloud_GetFileServer) error {
	s.ChanDelete <- struct{}{}
	defer func() {
		<-s.ChanDelete
	}()

	data, err := s.Worker.Get(request.GetNameFile(), request.GetFileFormat())
	if err != nil {
		return status.Error(codes.Internal, "can't get file")
	}

	buff := &proto.GetFileResponce{}

	//TODO: change 1024 to normal value
	for currentByte := 0; currentByte < len(data); currentByte += 1024 {
		if currentByte+1024 > len(data) {
			buff.File = data[currentByte:len(data)]
		} else {
			buff.File = data[currentByte : currentByte+1024]
		}
		if err := stream.Send(buff); err != nil {
			return err
		}
	}

	return nil
}

func (s *CloudServer) GetFullData(*proto.GetFullDataRequest, proto.Cloud_GetFullDataServer) error {
	panic("implement me")
}
