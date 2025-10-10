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

