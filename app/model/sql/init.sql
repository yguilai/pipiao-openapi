use pipiao;
-- goctl model mysql ddl --src "./init.sql" --dir "../" --style "goZero" --cache
create table openapi_auth
(
    `id`          bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `app_id`      varchar(55)  NOT NULL DEFAULT '' COMMENT 'appId',
    `app_key`     varchar(255) NOT NULL DEFAULT '' COMMENT 'appKey',
    `status`      tinyint(1)   NOT NULL DEFAULT 1 COMMENT '状态',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_idx_app_id` (`app_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = UTF8MB4 COMMENT ='Openapi鉴权表';

create table wf_item
(
    `id`          bigint(20)   not null auto_increment comment '自增id',
    `key`         varchar(100) not null default '' comment '词条键',
    `lang`        varchar(10)  not null default '' comment '语言缩写',
    `name`        varchar(50)  not null default '' comment '词条名称',
    `description` varchar(255) not null default '' comment '词条说明',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_idx_key_lang` (`key`, `lang`),
    UNIQUE KEY `uni_idx_name_lang` (`name`, `lang`)
) ENGINE = InnoDB
  DEFAULT CHARSET = UTF8MB4 COMMENT ='warframe词条表';