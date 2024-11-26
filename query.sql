CREATE TABLE IF NOT EXISTS User (
  user_id CHAR(32) PRIMARY KEY,
  email VARCHAR(100) NOT NULL,
  username VARCHAR(50) NOT NULL,
  password VARCHAR(255) NOT NULL,
  -- Hashed password
  role VARCHAR(50) NOT NULL,
  creation_date DATETIME NOT NULL,
  session_id CHAR(32),
  session_expire DATETIME
);
CREATE TABLE IF NOT EXISTS Post (
  post_id CHAR(32) PRIMARY KEY,
  title VARCHAR(50) NOT NULL,
  content TEXT NOT NULL,
  -- Content of the post
  user_id CHAR(32) NOT NULL,
  creation_date DATETIME NOT NULL,
  update_date DATETIME,
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Comment (
  comment_id CHAR(32) PRIMARY KEY,
  content TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  update_date DATETIME,
  user_id CHAR(32) NOT NULL,
  post_id CHAR(32) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
  FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Category (
  category_id CHAR(32) PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS User_Like (
  like_id CHAR(32) PRIMARY KEY,
  isLiked BOOLEAN NOT NULL,
  user_id CHAR(32) NOT NULL,
  post_id CHAR(32),
  -- If the user likes a post
  comment_id CHAR(32),
  -- If the user likes a comment
  FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE
  -- FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE,
  -- FOREIGN KEY (comment_id) REFERENCES Comment(comment_id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS Post_Category (
  post_id CHAR(32),
  category_id CHAR(32),
  PRIMARY KEY (post_id, category_id),
  FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES Category(category_id) ON DELETE CASCADE
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


CREATE TABLE IF NOT EXISTS Request(
    request_id CHAR(32) PRIMARY KEY,
    user_id CHAR(32) NOT NULL,
    status VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    creation_date DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE
);