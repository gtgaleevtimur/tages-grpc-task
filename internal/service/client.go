// Package client - содержит в себе все методы и структуры работы клиента сервиса.
package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	pb "tages/proto"

	rl "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RunClient - запуск клиентского интерфейса приложения.
func RunClient() {
	fmt.Println("Welcome to our service!")
	fmt.Println("To interact with the service, enter one of the suggested commands:")
	fmt.Println("Enter >>Upload<< to upload file to server.")
	fmt.Println("Enter >>Download<< to download file from server.")
	fmt.Println("Enter >>List<< to watch uploaded file at server.")
	fmt.Println("Enter >>Exit<< to end the program.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Command->")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		command = strings.TrimSpace(command)
		switch command {
		case "Upload":
			Upload()
		case "Download":
			Download()
		case "List":
			List()
		case "Exit":
			fmt.Println("Exit from programm.")
			os.Exit(0)
		default:
			fmt.Println("Command not allowed.")
		}
	}
}

// List - функция формирующая запрос серверу о предоставлении информации о хранимых файлах на сервере.
func List() {
	fmt.Println("Getting the list of downloaded files started.")
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(rl.StreamClientInterceptor(rl.NewLimiter(100))))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client := pb.NewTagesClient(conn)
	connData, err := client.List(context.Background(), &pb.ListRequest{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		data, err := connData.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return
		}
		fmt.Print(data.GetStr())
	}
}

// Upload - формирует запрос о загрузке файла на сервер.
func Upload() {
	fmt.Print("File name->")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	name = strings.TrimSpace(name)
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(rl.StreamClientInterceptor(rl.NewLimiter(10))))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client := pb.NewTagesClient(conn)
	uploadConn, err := client.Upload(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	wd, err := newPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	file, err := os.Open(wd + "/pkg/client_cache/" + name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	buffer := make([]byte, 100000)
	for {
		counter, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return
		}
		resp := &pb.UploadRequest{
			Str:  name,
			Data: buffer[:counter],
		}
		err = uploadConn.Send(resp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	response, err := uploadConn.CloseAndRecv()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(response.GetStr())
}

// Download - формирует запрос о скачивании файла с сервера.
func Download() {
	fmt.Print("File name->")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	name = strings.TrimSpace(name)
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(rl.StreamClientInterceptor(rl.NewLimiter(10))))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client := pb.NewTagesClient(conn)
	downloadConn, err := client.Download(context.Background(), &pb.StringForm{
		Str: "/pkg/server_cache/" + name,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var file *os.File
	wd, err := newPath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		connData, err := downloadConn.Recv()
		if err != nil {
			if err == io.EOF {
				file.Close()
				fmt.Println("Download file from server is done")
				break
			}
			fmt.Println(err.Error())
			return
		}

		file, err = os.OpenFile(wd+"/pkg/client_cache/"+name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		_, err = file.Write(connData.Data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func newPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for !strings.HasSuffix(wd, "tages-task") {
		wd = filepath.Dir(wd)
	}
	return wd, nil
}
