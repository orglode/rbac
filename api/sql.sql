CREATE TABLE `account` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '账号姓名',
`mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '电话',
`account_number` varchar(50) NOT NULL DEFAULT '' COMMENT '登录账号',
`password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码',
`account_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '账户类型 // 1 管理员 ',
`last_login_time` int(11) NOT NULL DEFAULT 0 COMMENT '最后一次登录时间',
`opertion_id` int(11) NOT NULL DEFAULT 0 COMMENT '操作人',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
`update_time` int(11)  NOT NULL DEFAULT 0 COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `account_number` (`account_number`),
UNIQUE KEY `mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='后台用户表';

CREATE TABLE `role` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '账号姓名',
    `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户角色表';

CREATE TABLE `user_role` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`account_id` int(11) NOT NULL DEFAULT 0 COMMENT '账号ID',
`role_id` int(11) NOT NULL DEFAULT 0 COMMENT '角色ID',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户角色权限关闭表';

CREATE TABLE `role_page` (
                             `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                             `role_id` int(11) NOT NULL DEFAULT 0 COMMENT '账号ID',
                             `page_id` int(11) NOT NULL DEFAULT 0 COMMENT '账号ID',
                             PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='角色权限关联表';


CREATE TABLE `page` (
                        `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                        `page_name` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单名称',
                        `url` varchar(255) NOT NULL DEFAULT '' COMMENT '路由',
                        `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
                        PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='页面菜单栏';