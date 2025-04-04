package app_grpc_register

import (
	v1alpha1 "github.com/agntcy/identity/internal/pkg/generated/agntcy/identity/v1alpha1"
)

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GrpcServiceRegister struct {
	DidServiceServer v1alpha1.DidServiceServer

	IssuerServiceServer v1alpha1.IssuerServiceServer

	VcServiceServer v1alpha1.VcServiceServer
}

func (r GrpcServiceRegister) RegisterGrpcHandlers(grpcServer *grpc.Server) {

	if r.DidServiceServer != nil {
		v1alpha1.RegisterDidServiceServer(grpcServer, r.DidServiceServer)
	}

	if r.IssuerServiceServer != nil {
		v1alpha1.RegisterIssuerServiceServer(grpcServer, r.IssuerServiceServer)
	}

	if r.VcServiceServer != nil {
		v1alpha1.RegisterVcServiceServer(grpcServer, r.VcServiceServer)
	}

}

func (r GrpcServiceRegister) RegisterHttpHandlers(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {

	if r.DidServiceServer != nil {
		err := v1alpha1.RegisterDidServiceHandler(ctx, mux, conn)
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

	if r.VcServiceServer != nil {
		err := v1alpha1.RegisterVcServiceHandler(ctx, mux, conn)
		if err != nil {
			return err
		}
	}

	return nil
}
