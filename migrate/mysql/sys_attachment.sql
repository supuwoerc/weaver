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

