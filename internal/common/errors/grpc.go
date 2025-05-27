package errors

import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// MapToGRPC 将领域错误映射为gRPC错误
func MapToGRPC(err error) error {
    if err == nil {
        return nil
    }

    domainErr, ok := err.(*DomainError)
    if !ok {
        return status.Error(codes.Internal, "Internal server error")
    }

    grpcCode := mapTypeToGRPCCode(domainErr.Type)
    return status.Error(grpcCode, domainErr.Message)
}

// mapTypeToGRPCCode 映射错误类型到gRPC状态码
func mapTypeToGRPCCode(errorType ErrorType) codes.Code {
    switch errorType {
    case ErrorTypeValidation:
        return codes.InvalidArgument
    case ErrorTypeNotFound:
        return codes.NotFound
    case ErrorTypeAlreadyExists:
        return codes.AlreadyExists
    case ErrorTypeUnauthorized:
        return codes.Unauthenticated
    case ErrorTypeForbidden:
        return codes.PermissionDenied
    case ErrorTypeRateLimit:
        return codes.ResourceExhausted
    case ErrorTypeExternal:
        return codes.Unavailable
    default: // ErrorTypeInternal
        return codes.Internal
    }
}