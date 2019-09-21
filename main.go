package main // import "github.com/roleypoly/auth"

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/joho/godotenv/autoload"
	"github.com/roleypoly/gripkit"
	"google.golang.org/grpc"

	pbBackend "github.com/roleypoly/rpc/auth/backend"
	pbClient "github.com/roleypoly/rpc/auth/client"
)

var (
	discordClientID     = os.Getenv("DISCORD_CLIENT_ID")
	discordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")

	rootUsers    = parseRoot(os.Getenv("ROOT_USERS"))
	svcPort      = os.Getenv("AUTH_SVC_PORT")
	tlsCertPath  = os.Getenv("TLS_CERT_PATH")
	tlsKeyPath   = os.Getenv("TLS_KEY_PATH")
	sharedSecret = os.Getenv("SHARED_SECRET")
)

func parseRoot(s string) []string {
	return strings.Split(s, ",")
}

func main() {
	go startGripkit()

	fmt.Println("roleypoly-auth: started grpc on", svcPort)

	syscallExit := make(chan os.Signal, 1)
	signal.Notify(
		syscallExit,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
		os.Kill,
	)
	<-syscallExit
}

func defaultAuthFunc(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func startGripkit() {
	grpcAuthBackend := &AuthBackendService{}

	grpcAuthClient := &AuthClientService{}
	gk := gripkit.Create(
		gripkit.WithHTTPOptions(gripkit.HTTPOptions{
			Addr:        svcPort,
			TLSCertPath: tlsCertPath,
			TLSKeyPath:  tlsKeyPath,
		}),
		gripkit.WithGrpcWeb(
			grpcweb.WithOriginFunc(func(o string) bool { return true }),
		),
		gripkit.WithDebug(),
		gripkit.WithOptions(
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_auth.UnaryServerInterceptor(defaultAuthFunc),
				),
			),
		),
	)
	pbClient.RegisterAuthClientServer(gk.Server, grpcAuthClient)
	pbBackend.RegisterAuthBackendServer(gk.Server, grpcAuthBackend)

	err := gk.Serve()
	if err != nil {
		log.Fatalln("grpc server failed to start", err)
	}
}
