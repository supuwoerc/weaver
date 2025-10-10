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
    status     tinyint(1)       not null comment '账户状态',
    created_at datetime(3)      not null comment '创建时间',
    updated_at datetime(3)      not null comment '更新时间',
    deleted_at bigint default 0 not null comment '删除标志',
    constraint uni_sys_user_email
        unique (email)
)
    comment '用户表';

create index idx_sys_user_deleted_at
    on sys_user (deleted_at);

INSERT INTO gin_web.sys_user (id, email, password, nickname, avatar, gender, about, birthday, status, created_at, updated_at, deleted_at) VALUES (1, '2293172479@qq.com', '$2a$10$2zfaNgU3vbvCheGwEVSXlOaQ7d9aQEKKD1EkHEs.rJ087.s/dUS.W', null, null, null, null, null, 2, '2025-08-31 14:21:49.236', '2025-08-31 14:21:49.236', 0);
