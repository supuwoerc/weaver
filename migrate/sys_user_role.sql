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

INSERT INTO gin_web.sys_user_role (user_id, role_id) VALUES (1, 1);
