--- dict {'Id': '15476160', 'PostId': '8440', 'VoteTypeId': '2', 'CreationDate': '2011-05-29T00:00:00.000'}
CREATE TABLE IF NOT EXISTS votes (
    id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    vote_type_id SMALLINT NOT NULL,
    creation_date TIMESTAMP NOT NULL
);

--- row {'Id': '164363', 'TagName': 'liferay-blade-cli', 'Count': '6', 'ExcerptPostId': '78218357', 'WikiPostId': '78218356'}
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    count INTEGER NOT NULL,
    excerpt_post_id BIGINT,
    wiki_post_id BIGINT
);

--- 19: {'Id': '3', 'Reputation': '15767', 'CreationDate': '2008-07-31T14:22:31.287',
-- 'DisplayName': 'Jarrod Dixon', 'LastAccessDate': '2023-04-29T21:17:53.580', 'WebsiteUrl': 'http://jarroddixon.com',
-- 'Location': 'Johnson City, TN, USA', 'AboutMe': '', 'Views': '30431', 'UpVotes': '7765', 'DownVotes': '91', 'AccountId': '3'}
CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    reputation INTEGER NOT NULL,
    views INTEGER NOT NULL,
    up_votes INTEGER NOT NULL,
    down_votes INTEGER NOT NULL
);

--- # {'Id': '38784', 'PostTypeId': '1', 'AcceptedAnswerId': '41285', 'CreationDate': '2008-09-02T03:49:17.920',
-- 'Score': '7', 'ViewCount': '3474', 'Body': "<p.......... p>\n", 'OwnerUserId': '4149',
-- 'LastEditorUserId': '4779472', 'LastEditDate': '2015-06-27T14:21:16.453', 'LastActivityDate': '2021-11-25T10:14:44.487',
-- 'Title': 'Visual Studio equivalent to Delphi bookmarks',
-- 'Tags': '|visual-studio|delphi|brief-bookmarks|', 'AnswerCount': '8',
-- 'CommentCount': '1', 'FavoriteCount': '0', 'ContentLicense': 'CC BY-SA 3.0'}
CREATE TABLE IF NOT EXISTS posts (
    id BIGINT NOT NULL,
    post_type_id BIGINT NOT NULL,
    creation_date TIMESTAMP NOT NULL,
    score INTEGER NOT NULL,
    view_count INTEGER NOT NULL,
    owner_user_id BIGINT NOT NULL,
    tags VARCHAR(200) NOT NULL,
    answer_count SMALLINT NOT NULL,
    comment_count SMALLINT NOT NULL,
    favorite_count SMALLINT NOT NULL
);

select count(*) from users; --22261600
select count(*) from tags; --65000
select count(*) from posts;
select count(*) from votes; --90537500