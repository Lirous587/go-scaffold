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
    avatar_url VARCHAR(500),
    email_verified BOOLEAN DEFAULT false,
    github_id VARCHAR(100), -- GitHub用户ID
    google_id VARCHAR(100), -- Google用户ID
    gitlab_id VARCHAR(100), -- GitLab用户ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'active',
    -- 确保OAuth ID唯一
    UNIQUE(github_id),
    UNIQUE(google_id),
    UNIQUE(gitlab_id),
    -- 添加用户名约束（字母数字下划线，长度限制）
    CONSTRAINT username_format CHECK (username ~ '^[a-zA-Z0-9_-]{3,30}$' OR username IS NULL)
);

-- 组织表（分片键）
CREATE TABLE organizations (
    org_id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL,
    plan_type VARCHAR(20) DEFAULT 'free',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    PRIMARY KEY (org_id) -- 分片键必须在主键中
);

-- 组织成员表（按 org_id 分片）
CREATE TABLE organization_members (
    org_id UUID NOT NULL, -- 分片键
    user_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- owner, admin, member
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    PRIMARY KEY (org_id, user_id)  -- 包含分片键的复合主键
);

-- 项目表（按 org_id 分片）
CREATE TABLE projects (
    project_id UUID DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL, -- 分片键
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    PRIMARY KEY (org_id, project_id)  -- 包含分片键的复合主键
);

-- 项目成员表（按 org_id 分片）
CREATE TABLE project_members (
    org_id UUID NOT NULL, -- 分片键
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role VARCHAR(20) DEFAULT 'member', -- admin, member
    added_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (org_id, project_id, user_id)  -- 包含分片键的复合主键
);

-- 使用量统计（按 org_id 分片）
CREATE TABLE usage_stats (
    org_id UUID NOT NULL, -- 分片键
    metric_name VARCHAR(50) NOT NULL, -- members, projects, api_calls
    current_value INTEGER DEFAULT 0,
    period_start DATE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (org_id, metric_name, period_start)  -- 包含分片键的复合主键
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

-- 为分布式表添加局部索引
CREATE INDEX idx_organizations_owner ON organizations(owner_id);
CREATE INDEX idx_organization_members_user_id ON organization_members(user_id);
CREATE INDEX idx_projects_created_by ON projects(created_by);
CREATE INDEX idx_project_members_user_id ON project_members(user_id);
CREATE INDEX idx_usage_stats_metric_period ON usage_stats(metric_name, period_start);

-- 插入默认计划数据
INSERT INTO plans (plan_type, name, max_members, max_projects, price_monthly, features) VALUES 
('free', 'Free Plan', 1, 5, 0.00, '{"basic_support": true}'),
('pro', 'Pro Plan', 10, 50, 29.00, '{"priority_support": true, "advanced_analytics": true}'),
('enterprise', 'Enterprise Plan', 20, 100, 99.00, '{"custom_support": true, "sso": true, "advanced_analytics": true}');
