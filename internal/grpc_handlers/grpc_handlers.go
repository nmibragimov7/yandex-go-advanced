package grpchandlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/models"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/storage/db/shortener"
	"yandex-go-advanced/internal/util"
	pb "yandex-go-advanced/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HandlerProvider - struct that contains the necessary handler settings
type HandlerProvider struct {
	pb.UnimplementedShortenerServiceServer
	Config  *config.Config
	Storage storage.Storage
	Sugar   *zap.SugaredLogger
	Session *session.SessionProvider
}

const (
	cookieName         = "user_token"
	shortenerTable     = "shortener"
	statisticsTable    = "statistics"
	duplicateErrorKey  = "duplicate error"
	saveErrorKey       = "failed to store record"
	getErrorKey        = "failed to get record"
	invalidErrorKey    = "invalid record"
	permissionErrorKey = "permission denied"
)

// MainPage - base handler for short url
func (p *HandlerProvider) MainPage(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	userID := ctx.Value("userID").(int64)
	var err error

	url := in.Url

	key := util.GetKey()
	record := &pb.ShortenRecord{
		ShortUrl:    key,
		OriginalUrl: url,
		UserId:      userID,
		IsDeleted:   false,
	}

	_, err = p.Storage.Set(shortenerTable, record)
	if err != nil {
		var duplicateError *shortener.DuplicateError
		if errors.As(err, &duplicateError) {
			p.Sugar.Error(duplicateErrorKey)
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("duplicate error: %s", err.Error()))
		}

		p.Sugar.Error(saveErrorKey)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to store record: %s", err.Error()))
	}

	response := pb.ShortenResponse{
		Url: *p.Config.BaseURL + "/" + key,
	}

	return &response, nil
}

// GetItem - handler for get url by id
func (p *HandlerProvider) GetItem(_ context.Context, in *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	rec, err := p.Storage.Get(shortenerTable, in.GetId())
	if err != nil {
		p.Sugar.Error(getErrorKey)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get record: %s", err.Error()))
	}

	record, ok := rec.(*pb.ShortenRecord)
	if !ok {
		p.Sugar.Error(invalidErrorKey)
		return nil, status.Error(codes.Internal, invalidErrorKey)
	}

	if record.IsDeleted {
		p.Sugar.Error("record is deleted")
		return nil, status.Error(codes.NotFound, "record is deleted")
	}

	if record.OriginalUrl == "" {
		p.Sugar.Error("record not found")
		return nil, status.Error(codes.NotFound, "record not found")
	}

	return &pb.GetItemResponse{Url: record.OriginalUrl}, nil
}

// ShortenHandler - handler for short url by json
func (p *HandlerProvider) ShortenHandler(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	userID := ctx.Value("userID").(int64)
	var err error

	url := in.Url

	key := util.GetKey()
	record := &pb.ShortenRecord{
		ShortUrl:    key,
		OriginalUrl: url,
		UserId:      userID,
		IsDeleted:   false,
	}

	_, err = p.Storage.Set(shortenerTable, record)
	if err != nil {
		var duplicateError *shortener.DuplicateError
		if errors.As(err, &duplicateError) {
			p.Sugar.Error(duplicateErrorKey)
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("duplicate error: %s", err.Error()))
		}

		p.Sugar.Error(saveErrorKey)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to store record: %s", err.Error()))
	}

	response := pb.ShortenResponse{
		Url: *p.Config.BaseURL + "/" + key,
	}

	return &response, nil
}

// PingHandler - handler for ping storage
func (p *HandlerProvider) PingHandler(ctx context.Context, _ *emptypb.Empty) (*pb.PingResponse, error) {
	if p.Config.DataBase == nil {
		p.Sugar.Error("data base not configured")
		return nil, status.Error(codes.FailedPrecondition, "data base not configured")
	}

	err := p.Storage.Ping(ctx)
	if err != nil {
		p.Sugar.Error("failed to ping database")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to ping database: %s", err.Error()))
	}

	response := pb.PingResponse{
		Message: "database is connected",
	}

	return &response, nil
}

// ShortenBatchHandler - handler for short url batches
func (p *HandlerProvider) ShortenBatchHandler(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	userID := ctx.Value("userID").(int64)
	var err error

	values := make([]interface{}, 0, len(in.Items))
	result := make([]*pb.ShortenBatchResult, 0, len(in.Items))
	for _, value := range in.Items {
		key := util.GetKey()
		values = append(values, &pb.ShortenRecord{
			OriginalUrl: value.OriginalUrl,
			ShortUrl:    key,
			UserId:      userID,
		})
		result = append(result, &pb.ShortenBatchResult{
			CorrelationId: value.CorrelationId,
			ShortUrl:      *p.Config.BaseURL + "/" + key,
		})
	}

	err = p.Storage.SetAll(shortenerTable, values)
	if err != nil {
		p.Sugar.Error(saveErrorKey)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to store records: %s", err.Error()))
	}

	return &pb.ShortenBatchResponse{Results: result}, nil
}

// UserUrlsHandler - handler for get user short urls
func (p HandlerProvider) UserUrlsHandler(ctx context.Context, _ *emptypb.Empty) (*pb.UserUrlsResponse, error) {
	userID := ctx.Value("userID").(int64)
	var err error

	rcs, err := p.Storage.GetAll(shortenerTable, userID)
	if err != nil {
		p.Sugar.Error(getErrorKey)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get records: %s", err.Error()))
	}

	if len(rcs) == 0 {
		p.Sugar.Error("records not found")
		return nil, status.Error(codes.NotFound, "records not found")
	}

	records := make([]*pb.UserUrlsResult, 0, len(rcs))
	for _, rc := range rcs {
		value, ok := rc.(*pb.UserUrlsResult)
		if !ok {
			p.Sugar.Error(invalidErrorKey)
			return nil, status.Error(codes.Internal, invalidErrorKey)
		}
		records = append(records, &pb.UserUrlsResult{
			ShortUrl:    *p.Config.BaseURL + "/" + value.ShortUrl,
			OriginalUrl: value.OriginalUrl,
		})
	}

	return &pb.UserUrlsResponse{Results: records}, nil
}

// UserUrlsDeleteHandler - handler for remove user short urls
func (p HandlerProvider) UserUrlsDeleteHandler(ctx context.Context, in *pb.UserUrlsDeleteRequest) (*pb.UserUrlsDeleteResponse, error) {
	userID := ctx.Value("userID").(int64)

	generate := func(userID int64, key string) chan interface{} {
		out := make(chan interface{}, 1)
		go func() {
			defer close(out)
			val := &models.ShortenBatchUpdateRequest{
				ShortURL: key,
				UserID:   userID,
			}
			out <- val
		}()

		return out
	}

	values := make([]chan interface{}, 0, len(in.Items))
	for _, value := range in.Items {
		values = append(values, generate(userID, value.ShortUrl))
	}

	go func() {
		done := make(chan struct{})
		defer close(done)
		p.Storage.AddToChannel(shortenerTable, done, values...)
	}()

	return &pb.UserUrlsDeleteResponse{Message: "request is accepted"}, nil
}

// TrustedSubnetHandler - handler for get all shorten urls, users by trusted subnet
func (p *HandlerProvider) TrustedSubnetHandler(_ context.Context, in *pb.TrustedSubnetRequest) (*pb.TrustedSubnetResponse, error) {
	xRealIP := in.XRealIp
	ip := net.ParseIP(strings.TrimSpace(xRealIP))

	if *p.Config.TrustedSubnet == "" {
		p.Sugar.Error(permissionErrorKey)
		return nil, status.Error(codes.PermissionDenied, permissionErrorKey)
	}

	_, subnet, err := net.ParseCIDR(*p.Config.TrustedSubnet)
	if err != nil {
		p.Sugar.Error(permissionErrorKey)
		return nil, status.Error(codes.PermissionDenied, permissionErrorKey)
	}

	if ip == nil || !subnet.Contains(ip) {
		p.Sugar.Error(permissionErrorKey)
		return nil, status.Error(codes.PermissionDenied, permissionErrorKey)
	}

	rec, err := p.Storage.GetStat(statisticsTable)
	if err != nil {
		p.Sugar.Error("failed to get statistics")
		return nil, status.Error(codes.Internal, "failed to get statistics")
	}

	record, ok := rec.(*pb.TrustedSubnetResponse)
	if !ok {
		p.Sugar.Error(invalidErrorKey)
		return nil, status.Error(codes.Internal, invalidErrorKey)
	}

	return record, nil
}
