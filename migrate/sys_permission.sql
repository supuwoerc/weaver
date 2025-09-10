create table sys_permission
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(20)      not null comment '权限名',
    resource   varchar(255)     not null comment '资源名',
    type       tinyint(1)       not null comment '资源类型',
    creator_id bigint unsigned  not null comment '创建人ID',
    updater_id bigint unsigned  not null comment '更新人ID',
    created_at datetime(3)      not null comment '创建时间',
    updated_at datetime(3)      not null comment '更新时间',
    deleted_at bigint default 0 not null comment '删除时间',
    constraint uni_sys_permission_name_deleted_at
        unique (name, deleted_at),
    constraint uni_sys_permission_resource_deleted_at
        unique (resource, deleted_at)
)
    comment '权限表';

create index idx_sys_permission_deleted_at
    on sys_permission (deleted_at);

INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (1, '系统设置-用户管理', '/setting/user', 1, 1, 1, '2025-09-05 07:48:01.000', '2025-09-05 07:48:03.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (2, '系统设置-部门管理', '/setting/department', 1, 1, 1, '2025-09-09 21:00:40.000', '2025-09-09 21:00:42.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (3, '系统设置-角色管理', '/setting/role', 1, 1, 1, '2025-09-09 21:01:08.000', '2025-09-09 21:01:09.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (4, '系统设置-权限管理', '/setting/permission', 1, 1, 1, '2025-09-09 21:01:31.000', '2025-09-09 21:01:33.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (5, '系统设置-用户管理-菜单', 'router.setting.user', 2, 1, 1, '2025-09-09 21:04:58.000', '2025-09-09 21:05:01.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (6, '系统设置-部门管理-菜单', 'router.setting.department', 2, 1, 1, '2025-09-09 21:07:16.000', '2025-09-09 21:07:18.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (7, '系统设置-角色管理-菜单', 'router.setting.role', 2, 1, 1, '2025-09-09 21:09:23.000', '2025-09-09 21:09:25.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (8, '系统设置-权限管理-菜单', 'router.setting.permission', 2, 1, 1, '2025-09-09 21:10:03.000', '2025-09-09 21:10:05.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (9, '系统设置-用户管理-列表接口', '/api/v1/user/list', 4, 1, 1, '2025-09-09 21:13:32.000', '2025-09-09 21:13:33.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (10, '系统设置-部门管理-创建接口', '/api/v1/department/create', 4, 1, 1, '2025-09-09 21:14:26.000', '2025-09-09 21:14:27.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (11, '系统设置-权限管理-创建接口', '/api/v1/permission/create', 4, 1, 1, '2025-09-09 21:16:38.000', '2025-09-09 21:16:46.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (12, '系统设置-权限管理-列表接口', '/api/v1/permission/list', 4, 1, 1, '2025-09-09 21:16:40.000', '2025-09-09 21:16:48.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (13, '系统设置-权限管理-详情接口', '/api/v1/permission/detail', 4, 1, 1, '2025-09-09 21:16:41.000', '2025-09-09 21:16:49.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (14, '系统设置-权限管理-更新接口', '/api/v1/permission/update', 4, 1, 1, '2025-09-09 21:16:43.000', '2025-09-09 21:16:50.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (15, '系统设置-权限管理-删除接口', '/api/v1/permission/delete', 4, 1, 1, '2025-09-09 21:16:44.000', '2025-09-09 21:16:52.000', 0);
INSERT INTO gin_web.sys_permission (id, name, resource, type, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (16, '系统设置-角色管理-列表接口', '/api/v1/role/list', 4, 1, 1, '2025-09-10 08:29:46.000', '2025-09-10 08:29:48.000', 0);
