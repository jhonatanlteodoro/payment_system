from fastapi import Depends, APIRouter, Request
from ariadne import QueryType, make_executable_schema, InterfaceType, load_schema_from_path
from ariadne.asgi import GraphQL
from src.types.container import ContainerInterface
from src.types.votes import VotesQueryInterface
from src.types.posts import PostsQueryInterface
from src.container import get_container

query = QueryType()
type_defs = load_schema_from_path("./schema.graphql")


# Create interface resolver for Node
node_interface = InterfaceType("Node")
@node_interface.type_resolver
def resolve_node_type(obj, info, *args):
    # Look at the object and determine its type
    if info.path.prev.key in ("getPost", "listPosts"):
        return "Post"

    if info.path.prev.key in ("getVote", "listVotes"):
        return "Vote"

    # Alternative: if you include a type hint in your data
    if "resolve_type" in obj:
        return obj["resolve_type"]

    return None


@query.field("listPosts")
def resolve_list_posts(_, info):
    container = info.context["container"]
    data = [post.to_graphql_data() for post in container.get_instance(PostsQueryInterface).get_posts()]
    return {"success": True, "error": "", "data": data}

@query.field("getPost")
def resolve_get_post(_, info, id: str):
    container = info.context["container"]
    data = container.get_instance(PostsQueryInterface).get_post(int(id))
    return {"success": True, "error": "", "data": data.to_graphql_data()}

@query.field("listVotes")
def resolve_list_votes(_, info):
    container = info.context["container"]
    data = [post.to_graphql_data() for post in container.get_instance(VotesQueryInterface).get_votes()]
    return {"success": True, "error": "", "data": data}

@query.field("getVote")
def resolve_get_vote(_, info, id: str):
    container = info.context["container"]
    data = container.get_instance(VotesQueryInterface).get_vote(int(id))
    return {"success": True, "error": "", "data": data.to_graphql_data()}

def get_context_value(request: Request, _data) -> dict:
    return {
        "request": request,
        "container": request.scope["container"],
    }


router = APIRouter(
    prefix="/graphql",
)

@router.get("/")
@router.options("/")
async def handle_graphql_explorer(request: Request):
    return await graphql_app.handle_request(request)

@router.post("/")
async def handle_graphql_query(
    request: Request,
    container: ContainerInterface = Depends(get_container)
):
    request.scope["container"] = container
    return await graphql_app.handle_request(request)


schema = make_executable_schema(type_defs, query, node_interface)
graphql_app = GraphQL(
    schema,
    debug=True,
    context_value=get_context_value,
)
