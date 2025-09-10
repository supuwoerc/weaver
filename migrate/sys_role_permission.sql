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

INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 1);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 2);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 3);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 4);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 5);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 6);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 7);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 8);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 9);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 10);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 11);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 12);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 13);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 14);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 15);
INSERT INTO gin_web.sys_role_permission (role_id, permission_id) VALUES (1, 16);
