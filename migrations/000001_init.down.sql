DROP INDEX IF EXISTS idx_change_event_row_group_created_at;
DROP INDEX IF EXISTS idx_config_snapshot_row_group_created_at;
DROP INDEX IF EXISTS idx_gray_record_row_group_matched_at;
DROP INDEX IF EXISTS idx_subscription_row_app_env_group;
DROP INDEX IF EXISTS idx_gray_rule_row_group_status;
DROP INDEX IF EXISTS idx_release_record_row_group_created_at;
DROP INDEX IF EXISTS idx_config_version_row_group_created_at;
DROP INDEX IF EXISTS idx_config_item_row_group;
DROP INDEX IF EXISTS idx_config_group_row_app_env;

DROP TABLE IF EXISTS change_event_row;
DROP TABLE IF EXISTS config_snapshot_row;
DROP TABLE IF EXISTS gray_record_row;
DROP TABLE IF EXISTS subscription_row;
DROP TABLE IF EXISTS gray_rule_row;
DROP TABLE IF EXISTS release_record_row;
DROP TABLE IF EXISTS config_version_row;
DROP TABLE IF EXISTS config_item_row;
DROP TABLE IF EXISTS config_group_row;
DROP TABLE IF EXISTS environment_row;
DROP TABLE IF EXISTS app_row;

DROP EXTENSION IF EXISTS pgcrypto;
