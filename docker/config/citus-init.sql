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


-- 计划配置（引用表）- 必须先创建
CREATE TABLE plans (
    plan_type VARCHAR(20) PRIMARY KEY, -- free, pro, enterprise
    name VARCHAR(100) NOT NULL,
    max_members INTEGER DEFAULT 5,
    max_projects INTEGER DEFAULT 10,
    price_monthly DECIMAL(10,2) DEFAULT 0,
    features JSONB DEFAULT '{}'::jsonb
);

-- 用户表（引用表，所有节点都有完整副本）
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL, -- 主邮箱，用于账户识别
    password_hash VARCHAR(255), -- 可空，纯OAuth用户没有密码
    name VARCHAR(255) NOT NULL, -- 显示名称，支持中文等
    username VARCHAR(100) UNIQUE, -- 用户名，唯一标识，可空
    avatar_url VARCHAR(500), -- 头像URL，可空
    email_verified BOOLEAN NOT NULL DEFAULT false, -- 邮箱验证状态，必须有值
    github_id VARCHAR(100), -- GitHub用户ID，可空
    google_id VARCHAR(100), -- Google用户ID，可空
    gitlab_id VARCHAR(100), -- GitLab用户ID，可空
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- 创建时间，必须有值
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- 更新时间，必须有值
    last_login_at TIMESTAMP WITH TIME ZONE, -- 最后登录时间，可空
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 用户状态，必须有值
    -- 确保OAuth ID唯一（但允许为空）
    UNIQUE(github_id),
    UNIQUE(google_id),
    UNIQUE(gitlab_id),
    -- 添加用户名约束（字母数字下划线，长度限制）
    CONSTRAINT username_format CHECK (username ~ '^[a-zA-Z0-9_-]{3,30}$' OR username IS NULL),
    -- 状态约束
    CONSTRAINT valid_status CHECK (status IN ('active', 'inactive', 'suspended', 'deleted'))
);

-- 组织表（分片键）
CREATE TABLE organizations (
    org_id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL,
    plan_type VARCHAR(20) NOT NULL DEFAULT 'free',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (org_id), -- 分片键必须在主键中
    -- 计划类型约束
    CONSTRAINT valid_plan_type CHECK (plan_type IN ('free', 'pro', 'enterprise')),
    -- 状态约束
    CONSTRAINT valid_org_status CHECK (status IN ('active', 'suspended', 'deleted'))
);

-- 组织成员表（按 org_id 分片）
CREATE TABLE organization_members (
    org_id UUID NOT NULL, -- 分片键
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- owner, admin, member
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (org_id, user_id), -- 包含分片键的复合主键
    -- 角色约束
    CONSTRAINT valid_org_role CHECK (role IN ('owner', 'admin', 'member')),
    -- 状态约束
    CONSTRAINT valid_member_status CHECK (status IN ('active', 'inactive', 'pending', 'removed'))
);

-- 项目表（按 org_id 分片）
CREATE TABLE projects (
    project_id UUID DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL, -- 分片键
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (org_id, project_id), -- 包含分片键的复合主键
    -- 状态约束
    CONSTRAINT valid_project_status CHECK (status IN ('active', 'archived', 'deleted')),
    UNIQUE(org_id, name)
);

-- 项目成员表（按 org_id 分片）
CREATE TABLE project_members (
    org_id UUID NOT NULL, -- 分片键
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- admin, member
    added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    added_by UUID, -- 可空，系统自动添加的成员可能没有添加者
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    PRIMARY KEY (org_id, project_id, user_id), -- 包含分片键的复合主键
    -- 角色约束
    CONSTRAINT valid_project_role CHECK (role IN ('admin', 'member', 'viewer')),
    -- 状态约束
    CONSTRAINT valid_project_member_status CHECK (status IN ('active', 'inactive', 'removed'))
);

-- 使用量统计（按 org_id 分片）
CREATE TABLE usage_stats (
    org_id UUID NOT NULL, -- 分片键
    metric_name VARCHAR(50) NOT NULL, -- 用量指标 members, projects, api_calls 
    current_value INTEGER NOT NULL DEFAULT 0,
    period_start DATE NOT NULL,
    period_end DATE, -- 可空，当前期间可能还没结束
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (org_id, metric_name, period_start), -- 包含分片键的复合主键
    -- 指标名称约束
    CONSTRAINT valid_metric_name CHECK (metric_name IN ('members', 'projects', 'api_calls')),
    -- 数值约束
    CONSTRAINT non_negative_value CHECK (current_value >= 0),
    -- 日期约束
    CONSTRAINT valid_period CHECK (period_end IS NULL OR period_end >= period_start)
);

-- 设置引用表（在所有节点复制）
SELECT create_reference_table('plans');
SELECT create_reference_table('users');

-- 设置分布式表（按 org_id 分片）
SELECT create_distributed_table('organizations', 'org_id');
SELECT create_distributed_table('organization_members', 'org_id');
SELECT create_distributed_table('projects', 'org_id');
SELECT create_distributed_table('project_members', 'org_id');
SELECT create_distributed_table('usage_stats', 'org_id');


-- 为引用表添加索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_github_id ON users(github_id) WHERE github_id IS NOT NULL;
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
CREATE INDEX idx_users_gitlab_id ON users(gitlab_id) WHERE gitlab_id IS NOT NULL;
CREATE INDEX idx_users_username ON users(username) WHERE username IS NOT NULL;
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_last_login ON users(last_login_at) WHERE last_login_at IS NOT NULL;

-- 为分布式表添加局部索引
-- 组织表索引
CREATE INDEX idx_organizations_owner ON organizations(owner_id);
CREATE INDEX idx_organizations_plan_type ON organizations(plan_type);
CREATE INDEX idx_organizations_status ON organizations(status);
CREATE INDEX idx_organizations_created_at ON organizations(created_at);
CREATE INDEX idx_organizations_name ON organizations(name); -- 组织名称搜索

-- 组织成员表索引
CREATE INDEX idx_organization_members_user_id ON organization_members(user_id);
CREATE INDEX idx_organization_members_role ON organization_members(role);
CREATE INDEX idx_organization_members_status ON organization_members(status);
CREATE INDEX idx_organization_members_joined_at ON organization_members(joined_at);
-- 复合索引：查找用户的活跃组织
CREATE INDEX idx_org_members_user_status ON organization_members(user_id, status) WHERE status = 'active';
-- 复合索引：查找组织的某种角色成员
CREATE INDEX idx_org_members_role_status ON organization_members(role, status) WHERE status = 'active';

-- 项目表索引
CREATE INDEX idx_projects_created_by ON projects(created_by);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_updated_at ON projects(updated_at);
CREATE INDEX idx_projects_name ON projects(name); -- 项目名称搜索
-- 复合索引：组织内的活跃项目
CREATE INDEX idx_projects_org_status ON projects(org_id, status) WHERE status = 'active';

-- 项目成员表索引
CREATE INDEX idx_project_members_user_id ON project_members(user_id);
CREATE INDEX idx_project_members_role ON project_members(role);
CREATE INDEX idx_project_members_status ON project_members(status);
CREATE INDEX idx_project_members_added_at ON project_members(added_at);
CREATE INDEX idx_project_members_added_by ON project_members(added_by) WHERE added_by IS NOT NULL;
-- 复合索引：查找用户参与的活跃项目
CREATE INDEX idx_proj_members_user_status ON project_members(user_id, status) WHERE status = 'active';
-- 复合索引：项目内的活跃成员
CREATE INDEX idx_proj_members_project_status ON project_members(org_id, project_id, status) WHERE status = 'active';

-- 使用量统计表索引
CREATE INDEX idx_usage_stats_metric_period ON usage_stats(metric_name, period_start);
CREATE INDEX idx_usage_stats_updated_at ON usage_stats(updated_at);
CREATE INDEX idx_usage_stats_period_end ON usage_stats(period_end) WHERE period_end IS NOT NULL;
-- 复合索引：查询特定组织的使用量趋势
CREATE INDEX idx_usage_stats_org_metric_period ON usage_stats(org_id, metric_name, period_start);
-- 复合索引：查询当前活跃统计
CREATE INDEX idx_usage_stats_active_period ON usage_stats(metric_name, period_start) WHERE period_end IS NULL;

-- 计划表索引（虽然是引用表，数据量小，但也可以加一些）
CREATE INDEX idx_plans_price ON plans(price_monthly);
CREATE INDEX idx_plans_max_members ON plans(max_members);
CREATE INDEX idx_plans_max_projects ON plans(max_projects);

-- 插入默认计划数据
INSERT INTO plans (plan_type, name, max_members, max_projects, price_monthly, features) VALUES 
('free', 'Free Plan', 1, 5, 0.00, '{"basic_support": true}'),
('pro', 'Pro Plan', 10, 50, 29.00, '{"priority_support": true, "advanced_analytics": true}'),
('enterprise', 'Enterprise Plan', 20, 100, 99.00, '{"custom_support": true, "sso": true, "advanced_analytics": true}');
