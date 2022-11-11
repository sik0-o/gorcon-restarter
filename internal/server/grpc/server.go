package grpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/sik0-o/gorcon-restarter/v2/internal/config"
	pb "github.com/sik0-o/gorcon-restarter/v2/internal/proto/grpcrcon"
)

type GRPCServer struct {
	pb.GRPCRCONServiceServer
	config *config.Config
	logger *zap.Logger
}

func RunnServe(logger *zap.Logger, config *config.Config, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("unable listen", zap.Error(err))
	}

	logger.Info("listening at address", zap.String("addr", addr))

	s := grpc.NewServer()

	pb.RegisterGRPCRCONServiceServer(s, &GRPCServer{
		config: config,
		logger: logger,
	})
	if err = s.Serve(lis); err != nil {
		logger.Fatal("not able to register rpc", zap.Error(err))
	}
}
