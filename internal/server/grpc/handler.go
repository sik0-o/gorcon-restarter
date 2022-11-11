package grpc

import (
	"context"
	"errors"
	"log"

	be "github.com/sik0-o/gorcon-restarter/v2/internal/battleye"
	"github.com/sik0-o/gorcon-restarter/v2/internal/config"
	pb "github.com/sik0-o/gorcon-restarter/v2/internal/proto/grpcrcon"
	"github.com/sik0-o/gorcon-restarter/v2/internal/service"
	"github.com/sik0-o/gorcon-restarter/v2/internal/service/action"
	"go.uber.org/zap"
)

func (s *GRPCServer) statusResponseError(err error, serverName string) *pb.StatusResponse {
	return &pb.StatusResponse{
		Status:     pb.ResponseStatus_ERROR,
		Message:    err.Error(),
		ServerName: serverName,
	}
}

func (s *GRPCServer) findServer(serverName string) (*config.ServerConf, error) {
	serverConf, ok := s.config.Servers[serverName]
	if !ok {
		return nil, errors.New("server not found")
	}

	return &serverConf, nil
}

func (s *GRPCServer) Lock(ctx context.Context, in *pb.LockRequest) (*pb.StatusResponse, error) {
	serverName := in.ServerName
	serverConf, err := s.findServer(serverName)
	if err != nil {
		return s.statusResponseError(err, serverName), nil
	}

	beService, err := be.NewService(s.logger, serverName, *serverConf)
	if err != nil {
		s.logger.Error("beService error", zap.Error(err))
		return &pb.StatusResponse{
			Status:     pb.ResponseStatus_ERROR,
			Message:    "beService error",
			ServerName: serverName,
		}, nil
	}

	commandService := service.New(beService)

	if err := commandService.Exec(action.NewLock()); err != nil {
		return &pb.StatusResponse{
			Status:     pb.ResponseStatus_ERROR,
			Message:    "commandService exec error",
			ServerName: serverName,
		}, nil
	}

	return &pb.StatusResponse{
		Status:     pb.ResponseStatus_SUCCESS,
		Message:    "success",
		ServerName: serverName,
	}, nil
}

func (s *GRPCServer) Announce(ctx context.Context, in *pb.AnnounceRequest) (*pb.StatusResponse, error) {
	serverName := in.ServerName
	serverConf, ok := s.config.Servers[serverName]
	if !ok {
		return &pb.StatusResponse{
			Status:     pb.ResponseStatus_ERROR,
			Message:    "server not found",
			ServerName: serverName,
		}, nil
	}

	beService, err := be.NewService(s.logger, serverName, serverConf)
	if err != nil {
		s.logger.Error("beService error", zap.Error(err))
		return &pb.StatusResponse{
			Status:     pb.ResponseStatus_ERROR,
			Message:    "beService error",
			ServerName: serverName,
		}, nil
	}

	commandService := service.New(beService)

	if err := commandService.Exec(action.NewAnnounce(in.Announce)); err != nil {
		return &pb.StatusResponse{
			Status:     pb.ResponseStatus_ERROR,
			Message:    "commandService exec error",
			ServerName: serverName,
		}, nil
	}

	return &pb.StatusResponse{
		Status:     pb.ResponseStatus_SUCCESS,
		Message:    "success",
		ServerName: serverName,
	}, nil
}

func (s *GRPCServer) Servers(ctx context.Context, in *pb.ServersRequest) (*pb.ServersResponse, error) {
	log.Printf("Request", in)

	servers := make([]*pb.ServersResponse_ServerConfig, 0)
	for serverName, serverConf := range s.config.Servers {
		servers = append(servers, &pb.ServersResponse_ServerConfig{
			Name:     serverName,
			Host:     serverConf.Host,
			Port:     int32(serverConf.Port),
			Password: serverConf.Password,
			RestartConf: &pb.ServersResponse_RestartConf{
				Period:     int64(serverConf.Restart.Period),
				ServerLock: int64(serverConf.Restart.ServerLock),
				Announcements: &pb.ServersResponse_Announcements{
					At:  serverConf.Restart.Announcements.At,
					Min: serverConf.Restart.Announcements.Min,
					Sec: serverConf.Restart.Announcements.Sec,
				},
			},
		})
	}

	return &pb.ServersResponse{
		Greeting: "hello " + in.Name,
		Servers:  servers,
	}, nil
}
