package grpcHandlers

import (
	"context"
	"errors"
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
	"google.golang.org/grpc/metadata"
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
	cookieName      = "user_token"
	shortenerTable  = "shortener"
	statisticsTable = "statistics"
)

// MainPage - base handler for short url
func (p *HandlerProvider) MainPage(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	var userID int64
	var err error
	if *p.Config.DataBase != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			p.Sugar.Error("no metadata in context")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		cookieHeaders := md.Get("cookie")
		if len(cookieHeaders) == 0 {
			p.Sugar.Error("no cookie header")
			return nil, status.Error(codes.Unauthenticated, "missing cookie")
		}

		var token string
		cookies := strings.Split(cookieHeaders[0], "; ")
		for _, c := range cookies {
			if strings.HasPrefix(c, cookieName+"=") {
				token = strings.TrimPrefix(c, cookieName+"=")
				break
			}
		}

		if userID, err = p.Session.ParseCookie(token); err != nil {
			p.Sugar.Error("missing cookie")
			return nil, status.Error(codes.Unauthenticated, "missing cookie")
		}
	}

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
			p.Sugar.Error("duplicate error")
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		p.Sugar.Error("failed to store record")
		return nil, status.Error(codes.Internal, err.Error())
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
		p.Sugar.Error("failed to get record")
		return nil, status.Error(codes.Internal, err.Error())
	}

	record, ok := rec.(*pb.ShortenRecord)
	if !ok {
		p.Sugar.Error("invalid record")
		return nil, status.Error(codes.Internal, "invalid record")
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
	var userID int64
	var err error
	if *p.Config.DataBase != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			p.Sugar.Error("no metadata in context")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		cookieHeaders := md.Get("cookie")
		if len(cookieHeaders) == 0 {
			p.Sugar.Error("no cookie header")
			return nil, status.Error(codes.Unauthenticated, "missing cookie")
		}

		var token string
		cookies := strings.Split(cookieHeaders[0], "; ")
		for _, c := range cookies {
			if strings.HasPrefix(c, cookieName+"=") {
				token = strings.TrimPrefix(c, cookieName+"=")
				break
			}
		}

		if userID, err = p.Session.ParseCookie(token); err != nil {
			p.Sugar.Error("missing cookie")
			return nil, status.Error(codes.Unauthenticated, "missing cookie")
		}
	}

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
			p.Sugar.Error("duplicate error")
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		p.Sugar.Error("failed to store record")
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := pb.PingResponse{
		Message: "database is connected",
	}

	return &response, nil
}

// ShortenBatchHandler - handler for short url batches
func (p *HandlerProvider) ShortenBatchHandler(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	var userID int64
	var err error
	if *p.Config.DataBase != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			p.Sugar.Error("no metadata in context")
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		cookieHeaders := md.Get("cookie")
		if len(cookieHeaders) == 0 {
			p.Sugar.Error("no cookie header")
			return nil, status.Error(codes.Unauthenticated, "missing cookie")
		}

		var token string
		cookies := strings.Split(cookieHeaders[0], "; ")
		for _, c := range cookies {
			if strings.HasPrefix(c, cookieName+"=") {
				token = strings.TrimPrefix(c, cookieName+"=")
				break
			}
		}

		if userID, err = p.Session.ParseCookie(token); err != nil {
			p.Sugar.Error("missing cookie")
			return nil, status.Error(codes.Unauthenticated, "missing cookie")
		}
	}

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
		p.Sugar.Error("failed to store records")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ShortenBatchResponse{Results: result}, nil
}

// UserUrlsHandler - handler for get user short urls
func (p HandlerProvider) UserUrlsHandler(ctx context.Context, _ *emptypb.Empty) (*pb.UserUrlsResponse, error) {
	var userID int64
	var err error

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		p.Sugar.Error("no metadata in context")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	cookieHeaders := md.Get("cookie")
	if len(cookieHeaders) == 0 {
		p.Sugar.Error("no cookie header")
		return nil, status.Error(codes.Unauthenticated, "missing cookie")
	}

	var token string
	cookies := strings.Split(cookieHeaders[0], "; ")
	for _, c := range cookies {
		if strings.HasPrefix(c, cookieName+"=") {
			token = strings.TrimPrefix(c, cookieName+"=")
			break
		}
	}

	if userID, err = p.Session.ParseCookie(token); err != nil {
		p.Sugar.Error("missing cookie")
		return nil, status.Error(codes.Unauthenticated, "missing cookie")
	}

	rcs, err := p.Storage.GetAll(shortenerTable, userID)
	if err != nil {
		p.Sugar.Error("failed to get shortener records")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(rcs) == 0 {
		p.Sugar.Error("shortener records not found")
		return nil, status.Error(codes.NotFound, "shortener records not found")
	}

	records := make([]*pb.UserUrlsResult, 0, len(rcs))
	for _, rc := range rcs {
		value, ok := rc.(pb.UserUrlsResult)
		if !ok {
			p.Sugar.Error("invalid shorten record")
			return nil, status.Error(codes.Internal, "invalid shorten record")
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
	var userID int64
	var err error

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		p.Sugar.Error("no metadata in context")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	cookieHeaders := md.Get("cookie")
	if len(cookieHeaders) == 0 {
		p.Sugar.Error("no cookie header")
		return nil, status.Error(codes.Unauthenticated, "missing cookie")
	}

	var token string
	cookies := strings.Split(cookieHeaders[0], "; ")
	for _, c := range cookies {
		if strings.HasPrefix(c, cookieName+"=") {
			token = strings.TrimPrefix(c, cookieName+"=")
			break
		}
	}

	if userID, err = p.Session.ParseCookie(token); err != nil {
		p.Sugar.Error("missing cookie")
		return nil, status.Error(codes.Unauthenticated, "missing cookie")
	}

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
		p.Sugar.Error("permission denied")
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	_, subnet, err := net.ParseCIDR(*p.Config.TrustedSubnet)
	if err != nil {
		p.Sugar.Error("permission denied")
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	if ip == nil || !subnet.Contains(ip) {
		p.Sugar.Error("permission denied")
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	rec, err := p.Storage.GetStat(statisticsTable)
	if err != nil {
		p.Sugar.Error("failed to get statistics")
		return nil, status.Error(codes.Internal, "failed to get statistics")
	}

	record, ok := rec.(*pb.TrustedSubnetResponse)
	if !ok {
		p.Sugar.Error("invalid record")
		return nil, status.Error(codes.Internal, "invalid record")
	}

	return record, nil
}
