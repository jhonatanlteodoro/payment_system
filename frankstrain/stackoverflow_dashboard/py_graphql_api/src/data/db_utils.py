from collections.abc import Callable
from typing import Any, Union, List, Tuple, TypeVar

from psycopg import Cursor
from psycopg_pool import ConnectionPool

DataClassType = TypeVar('DataClassType', bound="DataClass")

def db_data_factory(data_class: DataClassType) -> Callable[[Cursor], Callable[[Any], Any]]:
    def handle_cursor(cursor):
        columns = [info.name for info in cursor.description]

        def make(row):
            return data_class(**{columns[i]: row[i] for i in range(len(columns))})

        return make

    return handle_cursor


def execute_query(db: ConnectionPool, query: str, query_args: Tuple[str], factory: Callable[[Cursor], Any], multiple_objects=False) -> Union[List[Any], Any, None]:
    with db.connection() as conn:
        with conn.cursor(row_factory=factory) as curs:
            curs.execute(query, query_args)
            if multiple_objects:
                return curs.fetchall()
            return curs.fetchone()
