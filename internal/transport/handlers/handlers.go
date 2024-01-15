package handlers

import (
	"bytes"
	"context"
	"gRPCserver/internal/storage/postgres"
	"io"
	"os"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxChunkSize = 1024

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

func (s *CloudServer) UploadFile(stream proto.Cloud_UploadFileServer) error {
	s.ChanSave <- struct{}{}
	defer func() {
		<-s.ChanSave
	}()

	file_bytes := bytes.Buffer{}

	r, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, "reading namefile and formatfile from stream error")
	}
	name := r.GetNameFile()
	format := r.GetFileFormat()

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&proto.UploadFileResponce{FullName: name + "." + format})
			break
		}
		if err != nil {
			return status.Error(codes.Internal, "reading from stream error")
		}

		_, err = file_bytes.Write(r.GetFile())
		if err != nil {
			return status.Error(codes.Internal, "writing into internal buff error")
		}
	}

	if _, err := os.Stat("../../storage/" + name + "." + format); os.IsNotExist(err) {
		err = s.Worker.Write(file_bytes.Bytes(), name, format)
		if err != nil {
			return status.Error(codes.Internal, "writing to file error")
		}
	} else {
		err = s.Worker.Update(file_bytes.Bytes(), name, format)
		if err != nil {
			return status.Error(codes.Internal, "updating file error")
		}
	}
	return nil
}

func (s *CloudServer) DeleteFile(ctx context.Context, request *proto.DeleteFileRequest) (*proto.DeleteFileResponce, error) {
	s.ChanDelete <- struct{}{}
	defer func() {
		<-s.ChanDelete
	}()

	name := request.GetNameFile()
	format := request.GetFileFormat()

	err := s.Worker.Delete(name, format)

	if err != nil {
		return nil, status.Error(codes.Internal, "deleting file error")
	}

	return &proto.DeleteFileResponce{FullName: name + format}, nil

}

// TODO: is it new chan and new operation???
func (s *CloudServer) GetFile(request *proto.GetFileRequest, stream proto.Cloud_GetFileServer) error {
	s.ChanDelete <- struct{}{}
	defer func() {
		<-s.ChanDelete
	}()

	data, err := s.Worker.Get(request.GetNameFile(), request.GetFileFormat())
	if err != nil {
		return status.Error(codes.Internal, "getting file error")
	}

	buff := &proto.GetFileResponce{}

	for currentByte := 0; currentByte < len(data); currentByte += maxChunkSize {
		if currentByte+maxChunkSize > len(data) {
			buff.File = data[currentByte:]
		} else {
			buff.File = data[currentByte : currentByte+maxChunkSize]
		}
		if err := stream.Send(buff); err != nil {
			return status.Error(codes.Internal, "sending file error")
		}
	}

	return nil
}

func (s *CloudServer) GetFullData(request *proto.GetFullDataRequest, stream proto.Cloud_GetFullDataServer) error {
	s.ChanCheck <- struct{}{}
	defer func() {
		<-s.ChanCheck
	}()

	data, err := s.Worker.GetFullData()
	if err != nil {
		return status.Error(codes.Internal, "geting full data error")
	}

	if err := stream.Send(&proto.GetFullDataResponce{Size: int64(len(data))}); err != nil {
		return status.Error(codes.Internal, "sending size of full data error")
	}

	for i := 0; i < len(data); i++ {
		if err := stream.Send(&proto.GetFullDataResponce{Name: data[i].Name, CreationDate: data[i].Creation_date, UpdatingDate: data[i].Update_date}); err != nil {
			return status.Error(codes.Internal, "sending data error")
		}
	}
	return nil
}
