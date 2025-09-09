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

