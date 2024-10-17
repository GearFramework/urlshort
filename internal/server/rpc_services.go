package server

import (
	"context"
	"errors"
	"net"

	"github.com/GearFramework/urlshort/internal/app"
	"github.com/GearFramework/urlshort/internal/pkg"
	"github.com/GearFramework/urlshort/internal/pkg/auth"
	"github.com/GearFramework/urlshort/internal/pkg/logger"
	"github.com/GearFramework/urlshort/internal/server/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *RPCServer) getUserIDFromCtx(ctx context.Context) (userID int, err error) {
	token := ctx.Value("Authorization")
	if token == nil {
		return 0, errors.New("user ID is missing")
	}
	_, ok := token.(string)
	if !ok {
		return 0, auth.ErrInvalidAuthorization
	}
	return s.api.Auth(token.(string))
}

// Ping check connection to storage
func (s *RPCServer) Ping(ctx context.Context, _ *proto.PingRequest) (*proto.PingResponse, error) {
	if err := s.api.(*app.ShortApp).Store.Ping(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.PingResponse{}, nil
}

// EncodeURL return short url for requested url
func (s *RPCServer) EncodeURL(ctx context.Context, in *proto.EncodeURLRequest) (*proto.EncodeURLResponse, error) {
	userID, err := s.getUserIDFromCtx(ctx)
	if err != nil {
		logger.Log.Errorf("RPC unauthorized: %v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := proto.EncodeURLResponse{}
	shortURL, _ := s.api.EncodeURL(ctx, userID, in.Url)
	response.Short = shortURL
	return &response, nil
}

// BatchEncodeURLs return short urls for urls in batch json content type request
func (s *RPCServer) BatchEncodeURLs(ctx context.Context, in *proto.BatchEncodeURLRequest) (*proto.BatchEncodeURLResponse, error) {
	userID, err := s.getUserIDFromCtx(ctx)
	if err != nil {
		logger.Log.Errorf("RPC unauthorized: %v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := proto.BatchEncodeURLResponse{}
	batch := []pkg.BatchURLs{}
	for _, v := range in.Urls {
		batch = append(batch, pkg.BatchURLs{
			OriginalURL:   v.OriginalUrl,
			CorrelationID: v.CorrelationId,
		})
	}
	for _, v := range s.api.BatchEncodeURL(ctx, userID, batch) {
		response.Urls = append(response.Urls, &proto.BatchEncodeURLResponse_BatchURL{
			ShortUrl:      v.ShortURL,
			CorrelationId: v.CorrelationID,
		})
	}
	return &response, nil
}

// DecodeURL return url by short code
func (s *RPCServer) DecodeURL(ctx context.Context, in *proto.DecodeURLRequest) (*proto.DecodeURLResponse, error) {
	url, err := s.api.DecodeURL(ctx, in.Code)
	if errors.Is(err, app.ErrShortURLIsDeleted) {
		logger.Log.Errorf("%s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err != nil {
		logger.Log.Errorf("%s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	logger.Log.Infof("Request short code: %s url: %s", in.Code, url)
	return &proto.DecodeURLResponse{OriginalUrl: url}, nil
}

// GetUserURLs handler of request on get all saved urls by user
func (s *RPCServer) GetUserURLs(ctx context.Context, _ *proto.GetUserURLsRequest) (*proto.GetUserURLsResponse, error) {
	userID, err := s.getUserIDFromCtx(ctx)
	if err != nil {
		logger.Log.Errorf("RPC unauthorized: %v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := proto.GetUserURLsResponse{
		Urls: []*proto.GetUserURLsResponse_UserURL{},
	}
	userURLs := s.api.GetUserURLs(ctx, userID)
	if len(userURLs) == 0 {
		return nil, status.Error(codes.Internal, "no content")
	}
	for _, v := range userURLs {
		response.Urls = append(response.Urls, &proto.GetUserURLsResponse_UserURL{
			OriginalUrl: v.URL,
			ShortUrl:    v.ShortURL,
		})
	}
	return &response, nil
}

// DeleteUserURLs remove user saved urls codes by codes from request
func (s *RPCServer) DeleteUserURLs(ctx context.Context, in *proto.DeleteUserURLsRequest) (*proto.DeleteUserURLsResponse, error) {
	userID, err := s.getUserIDFromCtx(ctx)
	if err != nil {
		logger.Log.Errorf("RPC unauthorized: %v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.api.DeleteUserURLs(ctx, userID, in.Urls)
	return &proto.DeleteUserURLsResponse{}, nil
}

// GetInternalStats return internal statistics about short urls and uers
func (s *RPCServer) GetInternalStats(ctx context.Context, in *proto.GetStatsRequest) (*proto.GetStatsResponse, error) {
	if err := validateUserIP(ctx, s.Conf.TrustedSubnet); err != nil {
		logger.Log.Errorf("unauthorized access: %s\n", err)
		return nil, status.Error(codes.PermissionDenied, "Forbidden")
	}
	stats, err := s.api.GetStats(ctx)
	if err != nil {
		logger.Log.Errorf("RPC internal stats error: %v\n", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.GetStatsResponse{Urls: int32(stats.URLs), Users: int32(stats.Users)}, nil
}

func validateUserIP(ctx context.Context, trustedSubnet string) error {
	_, trustNet, err := auth.GetTrustedIP(trustedSubnet)
	if err != nil {
		return err
	}
	userIP, err := getXRealIP(ctx)
	if err != nil {
		return err
	}
	if !trustNet.Contains(userIP) {
		return auth.ErrIPNotFromTrustedNetwork
	}
	return nil
}

func getXRealIP(ctx context.Context) (net.IP, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-Real-IP")
		if len(values) == 0 {
			return nil, auth.ErrIPNotFromTrustedNetwork
		}
		IP := values[0]
		if IP == "" {
			return nil, auth.ErrEmptyXRealIP
		}
		return auth.ParseIP(IP), nil
	}
	return nil, auth.ErrEmptyXRealIP
}
