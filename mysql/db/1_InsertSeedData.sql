INSERT INTO sample.users (name) VALUES ('user_1'), ('user_2'), ('user_3'), ('user_4'), ('user_5');
INSERT INTO sample.tweets (user_id, content) VALUES
(1, 'First tweet! by user_1'), (1, 'Second tweet! by user_1'), (1, 'Third tweet! by user_1'),
(2, 'First tweet! by user_2'), (2, 'Second tweet! by user_2'), (2, 'Third tweet! by user_2'),
(3, 'First tweet! by user_3'), (3, 'Second tweet! by user_3'), (3, 'Third tweet! by user_3'),
(4, 'First tweet! by user_4'), (4, 'Second tweet! by user_4'), (4, 'Third tweet! by user_4'),
(5, 'First tweet! by user_5'), (5, 'Second tweet! by user_5'), (5, 'Third tweet! by user_5');
INSERT INTO sample.follows (follower_id, followed_id) VALUES (1, 2), (1, 3), (2, 3), (2, 4), (4, 5), (5, 1), (5, 2), (5, 3), (5, 4);
