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

INSERT INTO gin_web.sys_role (id, name, creator_id, updater_id, created_at, updated_at, deleted_at) VALUES (1, '测试角色', 1, 1, '2025-09-05 07:48:29.000', '2025-09-05 07:48:31.000', 0);
