create table `server_conf`
(
    `id`          int(11)     not null auto_increment,
    `server_name` varchar(32) not null default '' comment '服务名称',
    `version`     int(11)     not null default 0 comment '版本号',
    `remark`      varchar(64) not null default '' comment '备注',
    `conf_text`   text        not null default '' comment '',
    `created_at`  datetime    not null default current_timestamp,
    `updated_at`  datetime    not null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `uniq_name_ver` (`server_name`, `version`) using btree
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_unicode_ci;


CREATE TABLE `event`
(
    `id`          int(11)     NOT NULL AUTO_INCREMENT,
    `log_id`      bigint(20)  NOT NULL DEFAULT 0 COMMENT '',
    `server_name` varchar(64) NOT NULL DEFAULT '' COMMENT '',
    `ip`          varchar(16) NOT NULL DEFAULT '' COMMENT '',
    `event_group` varchar(64) NOT NULL DEFAULT '' COMMENT '',
    `event_name`  varchar(64) NOT NULL DEFAULT '' COMMENT '',
    `value`       double      NOT NULL DEFAULT 0 COMMENT '',
    `data`        text        NOT NULL DEFAULT '' COMMENT '',
    `created_at`  datetime    NOT NULL,
    `updated_at`  datetime    NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_log_id` (`log_id`) USING BTREE,
    KEY `idx_server_name` (`server_name`) USING BTREE,
    KEY `idx_event_group` (`event_group`) USING BTREE,
    KEY `idx_event_name` (`event_name`) USING BTREE,
    KEY `idx_created_at` (`created_at`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;