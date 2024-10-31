Drop table if exists post;
Drop table if exists post_category;
Drop table if exists user_like;
Drop table if exists comment;
Drop table if exists category;
Drop table if exists user;
PRAGMA foreign_keys = ON;
insert into user (user_id, email, username, password, role, creation_date) values ('1', 'admin@admin.com','admin', 'admin', 'admin', '2021-01-01 00:00:00');
update user set role = 'admin' where username = "admin";
insert into category (category_id, name) values ('1', 'Shonen');
insert into category (category_id, name) values ('2', 'Naruto');

