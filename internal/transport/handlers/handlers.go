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
// TODO: make normal logging!!
func (s *CloudServer) UploadFile(stream proto.Cloud_UploadFileServer) error {
	s.ChanSave <- struct{}{}
	defer func() {
		<-s.ChanSave
	}()

	file_bytes := bytes.Buffer{}

	r, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, "can't read namefile and formatfile from stream")
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
			return status.Error(codes.Internal, "can't read from stream")
		}

		_, err = file_bytes.Write(r.GetFile())
		if err != nil {
			return status.Error(codes.Internal, "can't write into internal buff")
		}
	}

	if _, err := os.Stat(name + "." + format); os.IsNotExist(err) {
		err = s.Worker.Write(file_bytes.Bytes(), name, format)
		if err != nil {
			return status.Error(codes.Internal, "can't write file")
		}
	} else {
		err = s.Worker.Update(file_bytes.Bytes(), name, format)
		if err != nil {
			return status.Error(codes.Internal, "can't update file")
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
	name := request.GetNameFile()
	format := request.GetFileFormat()

	err := s.Worker.Delete(name, format)

	if err != nil {
		return nil, status.Error(codes.Internal, "can't delete file")
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
			return status.Error(codes.Internal, "can't send file")
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
		return status.Error(codes.Internal, "can't get full data")
	}

	if err := stream.Send(&proto.GetFullDataResponce{Size: int64(len(data))}); err != nil {
		return status.Error(codes.Internal, "can't send size of full data")
	}

	for i := 0; i < len(data); i++ {
		fmt.Println(data[i])
		if err := stream.Send(&proto.GetFullDataResponce{Name: data[i].Name, CreationDate: data[i].Creation_date, UpdatingDate: data[i].Update_date}); err != nil {
			return err
		}
	}

	return nil
}
