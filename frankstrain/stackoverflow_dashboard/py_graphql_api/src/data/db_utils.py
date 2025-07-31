from collections.abc import Callable
from contextlib import asynccontextmanager
from typing import Any, Union, List, Tuple, TypeVar

from fastapi import FastAPI
from psycopg import Cursor
from psycopg_pool import ConnectionPool, AsyncConnectionPool

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


db_url = "postgresql://secret_user:secret_password@localhost:5432/stackoverflow"
connection_pool = ConnectionPool(db_url, num_workers=2, open=False)

@asynccontextmanager
async def db_connection_as_lifespan(app: FastAPI):
    connection_pool.open()
    yield
    connection_pool.close()