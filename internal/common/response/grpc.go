package response

import (
    commonErrors "scaffold/internal/common/errors"
)

// HandleGRPCError 处理gRPC错误
func HandleGRPCError(err error) error {
    return commonErrors.MapToGRPC(err)
}