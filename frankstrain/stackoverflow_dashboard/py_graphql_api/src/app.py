from ariadne import QueryType, make_executable_schema, InterfaceType, load_schema_from_path
from ariadne.asgi import GraphQL
from fastapi import FastAPI, Request
from psycopg_pool import ConnectionPool

from src.data.posts import PostsQuery

db_url = "postgresql://secret_user:secret_password@localhost:5432/stackoverflow"
connection_pool = ConnectionPool(db_url, max_size=5, open=True, num_workers=5)
connection_pool.wait()
print("pool ready")
pp = PostsQuery(connection_pool)

type_defs = load_schema_from_path("./schema.graphql")
query = QueryType()

# Create interface resolver for Node
node_interface = InterfaceType("Node")
@node_interface.type_resolver
def resolve_node_type(obj, *args, **kwargs):
    # Look at the object and determine its type
    if args[0].path.prev.key in ("getPost", "listPosts"):
        return "Post"

    # Alternative: if you include a type hint in your data
    if "resolve_type" in obj:
        return obj["resolve_type"]

    return None


@query.field("listPosts")
def resolve_hello(*_):
    data = [post.to_graphql_data() for post in pp.resolve_list_posts()]
    return {"success": True, "error": "", "data": data}

@query.field("getPost")
def resolve_get_post(*args, **kwargs):
    data = pp.resolve_post(int(kwargs["id"]))
    posts_result = {"success": True, "error": "", "data": data.to_graphql_data()}
    return posts_result


schema = make_executable_schema(type_defs, query, node_interface)

graphql_app = GraphQL(
    schema,
    debug=True,
)

app = FastAPI()


@app.get("/graphql/")
@app.options("/graphql/")
async def handle_graphql_explorer(request: Request):
    return await graphql_app.handle_request(request)

@app.post("/graphql/")
async def handle_graphql_query(
    request: Request,
):
    return await graphql_app.handle_request(request)