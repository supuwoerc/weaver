create table sys_attachment
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(255)     not null comment '附件名称',
    type       tinyint          not null comment '文件类型',
    size       bigint default 0 not null comment '文件大小',
    hash       varchar(72)      not null comment '文件摘要',
    path       longtext         not null comment '文件路径',
    creator_id bigint unsigned  null comment '创建者ID,null为系统创建',
    created_at datetime(3)      not null comment '创建时间',
    updated_at datetime(3)      not null comment '更新时间',
    deleted_at bigint default 0 not null comment '删除标志'
)
    comment '附件表';

create index idx_sys_attachment_deleted_at
    on sys_attachment (deleted_at);

create index idx_sys_attachment_hash
    on sys_attachment (hash);

create index idx_sys_attachment_uid_name_type
    on sys_attachment (creator_id, name, type);

create table sys_department
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(20)      not null comment '部门名称',
    parent_id  bigint unsigned  null comment '父部门ID',
    ancestors  varchar(255)     null comment '祖先部门ID逗号拼接的字符串',
    creator_id bigint unsigned  not null comment '创建人ID',
    updater_id bigint unsigned  not null comment '更新人ID',
    created_at datetime(3)      not null comment '创建时间',
    updated_at datetime(3)      not null comment '更新时间',
    deleted_at bigint default 0 not null comment '删除标志',
    constraint idx_sys_department_name_deleted_at
        unique (name, deleted_at)
);

create index idx_sys_department_deleted_at
    on sys_department (deleted_at);

create index idx_sys_department_parent_id
    on sys_department (parent_id);

create table sys_department_leader
(
    department_id bigint unsigned not null comment '部门ID',
    user_id       bigint unsigned not null comment '用户ID',
    primary key (department_id, user_id)
);

create index idx_sys_department_leader_department_id
    on sys_department_leader (department_id);

create index idx_sys_department_leader_user_id
    on sys_department_leader (user_id);

create table sys_permission
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(20)      not null comment '权限名',
    resource   varchar(255)     not null comment '资源名',
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

create table sys_role
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    name       varchar(20)      not null comment '角色名',
    creator_id bigint unsigned  not null comment '创建人ID',
    updater_id bigint unsigned  not null comment '更新人ID',
    created_at datetime(3)      not null comment '创建时间',
    updated_at datetime(3)      not null comment '更新时间',
    deleted_at bigint default 0 not null comment '删除时间',
    constraint uni_sys_role_name_deleted_at
        unique (name, deleted_at)
)
    comment '角色表';

create index idx_sys_role_deleted_at
    on sys_role (deleted_at);

create table sys_role_permission
(
    role_id       bigint unsigned not null comment '角色ID',
    permission_id bigint unsigned not null comment '权限ID',
    primary key (role_id, permission_id)
)
    comment '角色-权限中间表';

create index idx_sys_role_permission_permission_id
    on sys_role_permission (permission_id);

create index idx_sys_role_permission_role_id
    on sys_role_permission (role_id);

create table sys_user
(
    id         bigint unsigned auto_increment comment '主键ID'
        primary key,
    email      varchar(50)      not null comment '邮箱',
    password   varchar(60)      not null comment '密码',
    nickname   varchar(20)      null comment '昵称',
    avatar     varchar(255)     null comment '头像文件URL',
    gender     tinyint unsigned null comment '性别',
    about      varchar(100)     null comment '关于',
    birthday   datetime(3)      null comment '生日',
    created_at datetime(3)      not null comment '创建时间',
    updated_at datetime(3)      not null comment '更新时间',
    deleted_at bigint default 0 not null comment '删除标志',
    constraint uni_sys_user_email
        unique (email)
)
    comment '用户表';

create index idx_sys_user_deleted_at
    on sys_user (deleted_at);

create table sys_user_department
(
    department_id bigint unsigned not null comment '部门ID',
    user_id       bigint unsigned not null comment '用户ID',
    primary key (department_id, user_id)
);

create index idx_sys_user_department_department_id
    on sys_user_department (department_id);

create index idx_sys_user_department_user_id
    on sys_user_department (user_id);

create table sys_user_role
(
    user_id bigint unsigned not null comment '用户ID',
    role_id bigint unsigned not null comment '角色ID',
    primary key (user_id, role_id)
)
    comment '用户-角色中间表';

create index idx_sys_user_role_role_id
    on sys_user_role (role_id);

create index idx_sys_user_role_user_id
    on sys_user_role (user_id);