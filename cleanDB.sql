PRAGMA foreign_keys = ON;
update user set role = 'admin' where username = "admin";
insert into category (category_id, name) values ('1', 'Shonen');
insert into category (category_id, name) values ('2', 'Naruto');

