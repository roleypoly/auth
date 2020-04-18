package authbackendrpc

import (
	"github.com/roleypoly/db/ent"
	"github.com/roleypoly/gripkit"
	authBackend "github.com/roleypoly/rpc/auth/backend"
)

func RegisterRPCHandlers(gk *gripkit.Gripkit, svc *AuthBackendService) {
	authBackend.RegisterAuthBackendServer(gk.Server, svc)
}

type AuthBackendService struct {
	authBackend.UnimplementedAuthBackendServer
	db *ent.Client
}

func NewAuthBackendService(db *ent.Client) *AuthBackendService {
	return &AuthBackendService{
		db: db,
	}
}
