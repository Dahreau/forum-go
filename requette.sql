CREATE TABLE User (
  user_id CHAR(32) PRIMARY KEY,
  email VARCHAR(100) NOT NULL,
  username VARCHAR(50) NOT NULL,
  password VARCHAR(255) NOT NULL, -- Hashed password
  prefix VARCHAR(50) NOT NULL,
  creation_date DATETIME NOT NULL,
  session_id CHAR(32),
  session_expire DATETIME
);

CREATE TABLE Post (
  post_id CHAR(32) PRIMARY KEY,
  title VARCHAR(50) NOT NULL,
  content TEXT NOT NULL, -- Content of the post
  user_id CHAR(32) NOT NULL,
  creation_date DATETIME NOT NULL,
  update_date DATETIME,
  FOREIGN KEY (user_id) REFERENCES User(user_id)
);

CREATE TABLE Comment (
  comment_id CHAR(32) PRIMARY KEY,
  content TEXT NOT NULL,
  creation_date DATETIME NOT NULL,
  update_date DATETIME,
  user_id CHAR(32) NOT NULL,
  post_id CHAR(32) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES User(user_id),
  FOREIGN KEY (post_id) REFERENCES Post(post_id)
);

CREATE TABLE Category (
  category_id CHAR(32) PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE Like (
  like_id CHAR(32) PRIMARY KEY,
  isLiked BOOLEAN NOT NULL,
  user_id CHAR(32) NOT NULL,
  post_id CHAR(32), -- If the user likes a post
  comment_id CHAR(32), -- If the user likes a comment
  FOREIGN KEY (user_id) REFERENCES User(user_id),
  FOREIGN KEY (post_id) REFERENCES Post(post_id),
  FOREIGN KEY (comment_id) REFERENCES Comment(comment_id)
);

CREATE TABLE Post_Category (
  post_id CHAR(32),
  category_id CHAR(32),
  PRIMARY KEY (post_id, category_id),
  FOREIGN KEY (post_id) REFERENCES Post(post_id),
  FOREIGN KEY (category_id) REFERENCES Category(category_id)
);