use dailyhub;
drop table if exists profile;
drop table if exists habit;
drop table if exists month;
drop table if exists day;
drop table if exists token_item;
drop table if exists daily_commit;

create table if not exists profile (
    username char(255),
    password char(255),
    avatar MEDIUMBLOB,
    description char(255),
    habits char(255),
    primary key(username)
); 

create table if not exists habit (
	id char(255),
    time_quantum char(255),
	name char(255),
	icon char(255),
	file bool,
    color char(255),
	reminder_time char(255),
	encourage    char(255),
	important    bool,
	notification bool,
    recent_punch_time char(255),
    last_recent_punch_time char(255),
	total_punch integer,
	currc_punch integer,
	oncec_punch integer,
	create_at char(255),
    primary key(id)
);

create table if not exists month (
	id char(255),
    plan_punch char(255),
    actual_punch char(255),
    miss_punch char(255),
    days char(255),
    primary key(id)
);

create table if not exists day(
	id char(255),
    time char(255),
    log char(255),
    primary key(id)
);

create table if not exists token_item(
	username char(255),
    dh_token char(255),
    primary key(username)
);

create table if not exists daily_commit(
	id char(255),
    commit_time char(255),
    commit_content text,
    primary key(id)
);