from datetime import datetime, timedelta

from fastapi import Depends, APIRouter, Request
from ariadne import QueryType, make_executable_schema, InterfaceType, load_schema_from_path, EnumType
from ariadne.asgi import GraphQL

from src.data.posts import PostType
from src.data.votes import VoteType
from src.types.container import ContainerInterface
from src.types.votes import VotesQueryInterface
from src.types.posts import PostsQueryInterface
from src.container import get_container

query = QueryType()
type_defs = load_schema_from_path("./schema.graphql")

vote_type_enum = EnumType("VoteType", {
    "UP_VOTE": VoteType.UP_VOTE,
    "DOWN_VOTE": VoteType.DOWN_VOTE,
    "FAVORITE": VoteType.FAVORITE,
})

post_type_enum = EnumType("PostType", {
    "QUESTION": PostType.QUESTION,
    "ANSWER": PostType.ANSWER,
})

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
async def resolve_list_posts(_, info):
    container = info.context["container"]
    data = [post.to_graphql_data() for post in container.get_instance(PostsQueryInterface).get_posts()]
    return {"success": True, "error": "", "data": data}

@query.field("getPost")
async def resolve_get_post(_, info, id: str):
    container = info.context["container"]
    data = container.get_instance(PostsQueryInterface).get_post(int(id))
    return {"success": True, "error": "", "data": data.to_graphql_data()}

@query.field("listVotes")
async def resolve_list_votes(_, info):
    container = info.context["container"]
    data = [post.to_graphql_data() for post in container.get_instance(VotesQueryInterface).get_votes()]
    return {"success": True, "error": "", "data": data}

@query.field("getVote")
async def resolve_get_vote(_, info, id: str):
    container = info.context["container"]
    data = container.get_instance(VotesQueryInterface).get_vote(int(id))
    return {"success": True, "error": "", "data": data.to_graphql_data()}


def requested_fields(info):
    selections = info.field_nodes[0].selection_set.selections
    fields = {selection.name.value for selection in selections}
    return fields

@query.field("postsHistogram")
async def resolve_posts_overtime(_, info, start_date: str, end_date: str):
    container = info.context["container"]

    start = datetime.fromisoformat(start_date)
    end = datetime.fromisoformat(end_date)

    if start > end or end - start > timedelta(days=30*6):
        return {"success": False, "errors": ["Start date must be before end date and range maximum 6 months."]}

    data = {
        "success": False,
        "errors": [],
    }
    fields = requested_fields(info)
    for field in fields:
        if field == "created_overtime":
            try:
                data[field] = [item.__dict__ for item in await container.get_instance(PostsQueryInterface).get_post_count_overtime(start, end)]
                continue
            except Exception as e:
                print(e)
                data["errors"].append("failed to get post creation overtime")

        if field == "most_used_tags_overtime":
            try:
                data[field] = [item.__dict__ for item in await container.get_instance(PostsQueryInterface).get_tags_count_withing_posts_overtime(start, end)]
            except Exception as e:
                print(e)
                data["errors"].append("failed to get users tags overtime")

    return data

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


schema = make_executable_schema(type_defs, query, node_interface, vote_type_enum, post_type_enum)
graphql_app = GraphQL(
    schema,
    debug=True,
    context_value=get_context_value,
)
