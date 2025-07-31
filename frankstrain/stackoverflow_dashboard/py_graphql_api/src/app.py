from fastapi import FastAPI
from src.container import get_container
from src.data.db_utils import db_connection_as_lifespan, connection_pool
from src.types.votes import VotesQueryInterface
from src.types.posts import PostsQueryInterface

def start_app():
    from src.data.posts import PostsQuery
    from src.data.votes import VotesQuery

    from src.resolvers import router
    print("here")
    app = FastAPI(
        title="StackOverflow Dashboard",
        description="Dashboard for StackOverflow",
        version="1.0.0",
        lifespan=db_connection_as_lifespan,
    )

    container = get_container()
    container.register_singleton(PostsQueryInterface, PostsQuery(connection_pool))
    container.register_singleton(VotesQueryInterface, VotesQuery(connection_pool))

    app.include_router(router)

    return app


app = start_app()