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

class PostsQuery:

    def __init__(self, db):
        self.db = db
        self.cachedQueries = dict()

    def _execute_query(self, query, query_args, multiple_objects=False):
        return execute_query(self.db, query, query_args, db_data_factory(Post), multiple_objects)

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
