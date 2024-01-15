package tests

import (
	"bytes"
	"fmt"
	main_test "gRPCserver/tests/main"
	"io"
	"strconv"
	"testing"

	proto "github.com/Ivan010403/proto/protoc/go"
)

func TestHandlerGetFile(t *testing.T) {
	ctx, cl := main_test.New(t)

	for i := 0; i < 10; i++ {
		stream, err := cl.CloudCl.GetFile(ctx, &proto.GetFileRequest{NameFile: "test" + strconv.Itoa(i), FileFormat: "txt"})
		if err != nil {
			t.Fatalf("failed create stream")
		}

		file_bytes := bytes.Buffer{}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("failed recv()")
			}

			_, err = file_bytes.Write(resp.GetFile())
			if err != nil {
				t.Fatalf("failed write()")
			}
		}
		fmt.Println(len(file_bytes.Bytes()))
	}
}
