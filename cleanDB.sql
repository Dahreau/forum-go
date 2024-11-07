PRAGMA foreign_keys = ON;
update user set role = 'admin' where username = "admin";
insert into category (category_id, name) values ('1', 'Shonen');
insert into category (category_id, name) values ('2', 'Naruto');


INSERT INTO User_like (like_id, user_id, post_id, comment_id, isLiked) VALUES (2,1165671184,1307312793,0,true);
delete from User_like;