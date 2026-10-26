-- RBAC 建库脚本（新项目一次性执行）

CREATE TABLE `account` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '账号姓名',
`mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '电话',
`account_number` varchar(50) NOT NULL DEFAULT '' COMMENT '登录账号',
`password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码',
`account_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '账户类型 1超级管理员',
`status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '账号状态 1启用 2禁用',
`last_login_time` int(11) NOT NULL DEFAULT 0 COMMENT '最后一次登录时间',
`opertion_id` int(11) NOT NULL DEFAULT 0 COMMENT '操作人',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
`update_time` int(11)  NOT NULL DEFAULT 0 COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `account_number` (`account_number`),
UNIQUE KEY `mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='后台用户表';

CREATE TABLE `login_log` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`account_id` int(11) NOT NULL DEFAULT 0 COMMENT '账号ID，未知为0',
`account_number` varchar(50) NOT NULL DEFAULT '' COMMENT '登录账号',
`login_ip` varchar(64) NOT NULL DEFAULT '' COMMENT '登录IP',
`user_agent` varchar(255) NOT NULL DEFAULT '' COMMENT 'UA',
`status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '1成功 2失败',
`message` varchar(100) NOT NULL DEFAULT '' COMMENT '说明',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`),
KEY `idx_account_id` (`account_id`),
KEY `idx_account_number` (`account_number`),
KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='登录日志';

CREATE TABLE `role` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名称',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户角色表';

CREATE TABLE `account_role` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`account_id` int(11) NOT NULL DEFAULT 0 COMMENT '账号ID',
`role_id` int(11) NOT NULL DEFAULT 0 COMMENT '角色ID',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`),
KEY `idx_account_id` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='用户角色关联表';

CREATE TABLE `role_page` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`role_id` int(11) NOT NULL DEFAULT 0 COMMENT '角色ID',
`page_id` int(11) NOT NULL DEFAULT 0 COMMENT '页面或按钮ID',
PRIMARY KEY (`id`),
KEY `idx_role_id` (`role_id`),
KEY `idx_page_id` (`page_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='角色权限关联表';

CREATE TABLE `page` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`page_name` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单/按钮名称',
`url` varchar(255) NOT NULL DEFAULT '' COMMENT '前端路由',
`parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '父级id',
`page_level` int(11) NOT NULL DEFAULT 0 COMMENT '1一级菜单 2二级菜单 3按钮',
`page_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '1菜单 2按钮',
`button_code` varchar(50) NOT NULL DEFAULT '' COMMENT '按钮编码 query/create/update/delete',
`api_method` varchar(16) NOT NULL DEFAULT '' COMMENT 'HTTP方法 GET/POST/PUT/DELETE',
`api_path` varchar(255) NOT NULL DEFAULT '' COMMENT 'Gin FullPath，如 /api/admin/account/:id',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`),
KEY `idx_parent_id` (`parent_id`),
KEY `idx_api` (`api_method`, `api_path`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='页面菜单与按钮';

-- ========== 初始超级管理员（账号 admin / 密码 admin，password 为 MD5） ==========
INSERT INTO `account` (`id`,`user_name`,`mobile`,`account_number`,`password`,`account_type`,`status`,`last_login_time`,`opertion_id`,`create_time`,`update_time`)
VALUES (1, '管理员', '13800000000', 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, 1, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- ========== 初始菜单 + 默认 CRUD 按钮（api_path 需与 Gin FullPath 一致） ==========

-- 一级菜单：系统管理
INSERT INTO `page` (`id`,`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`)
VALUES (1, '系统管理', '/system', 0, 1, 1, '', '', '', UNIX_TIMESTAMP());

-- 二级菜单
INSERT INTO `page` (`id`,`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
(2, '账号管理', '/system/account', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(3, '角色管理', '/system/role', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(4, '页面管理', '/system/page', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(5, '登录日志', '/system/login-log', 1, 2, 1, '', '', '', UNIX_TIMESTAMP());

-- 账号管理 CRUD
INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 2, 3, 2, 'query',  'GET',    '/api/admin/account',     UNIX_TIMESTAMP()),
('新增', '', 2, 3, 2, 'create', 'POST',   '/api/admin/account',     UNIX_TIMESTAMP()),
('修改', '', 2, 3, 2, 'update', 'PUT',    '/api/admin/account/:id', UNIX_TIMESTAMP()),
('删除', '', 2, 3, 2, 'delete', 'DELETE', '/api/admin/account/:id', UNIX_TIMESTAMP()),
('分配角色', '', 2, 3, 2, 'assign_role', 'POST',   '/api/admin/account/role', UNIX_TIMESTAMP()),
('移除角色', '', 2, 3, 2, 'remove_role', 'DELETE', '/api/admin/account/role', UNIX_TIMESTAMP());

-- 角色管理 CRUD + 角色页面分配
INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 3, 3, 2, 'query',  'GET',    '/api/admin/role',              UNIX_TIMESTAMP()),
('新增', '', 3, 3, 2, 'create', 'POST',   '/api/admin/role',              UNIX_TIMESTAMP()),
('修改', '', 3, 3, 2, 'update', 'PUT',    '/api/admin/role/:id',          UNIX_TIMESTAMP()),
('删除', '', 3, 3, 2, 'delete', 'DELETE', '/api/admin/role/:id',          UNIX_TIMESTAMP()),
('全部角色', '', 3, 3, 2, 'query_all', 'GET', '/api/admin/role/all',        UNIX_TIMESTAMP()),
('查看角色页面', '', 3, 3, 2, 'role_page', 'GET', '/api/admin/role/page/:role_id', UNIX_TIMESTAMP()),
('分配角色页面', '', 3, 3, 2, 'role_page_assign', 'POST', '/api/admin/role/page', UNIX_TIMESTAMP());

-- 页面管理 CRUD
INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 4, 3, 2, 'query',  'GET',    '/api/admin/page',     UNIX_TIMESTAMP()),
('新增', '', 4, 3, 2, 'create', 'POST',   '/api/admin/page',     UNIX_TIMESTAMP()),
('修改', '', 4, 3, 2, 'update', 'PUT',    '/api/admin/page/:id', UNIX_TIMESTAMP()),
('删除', '', 4, 3, 2, 'delete', 'DELETE', '/api/admin/page/:id', UNIX_TIMESTAMP()),
('全部页面', '', 4, 3, 2, 'query_all', 'GET', '/api/admin/page/all', UNIX_TIMESTAMP());

-- 登录日志
INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 5, 3, 2, 'query', 'GET', '/api/admin/login/log', UNIX_TIMESTAMP());
