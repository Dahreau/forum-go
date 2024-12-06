PRAGMA foreign_keys = ON;
update user set role = 'admin' where username = "admin";
insert into category (category_id, name) values ('1', 'Shonen');
insert into category (category_id, name) values ('2', 'Naruto');


INSERT INTO User_like (like_id, user_id, post_id, comment_id, isLiked) VALUES (2,1165671184,1307312793,0,true);
delete from User_like;



CREATE TABLE IF NOT EXISTS Request(
    request_id CHAR(32) PRIMARY KEY,
    user_id CHAR(32) NOT NULL,
    status VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    creation_date DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Activity(
    activity_id   CHAR(32) PRIMARY KEY,               -- Identifiant unique pour l'activité
    user_id       CHAR(32) NOT NULL,                  -- ID de l'utilisateur destinataire de l'activité
    action_user_id CHAR(32) NOT NULL,                 -- ID de l'utilisateur ayant généré l'activité
    action_type   VARCHAR(50) NOT NULL,                  -- Type d'action (like, comment, etc.)
    post_id       CHAR(32) NOT NULL,                           -- ID du post concerné (ne peut pas être NULL)
    comment_id    CHAR(32),                           -- ID du commentaire concerné (peut être NULL)
    creation_date DATETIME NOT NULL, -- Date de création de l'activité
    details       TEXT,                           -- Informations supplémentaires
    is_read       BOOLEAN,          -- Indicateur si l'activité a été lue
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,      -- FK vers la table User
    FOREIGN KEY (action_user_id) REFERENCES User(user_id) ON DELETE CASCADE, -- FK vers la table User
    FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE     -- FK vers la table Post
    
);

update user set provider = 'local';