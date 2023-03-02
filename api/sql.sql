
CREATE TABLE `module` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(50) NOT NULL DEFAULT '' COMMENT '模块名称',
`api_path` varchar(200) NOT NULL DEFAULT '/' COMMENT '后端接口',
`parent` int(11) NOT NULL DEFAULT '0' COMMENT '所属父级信息',
`sort` int(11) NOT NULL DEFAULT '1' COMMENT '排序 数值越大越靠后',
`code` int(10) NOT NULL DEFAULT '0' COMMENT '按钮标记 0--菜单 1--查询 2-添加 3--修改 4-删除 5-修改其他状态，6–导出，其他标记 ...',
`type` varchar(200) NOT NULL DEFAULT '1' COMMENT '前端区分按钮 页面 1–菜单 2--按钮',
`root` varchar(200) NOT NULL DEFAULT '' COMMENT '前端路由使用',
`status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1--启用 2--禁用 3--菜单上不显示',
`operator_uid` int(11) NOT NULL DEFAULT '0' COMMENT '操作者id',
`operator` varchar(30) NOT NULL COMMENT '操作者',
`create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
`update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
PRIMARY KEY (`id`),
KEY `sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='权限模块表';

insert into module (id, name, api_path, parent, sort, code, `type`, `root`, `status`, operator_uid, operator)
values  (1,'系统管理', '/', 0, 1, 0, '1', 'auth', 1, 0, 'system'),
        (2,'角色管理', '/', 1, 1, 1, '1', 'roleMs', 1, 0, 'system'),
        (3,'页面与按钮管理', '/', 1, 1, 2, '1', 'pageButtonMs', 1, 0, 'system'),
        (4,'用户管理', '/', 1, 1, 1, '1', 'userMs', 1, 0, 'system'),
        (5,'角色类型管理', '/', 1, 1, 1, '1', 'roleTypeMs', 1, 0, 'system'),
        (6,'日志管理', '/', 1, 1, 1, '1', 'logMs', 1, 0, 'system'),
        (7,'新增', '/login/rbac/role/role_modify', 2, 1, 2, '2', '', 1, 0, 'system'),
        (8,'查询', '/login/rbac/role/role_list', 2, 1, 1, '2', '', 1, 0, 'system'),
        (9,'编辑', '/login/rbac/role/role_modify', 2, 1, 3, '2', '', 1, 0, 'system'),
        (10,'新增', '/login/rbac/module/module_modify', 3, 1, 2, '2', '', 1, 0, 'system'),
        (11,'编辑', '/login/rbac/module/module_modify', 3, 1, 3, '2', '', 1, 0, 'system'),
        (12,'查询', '/login/rbac/users/user_list', 4, 1, 1, '2', '', 1, 0, 'system'),
        (13,'新增', '/login/rbac/role/role_type_add', 5, 1, 2, '2', '', 1, 0, 'system'),
        (14,'查询', '/login/rbac/role/role_type_list', 5, 1, 1, '2', '', 1, 0, 'system'),
        (15,'编辑', '/login/rbac/users/user_modify', 4, 1, 3, '2', '', 1, 0, 'system'),
        (16,'新增', '/login/rbac/users/user_modify', 4, 1, 2, '2', '', 1, 0, 'system'),
        (17,'导出', '/web/op_log/lists', 6, 1, 6, '2', '', 1, 0, 'system'),
        (18,'查询', '/web/op_log/lists', 6, 1, 1, '2', '', 1, 0, 'system'),
        (19,'查看明细', '/web/op_log/lists', 6, 1, 4, '2', '', 1, 0, 'system'),
        (20,'系统初始化', '/', 0, 1, 0, '1', 'authorityInitialization', 1, 0, 'system'),
        (21,'查询', '/login/rbac/system/lists', 20, 1, 1, '2', '', 1, 0, 'system'),
        (22,'新增', '/login/rbac/system/add', 20, 1, 2, '2', '', 1, 0, 'system'),
        (23,'初始化权限', '/login/rbac/module/module_modify', 20, 1, 9, '2', '', 1, 0, 'system'),
        (24,'初始化账号', '/login/rbac/users/user_modify', 20, 1, 10, '2', '', 1, 0, 'system');







CREATE TABLE `role_type` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`type_name` varchar(30) NOT NULL DEFAULT '' COMMENT '角色类型名称',
`status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1：有效 2：删除',
`create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
`update_time` int(11) NOT NULL DEFAULT '0' COMMENT '更新时间',
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='角色类型表';




CREATE TABLE `role` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`pid` varchar(200) NOT NULL DEFAULT '0' COMMENT '父级角色id，多个逗号分隔',
`type_id` smallint(6) NOT NULL DEFAULT '1' COMMENT '角色类型id',
`name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名称',
`description` varchar(200) NOT NULL DEFAULT '' COMMENT '角色描述',
`operator_id` int(11) NOT NULL DEFAULT '0' COMMENT '后台操作人id',
`operator` varchar(100) NOT NULL DEFAULT '' COMMENT '操作人名称',
`page_type_id` smallint(6) NOT NULL DEFAULT '0' COMMENT '页面分类id',
`status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1--启用 2--禁用',
`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='角色表（含数据范围）';



CREATE TABLE `role_module` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`role_id` smallint(6) NOT NULL DEFAULT '0' COMMENT '角色ID',
`module_id` int(11) NOT NULL DEFAULT '0' COMMENT '模块id',
UNIQUE KEY `role_id` (`role_id`,`module_id`)
) ENGINE=InnoDB  AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='角色权限映射关系表';



CREATE TABLE `page_type` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(50) NOT NULL DEFAULT '' COMMENT '页面分类名称',
`module_ids` varchar(2000) NOT NULL DEFAULT '0' COMMENT '模块ids',
`status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1：有效 2：删除',
`operator` varchar(30) NOT NULL DEFAULT '' COMMENT '操作者',
`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
KEY `name` (`name`),
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='页面分类表';


CREATE TABLE `url_skip` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`skip_url` varchar(200) NOT NULL DEFAULT '' COMMENT '后端需要跳过校验的接口',
`status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1：有效 2：无效',
`desc` varchar(100) NOT NULL DEFAULT '' COMMENT 'url描述',
PRIMARY KEY (`id`),
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='系统下需要跳过校验的特殊url表';

CREATE TABLE `users` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`account` varchar(255) NOT NULL DEFAULT '' COMMENT '账号',
`pass_word` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
`username` varchar(50) NOT NULL DEFAULT '' COMMENT '账号姓名',
`phone` varchar(50) NOT NULL DEFAULT '' COMMENT '电话',
`description` varchar(65) NOT NULL DEFAULT '' COMMENT '描述',
`status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 1:启用 2:禁用',
`operator` varchar(50) NOT NULL DEFAULT '' COMMENT '操作者',
`last_login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次登录时间',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
`update_time` int(11)  NOT NULL DEFAULT 0 COMMENT '更新时间',
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='系统用户表';



CREATE TABLE `inke_user_role` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户id',
`role_id` int(11) NOT NULL DEFAULT '0' COMMENT '角色id',
KEY `sel` (`user_id`,`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='用户角色关系表';









