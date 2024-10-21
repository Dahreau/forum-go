-- Insertion des utilisateurs
INSERT INTO User (user_id, email, username, password, role, creation_date, session_id, session_expire)
VALUES ('1', 'john.doe@example.com', 'JohnDoe', 'hashed_password_123', 'user', '2024-10-01 12:00:00', NULL, NULL);

INSERT INTO User (user_id, email, username, password, role, creation_date, session_id, session_expire)
VALUES ('2', 'jane.smith@example.com', 'JaneSmith', 'hashed_password_456', 'user', '2024-10-01 12:05:00', NULL, NULL);

-- Insertion des posts
INSERT INTO Post (post_id, title, content, user_id, creation_date, update_date)
VALUES ('101', 'First Post', 'This is the content of the first post.', '1', '2024-10-02 10:00:00', NULL);

INSERT INTO Post (post_id, title, content, user_id, creation_date, update_date)
VALUES ('102', 'Second Post', 'Here is some content for the second post.', '2', '2024-10-02 11:00:00', NULL);

-- Insertion des commentaires
INSERT INTO Comment (comment_id, content, creation_date, update_date, user_id, post_id)
VALUES ('1001', 'This is a comment on the first post.', '2024-10-02 12:00:00', NULL, '2', '101');

INSERT INTO Comment (comment_id, content, creation_date, update_date, user_id, post_id)
VALUES ('1002', 'Another comment on the first post.', '2024-10-02 12:10:00', NULL, '1', '101');

-- Insertion des catégories
INSERT INTO Category (category_id, name)
VALUES ('1', 'Technology');

INSERT INTO Category (category_id, name)
VALUES ('2', 'Science');

-- Insertion des likes
INSERT INTO User_Like (like_id, isLiked, user_id, post_id, comment_id)
VALUES ('501', 1, '1', '102', NULL);  -- Like on Post

INSERT INTO User_Like (like_id, isLiked, user_id, post_id, comment_id)
VALUES ('502', 1, '2', NULL, '1001'); -- Like on Comment

-- Insertion des relations entre posts et catégories
INSERT INTO Post_Category (post_id, category_id)
VALUES ('101', '1');  -- First post belongs to Technology

INSERT INTO Post_Category (post_id, category_id)
VALUES ('102', '2');  -- Second post belongs to Science