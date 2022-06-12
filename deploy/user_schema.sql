create schema wonderful_hand;

create table user
(
    uid       bigint auto_increment comment '用户id'           primary key,
    name      varchar(128)                        not null comment '用户名称',
    password  varchar(64)                         not null comment '加密后的密码',
    nick_name varchar(128)                        not null comment '用户昵称',
    creat_at  timestamp default CURRENT_TIMESTAMP not null comment '创建时间',
    update_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    constraint user_name_uindex
        unique (name)
);

