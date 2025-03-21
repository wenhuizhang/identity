package app_grpc_register

import (
	v1alpha1 "github.com/agntcy/pyramid/internal/pkg/generated/agntcy/pyramid/v1alpha1"
)

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GrpcServiceRegister struct {
	CredentialsServiceServer v1alpha1.CredentialsServiceServer

	IdServiceServer v1alpha1.IdServiceServer

	IssuerServiceServer v1alpha1.IssuerServiceServer
}

func (r GrpcServiceRegister) RegisterGrpcHandlers(grpcServer *grpc.Server) {

	if r.CredentialsServiceServer != nil {
		v1alpha1.RegisterCredentialsServiceServer(grpcServer, r.CredentialsServiceServer)
	}

	if r.IdServiceServer != nil {
		v1alpha1.RegisterIdServiceServer(grpcServer, r.IdServiceServer)
	}

	if r.IssuerServiceServer != nil {
		v1alpha1.RegisterIssuerServiceServer(grpcServer, r.IssuerServiceServer)
	}

}

func (r GrpcServiceRegister) RegisterHttpHandlers(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {

	if r.CredentialsServiceServer != nil {
		err := v1alpha1.RegisterCredentialsServiceHandler(ctx, mux, conn)
		if err != nil {
			return err
		}
	}

	if r.IdServiceServer != nil {
		err := v1alpha1.RegisterIdServiceHandler(ctx, mux, conn)
		if err != nil {
			return err
		}
	}

	if r.IssuerServiceServer != nil {
		err := v1alpha1.RegisterIssuerServiceHandler(ctx, mux, conn)
		if err != nil {
			return err
		}
	}

	return nil
}
