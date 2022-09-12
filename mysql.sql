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
    `id`          int(11)      NOT NULL AUTO_INCREMENT,
    `log_id`      bigint(20)   NOT NULL DEFAULT 0,
    `server_name` varchar(64)  NOT NULL DEFAULT '',
    `ip`          varchar(16)  NOT NULL DEFAULT '',
    `group`       varchar(256) NOT NULL DEFAULT '',
    `name`        varchar(256) NOT NULL DEFAULT '',
    `value`       double       NOT NULL DEFAULT 0,
    `data`        text         NOT NULL DEFAULT '',
    `create_time` datetime     NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_log_id` (`log_id`) USING BTREE,
    KEY `idx_server_name` (`server_name`) USING BTREE,
    KEY `idx_group` (`group`) USING BTREE,
    KEY `idx_name` (`name`) USING BTREE,
    KEY `idx_create_time` (`create_time`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;