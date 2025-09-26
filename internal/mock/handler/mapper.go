package handler

import (
    "time"
    "scaffold/internal/mock/domain"
)

type MockResponse struct {
    ID          int64  `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
    DeletedAt   string `json:"deleted_at,omitempty"`
}

func domainMockToResponse(mock *domain.Mock) *MockResponse {
    if mock == nil {
        return nil
    }

    return &MockResponse{
        ID:          mock.ID,
        Title:       mock.Title,
        Description: mock.Description,
        CreatedAt:   mock.CreatedAt.Format(time.DateTime),
        UpdatedAt:   mock.UpdatedAt.Format(time.DateTime),
        DeletedAt:   mock.DeletedAt.Format(time.DateTime),
    }
}

func domainMocksToResponse(mocks []*domain.Mock) []*MockResponse {
    if len(mocks) == 0 {
        return nil
    }

    ret := make([]*MockResponse, 0, len(mocks))

    for _, mock := range mocks {
        if mock != nil {
            ret = append(ret, domainMockToResponse(mock))
        }
    }
    return ret
}

type CreateRequest struct {
    Title       string  `json:"title" binding:"required,max=30"`
    Description string  `json:"Description" binding:"max=60"`
}

type UpdateRequest struct {
    Title       string  `json:"title" binding:"required,max=30"`
    Description string  `json:"Description" binding:"max=60"`
}

type ListRequest struct {
    Page     int    `form:"page,default=1" binding:"min=1"`
    PageSize int    `form:"page_size,default=5" binding:"min=5,max=20"`
    KeyWord  string `form:"keyword" binding:"max=20"`
}

type MockListResponse struct {
    Total int64                         `json:"total"`
    List  []*MockResponse   `json:"list"`
}

func domainMockListToResponse(data *domain.MockList) *MockListResponse {
    if data == nil {
        return nil
    }

    return &MockListResponse{
        Total: data.Total,
        List:  domainMocksToResponse(data.List),
    }
}
