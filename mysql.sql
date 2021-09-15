CREATE TABLE `server_conf`
(
    `id`          int(11)     NOT NULL AUTO_INCREMENT,
    `server_name` varchar(32) NOT NULL DEFAULT '' COMMENT '服务名称',
    `version`     int(11)     NOT NULL DEFAULT 0 COMMENT '版本号',
    `remark`      varchar(64) NOT NULL DEFAULT '' COMMENT '备注',
    `conf_text`   text        NOT NULL DEFAULT '' COMMENT '',
    `created_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_name_ver` (`server_name`, `version`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;