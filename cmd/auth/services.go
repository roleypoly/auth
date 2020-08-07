package main

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	pbBackend "github.com/roleypoly/rpc/auth/backend"
	pbClient "github.com/roleypoly/rpc/auth/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthClientService struct {
	*pbClient.UnimplementedAuthClientServer
}

type AuthBackendService struct {
	*pbBackend.UnimplementedAuthBackendServer
}

func (ab *AuthBackendService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	secret, err := grpc_auth.AuthFromMD(ctx, "shared")
	if err != nil {
		return nil, err
	}

	if secret != sharedSecret {
		return nil, grpc.Errorf(codes.Unauthenticated, "Bad authorization string")
	}

	return ctx, nil
}

func (ab *AuthBackendService) GetSessionChallenge(ctx context.Context, req *pbBackend.UserSlug) (*pbBackend.AuthChallenge, error) {
	return &pbBackend.AuthChallenge{
		UserID:     req.UserID,
		MagicUrl:   "https://roleypoly.local/magic/ewqjior",
		MagicWords: "oh my god",
	}, nil
}
