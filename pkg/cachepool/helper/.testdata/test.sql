create schema awesome;

use awesome;

create table t
(
    yee varchar(1024) default '' not null,
    bar int                      null,
    foo datetime                 null
);

insert into r (yee, bar) values ('Hello', 1);
insert into r (yee, bar, foo) values ('Hi', 2, NOW());