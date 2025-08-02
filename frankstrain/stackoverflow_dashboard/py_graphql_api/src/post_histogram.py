from datetime import datetime

from src.types.posts import PostsQueryInterface


class PostHistogram:

    @staticmethod
    def get_response_shape():
        return {"errors": []}

    async def run(self, fields_requested, post_query: PostsQueryInterface, start_date: datetime, end_date: datetime):
        response = self.get_response_shape()
        for field in fields_requested:
            if hasattr(self, field):
                try:
                    response[field] = await getattr(self, field)(post_query=post_query, start_date=start_date, end_date=end_date)
                except Exception as e:
                    print(e)
                    response["errors"].append(f"failed to get {field}")

        return response

    @staticmethod
    async def created_overtime(post_query: PostsQueryInterface, start_date: datetime, end_date: datetime):
        data = await post_query.get_post_count_overtime(start_date=start_date, end_date=end_date)
        return [item.__dict__ for item in data]

    @staticmethod
    async def most_used_tags_overtime(post_query: PostsQueryInterface, start_date: datetime, end_date: datetime):
        data = await post_query.get_tags_count_withing_posts_overtime(start_date=start_date, end_date=end_date)
        return [item.__dict__ for item in data]
