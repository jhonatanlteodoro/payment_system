import enum
from dataclasses import dataclass
from datetime import datetime

from src.data.db_utils import db_data_factory, execute_query

class PostType(enum.IntEnum):
    QUESTION = 1
    ANSWER = 2

@dataclass
class Post:
    id: int
    score: int
    post_type_id: PostType
    creation_date: datetime
    view_count: int
    owner_user_id: int
    tags: str
    answer_count: int
    comment_count: int
    favorite_count: int

    def to_graphql_data(self) -> dict:
        data = self.__dict__
        data["id"] = str(data["id"])
        data["post_type_id"] = PostType(data["post_type_id"])
        return data

@dataclass
class PostOvertime:
    week: datetime
    count: int

@dataclass
class PostTagOvertime:
    week: datetime
    tag: str
    count: int

class PostsQuery:

    def __init__(self, db):
        self.db = db
        self.cachedQueries = dict()

    def _execute_query(self, query, query_args, multiple_objects=False, factory_item_class: any = Post):
        return execute_query(self.db, query, query_args, db_data_factory(factory_item_class), multiple_objects)

    def get_posts(self, page=1, page_limit=10):
        query = '''
            SELECT * FROM posts ORDER BY id DESC LIMIT %s OFFSET %s
        
        '''
        query_args = (page_limit, page*(page_limit if page_limit > 1 else 0))
        return self._execute_query(query, query_args, multiple_objects=True)


    def get_post(self, post_id):
        query = '''
            SELECT * FROM posts WHERE id = %s LIMIT 1
        '''
        return self._execute_query(query, (post_id,))

    async def get_post_count_overtime(self, start_date: datetime, end_date: datetime):
        string_start = start_date.strftime('%Y-%m-%d')
        string_end = end_date.strftime('%Y-%m-%d')
        query = '''
        with filteredByDate AS (
            SELECT DATE_TRUNC('week', creation_date) AS week
            FROM posts
                WHERE creation_date >= %s AND creation_date <= %s
        )
        
        SELECT week, count(*) FROM filteredByDate
            GROUP BY week ORDER BY week;
        '''
        return self._execute_query(query, (string_start, string_end), multiple_objects=True, factory_item_class=PostOvertime)

    async def get_tags_count_withing_posts_overtime(self, start_date: datetime, end_date: datetime):
        string_start = start_date.strftime('%Y-%m-%d')
        string_end = end_date.strftime('%Y-%m-%d')
        query = '''
        with filteredByDate AS (
            SELECT DATE_TRUNC('week', creation_date) AS week, tags
            FROM posts
                WHERE creation_date >= %s AND creation_date <= %s
                AND tags != ''
                AND LENGTH(TRIM(tags)) > 0
            ),
            tagsWithinDate AS (
                SELECT week, TRIM(unnest(string_to_array(tags, '|'))) AS tag
                FROM filteredByDate
            )
        SELECT week, tag, count(*)
        FROM tagsWithinDate WHERE tag != ''
        GROUP BY week, tag
        ORDER BY week DESC;
        '''
        return self._execute_query(query, (string_start, string_end), multiple_objects=True, factory_item_class=PostTagOvertime)