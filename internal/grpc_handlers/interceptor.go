package grpchandlers

import (
	"context"
	"strings"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/session"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// InterceptorProvider - struct that contains the necessary interceptor settings
type InterceptorProvider struct {
	Config  *config.Config
	Sugar   *zap.SugaredLogger
	Session *session.SessionProvider
}

type ctxKey string

const userIDKey ctxKey = "userID"

func (p *InterceptorProvider) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
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
				p.Sugar.Error("not authorized")
				return nil, status.Error(codes.Unauthenticated, "not authorized")
			}
		}

		ctx = context.WithValue(ctx, userIDKey, userID)
		return handler(ctx, req)
	}
}
