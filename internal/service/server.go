// Package server - содержит  себе все методы и структуры работы сервера.
package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"syscall"
	pb "tages/proto"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// File - структура загружаемого/выгружаемого файла.
type File struct {
	name   string
	buffer *bytes.Buffer
}

type server struct {
	pb.TagesServer
}

// sLog - серверный логгер, пишет логи в папку logs.
var sLog zerolog.Logger

// RunServer - запуск сервера приложения.
func RunServer() {
	file, err := os.OpenFile("./logs/"+logName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("erorrlog file opening error:%s\n", err.Error())
	}
	defer file.Close()
	sLog = zerolog.New(file).With().Timestamp().Logger()
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		sLog.Fatal().Err(err).Msg("error listening")
	}
	defer l.Close()
	grpcS := grpc.NewServer()
	pb.RegisterTagesServer(grpcS, &server{})
	sLog.Info().Msg("server start at localhost:8080")
	if err = grpcS.Serve(l); err != nil {
		sLog.Fatal().Err(err).Msg("error start server")
	}
}

// Upload - метод загрузки файла на сервер со стороный клиента.
func (t *server) Upload(stream pb.Tages_UploadServer) error {
	sLog.Info().Msg("upload file started")
	file := &File{
		buffer: &bytes.Buffer{},
		name:   "",
	}
	wd, err := newPath()
	if err != nil {
		sLog.Error().Err(err).Msg(err.Error())
	}
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				err = ioutil.WriteFile(wd+"/pkg/server_cache/"+file.name, file.buffer.Bytes(), 0644)
				if err != nil {
					sLog.Error().Err(err).Msg("create file error")
					return err
				}
				break
			}
			sLog.Error().Err(err).Msg(err.Error())
			return err
		}
		file.name = response.GetStr()
		_, err = file.buffer.Write(response.Data)
		if err != nil {
			sLog.Error().Err(err).Msg(err.Error())
			return err
		}
	}
	err = stream.SendAndClose(&pb.StringForm{Str: "Upload to server is done."})
	if err != nil {
		sLog.Error().Err(err).Msg("send response error")
		return err
	}
	sLog.Info().Msg("upload file ended")
	return nil
}

// Download - метод выгрузки файла с сервера на сторону клиента.
func (t *server) Download(request *pb.StringForm, conn pb.Tages_DownloadServer) error {
	sLog.Info().Msg("download file started")
	wd, err := newPath()
	if err != nil {
		sLog.Error().Err(err).Msg(err.Error())
	}
	file, err := os.Open(wd + request.GetStr())
	if err != nil {
		sLog.Error().Err(err).Msg(err.Error())
		return errors.New("no such file in directory, try agan")
	}
	defer file.Close()
	buffer := make([]byte, 100000)
	for {
		counter, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				sLog.Error().Err(err).Msg(err.Error())
				return err
			}
			break
		}
		response := &pb.DownloadResponse{
			Data: buffer[:counter],
		}
		err = conn.Send(response)
		if err != nil {
			sLog.Fatal().Err(err).Msg("send response error:")
			return err
		}
	}
	sLog.Info().Msg("download file ended")
	return nil
}

// List - метод чтения хранилища фалов сервера и передачи информации о хранимых файлах на сервере клиенту.
func (t *server) List(_ *pb.ListRequest, conn pb.Tages_ListServer) error {
	sLog.Info().Msg("downloaded files in list send started")
	wd, err := newPath()
	if err != nil {
		sLog.Error().Err(err).Msg(err.Error())
	}
	data, err := os.ReadDir(wd + "/pkg/server_cache")
	if err != nil {
		sLog.Fatal().Err(err).Msg("reading dir error")
	}
	for _, file := range data {
		if !file.IsDir() {
			f, err := os.Open(wd + "/pkg/server_cache/" + file.Name())
			if err != nil {
				sLog.Error().Err(err).Msg(err.Error())
				return err
			}
			stat, err := f.Stat()
			if err != nil {
				sLog.Error().Err(err).Msg(err.Error())
				return err
			}
			sysstat := stat.Sys().(*syscall.Stat_t)
			addTime := time.Unix(sysstat.Atim.Sec, 0)
			modTime := time.Unix(sysstat.Mtim.Sec, 0)
			res := fmt.Sprintf("name:%s |creation time:%v |modification time:%v \n", file.Name(), addTime.Format("2 January 2006 15:04"),
				modTime.Format("2 January 2006 15:04"),
			)
			result := &pb.StringForm{
				Str: res,
			}
			err = conn.Send(result)
			if err != nil {
				sLog.Error().Err(err).Msg("send response error")
				return err
			}
		}
	}
	sLog.Info().Msg("send list of downloaded files ended")
	return nil
}

func logName() string {
	day, month, year := time.Now().Date()
	name := fmt.Sprintf("%v-%v-%v.txt", day, month, year)
	return name
}
