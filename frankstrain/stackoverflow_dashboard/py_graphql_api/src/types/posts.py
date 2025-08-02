from abc import ABC, abstractmethod
from datetime import datetime


class PostsQueryInterface(ABC):

    @abstractmethod
    def get_post(self, post_id):
        pass

    @abstractmethod
    def get_posts(self, page, page_limit):
        pass

    @abstractmethod
    async def get_post_count_overtime(self, start_date: datetime, end_date: datetime):
        pass

    @abstractmethod
    async def get_tags_count_withing_posts_overtime(self, start_date: datetime, end_date: datetime):
        pass