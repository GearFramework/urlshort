package server

import (
	"context"
	"errors"
	"net"

	"github.com/GearFramework/urlshort/internal/config"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/auth"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	keyRPCAuthorization = "Authorization"
)

type AuthInterceptor struct {
	api pkg.APIShortener
}

func NewAuthInterceptor(api pkg.APIShortener) *AuthInterceptor {
	return &AuthInterceptor{api}
}

func (i *AuthInterceptor) Auth(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(keyRPCAuthorization)
		if len(values) > 0 {
			token = values[0]
			_, err = i.api.Auth(token)
			if err != nil && errors.Is(err, auth.ErrInvalidAuthorization) {
				return nil, status.Error(codes.Unauthenticated, "wrong user id format")
			}
			if err != nil && errors.Is(err, auth.ErrNeedAuthorization) {
				_, token, err = i.authNewUser(ctx)
				if err != nil {
					return nil, status.Error(codes.Unauthenticated, "wrong user id")
				}
			}
		} else {
			logger.Log.Infoln("empty token in grpc-context; need new token")
			if _, token, err = i.authNewUser(ctx); err != nil {
				return nil, status.Error(codes.Unauthenticated, "wrong user id format")
			}
		}
		ctx = context.WithValue(ctx, keyRPCAuthorization, token)
	}
	return handler(ctx, req)
}

func (i *AuthInterceptor) authNewUser(ctx context.Context) (int, string, error) {
	userID, token, err := i.api.CreateToken()
	if err != nil {
		logger.Log.Error(err.Error())
		return 0, "", err
	}
	logger.Log.Infof("Created user ID: %d", userID)
	return userID, token, nil

}

// RPCServer struct of rpc-server
type RPCServer struct {
	proto.UnimplementedShortlyServiceServer
	Conf *config.ServiceConfig
	rpc  *grpc.Server
	api  pkg.APIShortener
}

// NewRPCServer return new rpc-server
func NewRPCServer(c *config.ServiceConfig, api pkg.APIShortener) (*RPCServer, error) {
	return &RPCServer{
		Conf: c,
		api:  api,
	}, nil
}

// Up start rpc-server
func (s *RPCServer) Up() error {
	listen, err := net.Listen("tcp", s.Conf.GRPCAddress)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	i := NewAuthInterceptor(s.api)
	s.rpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(i.Auth),
	)
	proto.RegisterShortlyServiceServer(s.rpc, s)
	logger.Log.Infof("Start gRPC server at the %s\n", s.Conf.GRPCAddress)
	if err := s.rpc.Serve(listen); err != nil {
		logger.Log.Infof("Failed to Listen and Serve gRPC: %v\n", err)
		return err
	}
	return nil
}
