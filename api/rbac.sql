-- RBAC 建库脚本（新项目一次性执行）
-- 权限全部走角色；无 account_type 旁路
-- 用法：mysql -uroot < api/rbac.sql

CREATE DATABASE IF NOT EXISTS `rbac` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `rbac`;

DROP TABLE IF EXISTS `role_page`;
DROP TABLE IF EXISTS `account_role`;
DROP TABLE IF EXISTS `site_message_inbox`;
DROP TABLE IF EXISTS `site_message`;
DROP TABLE IF EXISTS `operation_log`;
DROP TABLE IF EXISTS `login_log`;
DROP TABLE IF EXISTS `page`;
DROP TABLE IF EXISTS `role`;
DROP TABLE IF EXISTS `account`;

CREATE TABLE `account` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`user_name` varchar(50) NOT NULL DEFAULT '' COMMENT '账号姓名',
`mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '电话',
`account_number` varchar(50) NOT NULL DEFAULT '' COMMENT '登录账号',
`password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码',
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

CREATE TABLE `operation_log` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`module` tinyint(4) NOT NULL DEFAULT 0 COMMENT '1角色 2角色页面权限 3账号角色 4页面资源',
`action` tinyint(4) NOT NULL DEFAULT 0 COMMENT '1新增 2修改 3删除 4分配 5移除',
`target_type` varchar(32) NOT NULL DEFAULT '' COMMENT '目标类型 role/page/role_page/account_role',
`target_id` int(11) NOT NULL DEFAULT 0 COMMENT '目标ID',
`target_name` varchar(100) NOT NULL DEFAULT '' COMMENT '目标名称',
`content` varchar(500) NOT NULL DEFAULT '' COMMENT '变更说明',
`operator_id` int(11) NOT NULL DEFAULT 0 COMMENT '操作人ID',
`operator_name` varchar(50) NOT NULL DEFAULT '' COMMENT '操作人姓名',
`operator_account` varchar(50) NOT NULL DEFAULT '' COMMENT '操作人账号',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`),
KEY `idx_module` (`module`),
KEY `idx_action` (`action`),
KEY `idx_operator_account` (`operator_account`),
KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='操作日志';

CREATE TABLE `site_message` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`title` varchar(120) NOT NULL DEFAULT '' COMMENT '标题',
`content` varchar(2000) NOT NULL DEFAULT '' COMMENT '内容',
`msg_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '1系统广播 2业务通知 3指定用户',
`biz_type` varchar(32) NOT NULL DEFAULT '' COMMENT '业务子类型编码',
`sender_id` int(11) NOT NULL DEFAULT 0 COMMENT '发送人ID',
`sender_name` varchar(50) NOT NULL DEFAULT '' COMMENT '发送人姓名',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`),
KEY `idx_msg_type` (`msg_type`),
KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='站内信';

CREATE TABLE `site_message_inbox` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`message_id` int(11) NOT NULL DEFAULT 0 COMMENT '站内信ID',
`account_id` int(11) NOT NULL DEFAULT 0 COMMENT '接收人账号ID',
`is_read` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0未读 1已读',
`read_time` int(11) NOT NULL DEFAULT 0 COMMENT '阅读时间',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '入箱时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_message_account` (`message_id`,`account_id`),
KEY `idx_account_read` (`account_id`,`is_read`),
KEY `idx_message_id` (`message_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='站内信收件箱';

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
`icon` varchar(64) NOT NULL DEFAULT '' COMMENT '菜单图标 Element Plus Icons 名，如 Setting',
`url` varchar(255) NOT NULL DEFAULT '' COMMENT '前端路由',
`parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '父级id',
`page_level` int(11) NOT NULL DEFAULT 0 COMMENT '1一级菜单 2二级菜单 3按钮',
`page_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '1菜单 2按钮',
`button_code` varchar(50) NOT NULL DEFAULT '' COMMENT '按钮编码 query/create/update/delete',
`api_method` varchar(16) NOT NULL DEFAULT '' COMMENT 'HTTP方法 GET/POST/PUT/DELETE',
`api_path` varchar(255) NOT NULL DEFAULT '' COMMENT 'Gin FullPath，如 /api/admin/account/:id',
`sort` int(11) NOT NULL DEFAULT 0 COMMENT '展示排序，数字越大越靠前',
`create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
PRIMARY KEY (`id`),
KEY `idx_parent_id` (`parent_id`),
KEY `idx_api` (`api_method`, `api_path`),
KEY `idx_sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=1  DEFAULT CHARSET=utf8 COMMENT='页面菜单与按钮';

-- ========== 初始角色：系统管理员 ==========
INSERT INTO `role` (`id`,`role_name`,`create_time`)
VALUES (1, '系统管理员', UNIX_TIMESTAMP());

-- ========== 初始账号 admin / admin（password 为 MD5） ==========
INSERT INTO `account` (`id`,`user_name`,`mobile`,`account_number`,`password`,`status`,`last_login_time`,`opertion_id`,`create_time`,`update_time`)
VALUES (1, '管理员', '13800000000', 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

INSERT INTO `account_role` (`account_id`,`role_id`,`create_time`)
VALUES (1, 1, UNIX_TIMESTAMP());

-- ========== 初始菜单 + 默认 CRUD 按钮（api_path 需与 Gin FullPath 一致） ==========

INSERT INTO `page` (`id`,`page_name`,`icon`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`)
VALUES (1, '系统管理', 'Setting', '/system', 0, 1, 1, '', '', '', UNIX_TIMESTAMP());

INSERT INTO `page` (`id`,`page_name`,`icon`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
(2, '账号管理', 'User', '/system/account', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(3, '角色管理', 'Avatar', '/system/role', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(4, '页面管理', 'Menu', '/system/page', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(5, '登录日志', 'Key', '/system/login-log', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(6, '操作日志', 'Document', '/system/operation-log', 1, 2, 1, '', '', '', UNIX_TIMESTAMP()),
(7, '站内信', 'Bell', '/system/message', 1, 2, 1, '', '', '', UNIX_TIMESTAMP());

INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 2, 3, 2, 'query',  'GET',    '/api/admin/account',     UNIX_TIMESTAMP()),
('新增', '', 2, 3, 2, 'create', 'POST',   '/api/admin/account',     UNIX_TIMESTAMP()),
('修改', '', 2, 3, 2, 'update', 'PUT',    '/api/admin/account/:id', UNIX_TIMESTAMP()),
('删除', '', 2, 3, 2, 'delete', 'DELETE', '/api/admin/account/:id', UNIX_TIMESTAMP()),
('分配角色', '', 2, 3, 2, 'assign_role', 'POST',   '/api/admin/account/role', UNIX_TIMESTAMP()),
('移除角色', '', 2, 3, 2, 'remove_role', 'DELETE', '/api/admin/account/role', UNIX_TIMESTAMP());

INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 3, 3, 2, 'query',  'GET',    '/api/admin/role',              UNIX_TIMESTAMP()),
('新增', '', 3, 3, 2, 'create', 'POST',   '/api/admin/role',              UNIX_TIMESTAMP()),
('修改', '', 3, 3, 2, 'update', 'PUT',    '/api/admin/role/:id',          UNIX_TIMESTAMP()),
('删除', '', 3, 3, 2, 'delete', 'DELETE', '/api/admin/role/:id',          UNIX_TIMESTAMP()),
('全部角色', '', 3, 3, 2, 'query_all', 'GET', '/api/admin/role/all',        UNIX_TIMESTAMP()),
('查看角色页面', '', 3, 3, 2, 'role_page', 'GET', '/api/admin/role/page/:role_id', UNIX_TIMESTAMP()),
('分配角色页面', '', 3, 3, 2, 'role_page_assign', 'POST', '/api/admin/role/page', UNIX_TIMESTAMP());

INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 4, 3, 2, 'query',  'GET',    '/api/admin/page',     UNIX_TIMESTAMP()),
('新增', '', 4, 3, 2, 'create', 'POST',   '/api/admin/page',     UNIX_TIMESTAMP()),
('修改', '', 4, 3, 2, 'update', 'PUT',    '/api/admin/page/:id', UNIX_TIMESTAMP()),
('删除', '', 4, 3, 2, 'delete', 'DELETE', '/api/admin/page/:id', UNIX_TIMESTAMP()),
('全部页面', '', 4, 3, 2, 'query_all', 'GET', '/api/admin/page/all', UNIX_TIMESTAMP());

INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 5, 3, 2, 'query', 'GET', '/api/admin/login/log', UNIX_TIMESTAMP());

INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('查询', '', 6, 3, 2, 'query', 'GET', '/api/admin/operation/log', UNIX_TIMESTAMP());

INSERT INTO `page` (`page_name`,`url`,`parent_id`,`page_level`,`page_type`,`button_code`,`api_method`,`api_path`,`create_time`) VALUES
('发送记录', '', 7, 3, 2, 'query',  'GET',  '/api/admin/message', UNIX_TIMESTAMP()),
('发送消息', '', 7, 3, 2, 'create', 'POST', '/api/admin/message', UNIX_TIMESTAMP());

-- ========== 系统管理员拥有全部菜单与按钮 ==========
INSERT INTO `role_page` (`role_id`,`page_id`)
SELECT 1, `id` FROM `page`;
