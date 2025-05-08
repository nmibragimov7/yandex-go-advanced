package grpc_handlers

import (
	"context"
	"errors"
	"strings"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
	"yandex-go-advanced/internal/storage/db/shortener"
	"yandex-go-advanced/internal/util"
	pb "yandex-go-advanced/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ShortenerHandler struct {
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

func (p *ShortenerHandler) MainPage(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
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

func (p *ShortenerHandler) GetItem(_ context.Context, in *pb.GetItemRequest) (*pb.GetItemResponse, error) {
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
