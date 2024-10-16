drop table if exists `community`;
create table `community` (
    `id` int(11) not null AUTO_increment,
    `community_id` int(10) unsigned not null,
    `community_name` varchar(128) collate utf8mb4_general_ci not null,
    `introduction` varchar(256) collate utf8mb4_general_ci not null ,
    `create_time` timestamp not null default current_timestamp,
    `update_time` timestamp not null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `idx_community_id` (`community_id`),
    unique key `idx_community_name` (`community_name`)

) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;


insert into `community` values ('1','1','Go','Golang','2016-11-01 08:10:10','2016-11-01 08:10:10');
insert into `community` values ('2','2','leetcode','Golang','2020-11-01 08:10:10','2020-11-01 08:10:10');
insert into `community` values ('3','3','CS:Go','Golang','2018-11-01 08:10:10','2018-11-01 08:10:10');
insert into `community` values ('4','4','LOL','Golang','2016-01-01 08:10:10','2016-01-01 08:10:10');




DROP TARLE IF EXISTS `post`;CREATE TABLE `post` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `post_id` bigint(20) NOT NULL COMMENT '帖子id',
                        `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
                        `content` varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
                        `author_id` bigint(20) NOT NULL COMMENT '作者的用户id',
                        `community_id` bigint(20) NOT NULL COMMENT '所属社区',
                        `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '帖子状态',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_post_id` (`post_id`),
                        KEY `idx_author_id` (`author_id`),
                        KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;