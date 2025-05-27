package handler

import (
    "context"

    "google.golang.org/protobuf/types/known/timestamppb"
    
    "scaffold/internal/common/response"
    "scaffold/internal/user/domain"
    pb "scaffold/internal/user/proto"
)

type GrpcHandler struct {
    pb.UnimplementedUserServiceServer
    userService domain.UserService
}

func NewGrpcHandler(userService domain.UserService) *GrpcHandler {
    return &GrpcHandler{
        userService: userService,
    }
}

// 认证相关
func (h *GrpcHandler) AuthenticateWithOAuth(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
    if req.OauthUserInfo == nil {
        return nil, response.HandleGRPCError(domain.ErrOAuthUserInfoMissing)
    }

    userInfo := &domain.OAuthUserInfo{
        Provider: req.OauthUserInfo.Provider,
        ID:       req.OauthUserInfo.Id,
        Login:    req.OauthUserInfo.Login,
        Name:     req.OauthUserInfo.Name,
        Email:    req.OauthUserInfo.Email,
        Avatar:   req.OauthUserInfo.Avatar,
    }

    session, err := h.userService.AuthenticateWithOAuth(req.Provider, userInfo)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    return &pb.AuthenticateResponse{
        Session: domainSessionToProto(session),
    }, nil
}

func (h *GrpcHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
    session, err := h.userService.RefreshUserSession(req.UserId, req.RefreshToken)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    return &pb.RefreshTokenResponse{
        Session: domainSessionToProto(session),
    }, nil
}

// 用户管理
func (h *GrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := h.userService.GetUser(req.UserId)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    return &pb.GetUserResponse{
        User: domainUserToProto(user),
    }, nil
}

func (h *GrpcHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
    if req.Updates == nil {
        validationErr := domain.ErrOAuthUserInfoMissing.
            WithDetail("field", "updates").
            WithDetail("message", "updates field is required")
        return nil, response.HandleGRPCError(validationErr)
    }

    updates := &domain.UserProfileUpdate{}
    if req.Updates.Name != nil {
        updates.Name = req.Updates.Name
    }
    if req.Updates.Username != nil {
        updates.Username = req.Updates.Username
    }
    if req.Updates.Avatar != nil {
        updates.Avatar = req.Updates.Avatar
    }

    user, err := h.userService.UpdateUserProfile(req.UserId, updates)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    return &pb.UpdateUserProfileResponse{
        User: domainUserToProto(user),
    }, nil
}

// 团队管理
func (h *GrpcHandler) CreateTeam(ctx context.Context, req *pb.CreateTeamRequest) (*pb.CreateTeamResponse, error) {
    if req.TeamInfo == nil {
        validationErr := domain.ErrOAuthUserInfoMissing.
            WithDetail("field", "team_info").
            WithDetail("message", "team_info field is required")
        return nil, response.HandleGRPCError(validationErr)
    }

    teamInfo := &domain.TeamCreateRequest{
        Name:        req.TeamInfo.Name,
        Description: req.TeamInfo.Description,
    }

    team, err := h.userService.CreateTeam(req.OwnerId, teamInfo)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    return &pb.CreateTeamResponse{
        Team: domainTeamToProto(team),
    }, nil
}

func (h *GrpcHandler) GetUserTeams(ctx context.Context, req *pb.GetUserTeamsRequest) (*pb.GetUserTeamsResponse, error) {
    teams, err := h.userService.GetUserTeams(req.UserId)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    protoTeams := make([]*pb.Team, len(teams))
    for i, team := range teams {
        protoTeams[i] = domainTeamToProto(team)
    }

    return &pb.GetUserTeamsResponse{
        Teams: protoTeams,
    }, nil
}

func (h *GrpcHandler) JoinTeam(ctx context.Context, req *pb.JoinTeamRequest) (*pb.JoinTeamResponse, error) {
    err := h.userService.JoinTeam(req.UserId, req.TeamId)
    if err != nil {
        return nil, response.HandleGRPCError(err)
    }

    return &pb.JoinTeamResponse{
        Success: true,
    }, nil
}

// 转换函数保持不变
func domainUserToProto(user *domain.User) *pb.User {
    if user == nil {
        return nil
    }

    protoUser := &pb.User{
        Id:            user.ID,
        Email:         user.Email,
        Name:          user.Name,
        Username:      user.Username,
        AvatarUrl:     user.AvatarURL,
        EmailVerified: user.EmailVerified,
        Status:        user.Status,
        CreatedAt:     timestamppb.New(user.CreatedAt),
        UpdatedAt:     timestamppb.New(user.UpdatedAt),
    }

    if user.LastLoginAt != nil {
        protoUser.LastLoginAt = timestamppb.New(*user.LastLoginAt)
    }

    return protoUser
}

func domainTeamToProto(team *domain.Team) *pb.Team {
    if team == nil {
        return nil
    }

    return &pb.Team{
        Id:          team.ID,
        OwnerId:     team.OwnerID,
        Name:        team.Name,
        Description: team.Description,
        Status:      team.Status,
        CreatedAt:   timestamppb.New(team.CreatedAt),
        UpdatedAt:   timestamppb.New(team.UpdatedAt),
    }
}

func domainSessionToProto(session *domain.UserSession) *pb.UserSession {
    if session == nil {
        return nil
    }

    return &pb.UserSession{
        User:         domainUserToProto(session.User),
        AccessToken:  session.AccessToken,
        RefreshToken: session.RefreshToken,
        ExpiresAt:    timestamppb.New(session.ExpiresAt),
    }
}