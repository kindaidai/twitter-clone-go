INSERT INTO sample.users (name, email, password, created_at, updated_at) VALUES
('user_1', 'user_1@example.com', 'password', current_timestamp, current_timestamp),
('user_2', 'user_2@example.com', 'password', current_timestamp, current_timestamp),
('user_3', 'user_3@example.com', 'password', current_timestamp, current_timestamp),
('user_4', 'user_4@example.com', 'password', current_timestamp, current_timestamp),
('user_5', 'user_5@example.com', 'password', current_timestamp, current_timestamp);
INSERT INTO sample.tweets (user_id, content, created_at, updated_at) VALUES
(1, 'First tweet! by user_1', current_timestamp, current_timestamp),
(1, 'Second tweet! by user_1', current_timestamp, current_timestamp),
(1, 'Third tweet! by user_1', current_timestamp, current_timestamp),
(2, 'First tweet! by user_2', current_timestamp, current_timestamp),
(2, 'Second tweet! by user_2', current_timestamp, current_timestamp),
(2, 'Third tweet! by user_2', current_timestamp, current_timestamp),
(3, 'First tweet! by user_3', current_timestamp, current_timestamp),
(3, 'Second tweet! by user_3', current_timestamp, current_timestamp),
(3, 'Third tweet! by user_3', current_timestamp, current_timestamp),
(4, 'First tweet! by user_4', current_timestamp, current_timestamp),
(4, 'Second tweet! by user_4', current_timestamp, current_timestamp),
(4, 'Third tweet! by user_4', current_timestamp, current_timestamp),
(5, 'First tweet! by user_5', current_timestamp, current_timestamp),
(5, 'Second tweet! by user_5', current_timestamp, current_timestamp),
(5, 'Third tweet! by user_5', current_timestamp, current_timestamp);
INSERT INTO sample.follows (follower_id, followed_id, created_at, updated_at) VALUES
(1, 2, current_timestamp, current_timestamp),
(1, 3, current_timestamp, current_timestamp),
(2, 3, current_timestamp, current_timestamp),
(2, 4, current_timestamp, current_timestamp),
(4, 5, current_timestamp, current_timestamp),
(5, 1, current_timestamp, current_timestamp),
(5, 2, current_timestamp, current_timestamp),
(5, 3, current_timestamp, current_timestamp),
(5, 4, current_timestamp, current_timestamp);
