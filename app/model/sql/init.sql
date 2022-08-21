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

create table wf_i18n_item
(
    `id`          bigint(20)   not null auto_increment comment '自增id',
    `unique_name` varchar(100) not null default '' comment '词条键',
    `lang`        varchar(10)  not null default '' comment '语言缩写',
    `name`        varchar(50)  not null default '' comment '词条名称',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_idx_key_lang` (`unique_name`, `lang`),
    UNIQUE KEY `uni_idx_name_lang` (`name`, `lang`)
) ENGINE = InnoDB
  DEFAULT CHARSET = UTF8MB4 COMMENT ='warframe词条国际化表';

create table wf_entry
(
    `id`          bigint(20)   not null auto_increment comment '自增id',
    `unique_name` varchar(100) not null default '' comment 'wf词条全局唯一名称',
    `category`    varchar(10)  not null default '' comment '词条分类',
    `name`        varchar(50)  not null default '' comment '词条英文名',
    `tradable`    tinyint(1)   not null default 0 comment '是否能交易',
    `create_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_idx_unique_name` (`unique_name`),
    KEY `idx_name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = UTF8MB4 COMMENT ='warframe词条表';