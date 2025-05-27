-- 启用 Citus 扩展
CREATE EXTENSION IF NOT EXISTS citus;

-- 使用重试逻辑添加 worker 节点
DO $$
DECLARE
    retry_count INTEGER := 0;
    max_retries INTEGER := 30;
    worker_added BOOLEAN := FALSE;
BEGIN
    WHILE retry_count < max_retries AND NOT worker_added LOOP
        BEGIN
            -- 尝试添加 Worker 节点
            PERFORM citus_add_node('citus-worker1', 5432);
            PERFORM citus_add_node('citus-worker2', 5432);
            PERFORM citus_add_node('citus-worker3', 5432);
            
            worker_added := TRUE;
            RAISE NOTICE 'Successfully added all worker nodes to the cluster';
            
        EXCEPTION WHEN OTHERS THEN
            retry_count := retry_count + 1;
            RAISE NOTICE 'Attempt % failed, retrying in 2 seconds...', retry_count;
            PERFORM pg_sleep(2);
        END;
    END LOOP;
    
    IF NOT worker_added THEN
        RAISE EXCEPTION 'Failed to add worker nodes after % attempts', max_retries;
    END IF;
END $$;

-- 验证集群状态
SELECT * FROM citus_get_active_worker_nodes();

-- 计划配置（引用表）
CREATE TABLE plans (
    plan_type VARCHAR(20) PRIMARY KEY, -- free, pro, enterprise
    name VARCHAR(100) NOT NULL,
    max_teams INTEGER NOT NULL DEFAULT 1, -- 最大团队数
    max_members_per_team INTEGER NOT NULL DEFAULT 5, -- 每个团队最大成员数
    max_projects_per_team INTEGER NOT NULL DEFAULT 10, -- 每个团队最大项目数
    max_invited_users INTEGER NOT NULL DEFAULT 5, -- 用户总共能邀请的成员数
    max_api_calls_monthly INTEGER NOT NULL DEFAULT 1000, -- 每月API调用限制
    price_monthly DECIMAL(10,2) NOT NULL DEFAULT 0,
    features JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 用户表（引用表）
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE,
    avatar_url VARCHAR(500),
    email_verified BOOLEAN NOT NULL DEFAULT false,
    github_id VARCHAR(100),
    google_id VARCHAR(100),
    gitlab_id VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    -- OAuth ID 唯一约束
    UNIQUE(github_id),
    UNIQUE(google_id),
    UNIQUE(gitlab_id),
    -- 用户名格式约束
    CONSTRAINT username_format CHECK (username ~ '^[a-zA-Z0-9_-]{3,30}$' OR username IS NULL),
    -- 状态约束
    CONSTRAINT valid_user_status CHECK (status IN ('active', 'inactive', 'suspended', 'deleted'))
);

-- 用户订阅表（引用表）
CREATE TABLE user_subscriptions (
    subscription_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- 订阅用户
    plan_type VARCHAR(20) NOT NULL, -- 当前计划
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 订阅状态
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE, -- 到期时间，NULL表示永久
    auto_renew BOOLEAN NOT NULL DEFAULT true, --自动续费
    payment_method JSONB, -- 支付方式信息
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- 状态约束
    CONSTRAINT valid_subscription_status CHECK (status IN ('active', 'expired', 'cancelled', 'suspended')),
    -- 确保一个用户同时只有一个活跃订阅
    UNIQUE(user_id) DEFERRABLE INITIALLY DEFERRED
);

-- 团队表（按 owner_id 分片）
CREATE TABLE teams (
    team_id UUID DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL, -- 分片键：团队所有者
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (owner_id, team_id), -- 分片键在主键中
    -- 状态约束
    CONSTRAINT valid_team_status CHECK (status IN ('active', 'archived', 'deleted')),
    -- 团队名称在同一用户下唯一
    UNIQUE(owner_id, name)
);

-- 团队成员表（按 owner_id 分片）
CREATE TABLE team_members (
    owner_id UUID NOT NULL, -- 分片键：团队所有者
    team_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    invited_by UUID, -- 谁邀请的
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (owner_id, team_id, user_id),
    -- 角色约束
    CONSTRAINT valid_team_role CHECK (role IN ('owner', 'admin', 'member')),
    -- 状态约束
    CONSTRAINT valid_member_status CHECK (status IN ('active', 'inactive', 'pending', 'removed'))
);

-- 项目表（按 owner_id 分片）
CREATE TABLE projects (
    project_id UUID DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL, -- 分片键：团队所有者
    team_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (owner_id, project_id),
    -- 状态约束
    CONSTRAINT valid_project_status CHECK (status IN ('active', 'archived', 'deleted')),
    -- 项目名称在同一团队下唯一
    UNIQUE(owner_id, team_id, name)
);

-- 项目成员表（按 owner_id 分片）
CREATE TABLE project_members (
    owner_id UUID NOT NULL, -- 分片键：团队所有者
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    added_by UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (owner_id, project_id, user_id),
    -- 角色约束
    CONSTRAINT valid_project_role CHECK (role IN ('admin', 'member', 'viewer')),
    -- 状态约束
    CONSTRAINT valid_project_member_status CHECK (status IN ('active', 'inactive', 'removed'))
);

-- 使用量统计表（按 user_id 分片）
CREATE TABLE usage_stats (
    user_id UUID NOT NULL, -- 分片键：用户ID
    metric_name VARCHAR(50) NOT NULL,
    current_value INTEGER NOT NULL DEFAULT 0,
    period_start DATE NOT NULL,
    period_end DATE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, metric_name, period_start),
    -- 指标约束
    CONSTRAINT valid_metric_name CHECK (metric_name IN (
        'teams', 'invited_users', 'total_projects', 'api_calls'
    )),
    -- 数值约束
    CONSTRAINT non_negative_value CHECK (current_value >= 0),
    -- 日期约束
    CONSTRAINT valid_period CHECK (period_end IS NULL OR period_end >= period_start)
);

-- 设置引用表
SELECT create_reference_table('plans');
SELECT create_reference_table('users');
SELECT create_reference_table('user_subscriptions');

-- 设置分布式表（按 owner_id/user_id 分片）
SELECT create_distributed_table('teams', 'owner_id');
SELECT create_distributed_table('team_members', 'owner_id');
SELECT create_distributed_table('projects', 'owner_id');
SELECT create_distributed_table('project_members', 'owner_id');
SELECT create_distributed_table('usage_stats', 'user_id');

-- 引用表索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_github_id ON users(github_id) WHERE github_id IS NOT NULL;
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
CREATE INDEX idx_users_gitlab_id ON users(gitlab_id) WHERE gitlab_id IS NOT NULL;
CREATE INDEX idx_users_username ON users(username) WHERE username IS NOT NULL;
CREATE INDEX idx_users_status ON users(status);

CREATE INDEX idx_subscriptions_user_id ON user_subscriptions(user_id);
CREATE INDEX idx_subscriptions_status ON user_subscriptions(status);
CREATE INDEX idx_subscriptions_expires_at ON user_subscriptions(expires_at) WHERE expires_at IS NOT NULL;

-- 分布式表索引
CREATE INDEX idx_teams_status ON teams(status);
CREATE INDEX idx_teams_name ON teams(name);

CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_team_members_status ON team_members(status);

CREATE INDEX idx_projects_team_id ON projects(team_id);
CREATE INDEX idx_projects_created_by ON projects(created_by);
CREATE INDEX idx_projects_status ON projects(status);

CREATE INDEX idx_proj_members_user_id ON project_members(user_id);
CREATE INDEX idx_proj_members_status ON project_members(status);

CREATE INDEX idx_usage_stats_metric_period ON usage_stats(metric_name, period_start);

-- 插入默认计划数据
INSERT INTO plans (
    plan_type, 
    name, 
    max_teams, 
    max_members_per_team, 
    max_projects_per_team, 
    max_invited_users, 
    max_api_calls_monthly, 
    price_monthly, 
    features
) VALUES 
('free', 'Free Plan', 1, 3, 3, 3, 1000, 0.00, '{"basic_support": true}'),
('pro', 'Pro Plan', 5, 10, 50, 25, 10000, 29.00, '{"priority_support": true, "advanced_analytics": true, "custom_domains": true}'),
('enterprise', 'Enterprise Plan', 20, 50, 200, 100, 100000, 99.00, '{"custom_support": true, "sso": true, "advanced_analytics": true, "api_access": true}');