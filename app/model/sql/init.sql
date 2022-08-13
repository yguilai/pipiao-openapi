use pipiao;
create table openapi_auth
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `app_id`      varchar(55)  NOT NULL DEFAULT '' COMMENT 'appId',
    `app_key`     varchar(255) NOT NULL DEFAULT '' COMMENT 'appKey',
    `status`      tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态',
    `create_time` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_idx_app_id`(`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Openapi鉴权表';

create table wf_dict
(
    `id`          bigint(20) not null auto_increment comment '自增id',
    `raw`         varchar(255) not null default '' comment '原始内容',
    `target`      varchar(255) not null default '' comment '翻译结果',
    `type`        int(5) not null default 0 comment '翻译类型',
    `create_time` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY           `idx_raw`(`raw`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='warframe词典表';