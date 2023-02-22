package service

import (
	"context"
	"github.com/stretchr/testify/require"
	rl "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"os"
	"strings"
	pb "tages/proto"
	"testing"
)

func TestDownload(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcS := grpc.NewServer()
	pb.RegisterTagesServer(grpcS, &server{})
	go grpcS.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(rl.StreamClientInterceptor(rl.NewLimiter(10))))
	require.NoError(t, err)
	client := pb.NewTagesClient(conn)
	downloadConn, err := client.Download(context.Background(), &pb.StringForm{
		Str: "/pkg/server_cache/test.jpg",
	})
	require.NoError(t, err)
	var file *os.File
	wd := testPath(t)
	for {
		connData, err := downloadConn.Recv()
		if err != nil {
			if err == io.EOF {
				file.Close()
				break
			}
			require.NoError(t, err)
		}
		file, err = os.OpenFile(wd+"/pkg/client_cache/test.jpg", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		require.NoError(t, err)
		counter, err := file.Write(connData.Data)
		require.NoError(t, err)
		require.NotNil(t, counter)
	}
}

func TestList(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcS := grpc.NewServer()
	pb.RegisterTagesServer(grpcS, &server{})
	go grpcS.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(rl.StreamClientInterceptor(rl.NewLimiter(100))))
	require.NoError(t, err)
	client := pb.NewTagesClient(conn)
	connData, err := client.List(context.Background(), &pb.ListRequest{})
	require.NoError(t, err)
	for {
		data, err := connData.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}
		require.NotNil(t, data.GetStr())
	}

}

func TestUpload(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer l.Close()
	address := l.Addr().String()
	grpcS := grpc.NewServer()
	pb.RegisterTagesServer(grpcS, &server{})
	go grpcS.Serve(l)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(rl.StreamClientInterceptor(rl.NewLimiter(10))))
	require.NoError(t, err)
	client := pb.NewTagesClient(conn)
	wd := testPath(t)
	file, err := os.Open(wd + "/pkg/client_cache/test.jpg")
	require.NoError(t, err)
	defer file.Close()
	uploadConn, err := client.Upload(context.Background())
	require.NoError(t, err)
	buffer := make([]byte, 100000)
	for {
		counter, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		resp := &pb.UploadRequest{
			Str:  "test.jpg",
			Data: buffer[:counter],
		}
		err = uploadConn.Send(resp)
		require.NoError(t, err)
	}
	response, err := uploadConn.CloseAndRecv()
	require.NoError(t, err)
	require.Equal(t, response.GetStr(), "Upload to server is done.")
}

func testPath(t *testing.T) string {
	wd, err := os.Getwd()
	require.NoError(t, err)
	wdSlice := strings.Split(wd, "/")
	wdSlice = wdSlice[:len(wdSlice)-2]
	result := strings.Join(wdSlice, "/")
	return result
}
