-- 启用 Citus 扩展
CREATE
    EXTENSION IF NOT EXISTS citus;

-- 使用重试逻辑添加 worker 节点
DO
$$
    DECLARE
        retry_count  INTEGER := 0;
        max_retries  INTEGER := 30;
        worker_added BOOLEAN := FALSE;
    BEGIN
        WHILE retry_count < max_retries AND NOT worker_added
            LOOP
                BEGIN
                    -- 尝试添加 Worker 节点
                    PERFORM citus_add_node('citus-worker1', 5432);
                    PERFORM citus_add_node('citus-worker2', 5432);
                    PERFORM citus_add_node('citus-worker3', 5432);

                    worker_added := TRUE;
                    RAISE NOTICE 'Successfully added all worker nodes to the cluster';

                EXCEPTION
                    WHEN OTHERS THEN
                        retry_count := retry_count + 1;
                        RAISE NOTICE 'Attempt % failed, retrying in 2 seconds...', retry_count;
                        PERFORM pg_sleep(2);
                END;
            END LOOP;

        IF NOT worker_added THEN
            RAISE EXCEPTION 'Failed to add worker nodes after % attempts', max_retries;
        END IF;
    END
$$;

-- 验证集群状态
SELECT *
FROM citus_get_active_worker_nodes();

-- 用户表
CREATE TABLE users
(
    user_id        UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    email          VARCHAR(255) UNIQUE      NOT NULL,
    password_hash  VARCHAR(255),
    name           VARCHAR(255)             NOT NULL,
    username       VARCHAR(100) UNIQUE,
    avatar_url     VARCHAR(500),
    email_verified BOOLEAN                  NOT NULL DEFAULT false,
    github_id      VARCHAR(100),
    google_id      VARCHAR(100),
    gitlab_id      VARCHAR(100),
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_login_at  TIMESTAMP WITH TIME ZONE,
    status         VARCHAR(20)              NOT NULL DEFAULT 'active',
    UNIQUE (github_id),
    UNIQUE (google_id),
    UNIQUE (gitlab_id),
    CONSTRAINT username_format CHECK (username ~ '^[a-zA-Z0-9_-]{3,30}$' OR username IS NULL),
    CONSTRAINT valid_user_status CHECK (status IN ('active', 'inactive', 'suspended', 'deleted'))
);

-- 用户表索引
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_github_id ON users (github_id) WHERE github_id IS NOT NULL;
CREATE INDEX idx_users_google_id ON users (google_id) WHERE google_id IS NOT NULL;
CREATE INDEX idx_users_gitlab_id ON users (gitlab_id) WHERE gitlab_id IS NOT NULL;
CREATE INDEX idx_users_username ON users (username) WHERE username IS NOT NULL;
CREATE INDEX idx_users_status ON users (status);




-- 验证分片配置
SELECT table_name, citus_table_type, distribution_column, shard_count
FROM citus_tables
ORDER BY table_name;
