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