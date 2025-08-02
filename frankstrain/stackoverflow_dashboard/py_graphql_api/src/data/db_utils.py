from collections.abc import Callable
from contextlib import asynccontextmanager
from typing import Any, Union, List, Tuple, TypeVar

from fastapi import FastAPI
from psycopg import Cursor, AsyncCursor
from psycopg_pool import AsyncConnectionPool

DataClassType = TypeVar('DataClassType', bound="DataClass")

def db_data_factory(data_class: DataClassType) -> Callable[[AsyncCursor], Callable[[Any], Any]]:
    def handle_cursor(cursor):
        columns = [info.name for info in cursor.description]

        def make(row):
            return data_class(**{columns[i]: row[i] for i in range(len(columns))})

        return make

    return handle_cursor


async def execute_query(db: AsyncConnectionPool, query: str, query_args: Tuple[str], factory: Callable[[Cursor], Any], multiple_objects=False) -> Union[List[Any], Any, None]:
    async with db.connection() as conn:
        async with conn.cursor(row_factory=factory) as curs:
            await curs.execute(query, query_args)
            if multiple_objects:
                return await curs.fetchall()
            return await curs.fetchone()


db_url = "postgresql://secret_user:secret_password@localhost:5432/stackoverflow"
connection_pool = AsyncConnectionPool(db_url, num_workers=2, open=False)

@asynccontextmanager
async def db_connection_as_lifespan(app: FastAPI):
    await connection_pool.open()
    yield
    await connection_pool.close()