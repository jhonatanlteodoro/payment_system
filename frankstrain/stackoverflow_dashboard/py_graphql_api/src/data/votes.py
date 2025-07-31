import enum
from dataclasses import dataclass
from datetime import datetime
from src.data.db_utils import db_data_factory, execute_query


class VoteType(enum.IntEnum):
    UP_VOTE = 2
    DOWN_VOTE = 3
    FAVORITE = 5

@dataclass
class Vote:
    id: int
    post_id: int
    vote_type_id: VoteType
    creation_date: datetime

    def to_graphql_data(self) -> dict:
        data = self.__dict__
        data["id"] = str(data["id"])
        data["vote_type_id"] = VoteType(data["vote_type_id"]).value
        return data


class VotesQuery:

    def __init__(self, db):
        self.db = db

    def _execute_query(self, query, query_args, multiple_objects=False):
        return execute_query(self.db, query, query_args, db_data_factory(Vote), multiple_objects)

    def get_votes(self, page=1, page_limit=10):
        query = '''
            SELECT * FROM votes ORDER BY creation_date DESC LIMIT %s OFFSET %s
        '''

        query_args = (page_limit, page*(page_limit if page_limit > 1 else 0))
        return self._execute_query(query, query_args, multiple_objects=True)

    def get_vote(self, vote_id):
        query = '''
            SELECT * FROM votes WHERE id = %s limit 1
        '''
        query_args = (vote_id,)
        return self._execute_query(query, query_args)
