CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS app_row (
    app_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_app_row_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS environment_row (
    env_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_environment_row_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS config_group_row (
    group_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID NOT NULL,
    env_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_config_group_row_app FOREIGN KEY (app_id) REFERENCES app_row(app_id) ON DELETE CASCADE,
    CONSTRAINT fk_config_group_row_env FOREIGN KEY (env_id) REFERENCES environment_row(env_id) ON DELETE CASCADE,
    CONSTRAINT uk_config_group_row_app_env_name UNIQUE (app_id, env_id, name)
);

CREATE TABLE IF NOT EXISTS config_item_row (
    item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL DEFAULT '',
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_config_item_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE,
    CONSTRAINT uk_config_item_row_group_key UNIQUE (group_id, key)
);

CREATE TABLE IF NOT EXISTS config_version_row (
    version_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    version_no BIGINT NOT NULL,
    draft_flag BOOLEAN NOT NULL DEFAULT TRUE,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_config_version_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE,
    CONSTRAINT uk_config_version_row_group_version_no UNIQUE (group_id, version_no)
);

CREATE TABLE IF NOT EXISTS release_record_row (
    release_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    version_id UUID NOT NULL,
    release_type TEXT NOT NULL,
    publish_scope JSONB NOT NULL DEFAULT '{}'::jsonb,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_release_record_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE,
    CONSTRAINT fk_release_record_row_version FOREIGN KEY (version_id) REFERENCES config_version_row(version_id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS gray_rule_row (
    gray_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    rule_type TEXT NOT NULL,
    rule_value JSONB NOT NULL DEFAULT '{}'::jsonb,
    target_scope JSONB NOT NULL DEFAULT '{}'::jsonb,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_gray_rule_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscription_row (
    subscription_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id TEXT NOT NULL,
    app_id UUID NOT NULL,
    env_id UUID NOT NULL,
    group_id UUID NOT NULL,
    checkpoint BIGINT NOT NULL DEFAULT 0,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_subscription_row_app FOREIGN KEY (app_id) REFERENCES app_row(app_id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_row_env FOREIGN KEY (env_id) REFERENCES environment_row(env_id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS gray_record_row (
    record_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL,
    user_id TEXT NOT NULL,
    matched_rule_id UUID NOT NULL,
    matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    result BOOLEAN NOT NULL,
    CONSTRAINT fk_gray_record_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS config_snapshot_row (
    snapshot_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID NOT NULL,
    env_id UUID NOT NULL,
    group_id UUID NOT NULL,
    items JSONB NOT NULL DEFAULT '{}'::jsonb,
    version_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_config_snapshot_row_app FOREIGN KEY (app_id) REFERENCES app_row(app_id) ON DELETE CASCADE,
    CONSTRAINT fk_config_snapshot_row_env FOREIGN KEY (env_id) REFERENCES environment_row(env_id) ON DELETE CASCADE,
    CONSTRAINT fk_config_snapshot_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE,
    CONSTRAINT fk_config_snapshot_row_version FOREIGN KEY (version_id) REFERENCES config_version_row(version_id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS change_event_row (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID NOT NULL,
    env_id UUID NOT NULL,
    group_id UUID NOT NULL,
    version_id UUID NOT NULL,
    version_no BIGINT NOT NULL,
    event_type TEXT NOT NULL,
    gray_flag BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_change_event_row_app FOREIGN KEY (app_id) REFERENCES app_row(app_id) ON DELETE CASCADE,
    CONSTRAINT fk_change_event_row_env FOREIGN KEY (env_id) REFERENCES environment_row(env_id) ON DELETE CASCADE,
    CONSTRAINT fk_change_event_row_group FOREIGN KEY (group_id) REFERENCES config_group_row(group_id) ON DELETE CASCADE,
    CONSTRAINT fk_change_event_row_version FOREIGN KEY (version_id) REFERENCES config_version_row(version_id) ON DELETE RESTRICT,
    CONSTRAINT uk_change_event_row_event_id UNIQUE (event_id)
);

CREATE INDEX IF NOT EXISTS idx_config_group_row_app_env ON config_group_row (app_id, env_id);
CREATE INDEX IF NOT EXISTS idx_config_item_row_group ON config_item_row (group_id);
CREATE INDEX IF NOT EXISTS idx_config_version_row_group_created_at ON config_version_row (group_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_release_record_row_group_created_at ON release_record_row (group_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_gray_rule_row_group_status ON gray_rule_row (group_id, status);
CREATE INDEX IF NOT EXISTS idx_subscription_row_app_env_group ON subscription_row (app_id, env_id, group_id);
CREATE INDEX IF NOT EXISTS idx_gray_record_row_group_matched_at ON gray_record_row (group_id, matched_at DESC);
CREATE INDEX IF NOT EXISTS idx_config_snapshot_row_group_created_at ON config_snapshot_row (group_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_change_event_row_group_created_at ON change_event_row (group_id, created_at DESC);
