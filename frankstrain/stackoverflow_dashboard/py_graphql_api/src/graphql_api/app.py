from ariadne import QueryType, make_executable_schema, InterfaceType, load_schema_from_path
from ariadne.asgi import GraphQL
from fastapi import FastAPI, Request


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
    post_sample = {"id": "id-id-id", "title": "title tile", "created_at": "2020..."}
    posts_list_result = [post_sample]
    return {"success": True, "error": "", "data": posts_list_result}

@query.field("getPost")
def resolve_get_post(*args, **kwargs):
    print(kwargs)
    post_sample = {"resolve_type": "Post" , "id": "id-id-id", "title": "title tile", "created_at": "2020..."}
    posts_result = {"success": True, "error": "", "data": post_sample}
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