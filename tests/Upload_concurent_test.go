package tests

import (
	"bytes"
	"crypto/rand"
	main_test "gRPCserver/tests/main"
	"io"
	"strconv"
	"sync"
	"testing"
	"time"

	proto "github.com/Ivan010403/proto/protoc/go"
)

func TestHandlerUpload_Concurent(t *testing.T) {
	ctx, cl := main_test.New(t)

	chunkSize := 1024

	wg := sync.WaitGroup{}

	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(i int) {
			defer wg.Done()
			t.Log("starting", i)

			data := makeBuff_(i)

			stream, err := cl.CloudCl.UploadFile(ctx)
			if err != nil {
				t.Fatalf("stream failed")
			}

			reader := bytes.NewReader(data)
			buffer := make([]byte, chunkSize)

			req := &proto.UploadFileRequest{NameFile: "test" + strconv.Itoa(i), FileFormat: "txt"}

			err = stream.Send(req)
			if err != nil {
				t.Fatalf("send failed")
			}

			for {
				n, err := reader.Read(buffer)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("read failed")
				}

				req := &proto.UploadFileRequest{File: buffer[:n]}

				err = stream.Send(req)

				if err != nil {
					t.Fatalf("send failed")
				}
			}

			_, err = stream.CloseAndRecv()
			if err != nil {
				t.Fatalf("close failed")
			}
			t.Log("done", i, time.Now())
		}(i)
	}

	wg.Wait()
}

func makeBuff_(i int) []byte {
	buff := make([]byte, 5000000)

	rand.Read(buff)
	return buff
}
