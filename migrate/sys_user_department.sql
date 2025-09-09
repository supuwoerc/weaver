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

