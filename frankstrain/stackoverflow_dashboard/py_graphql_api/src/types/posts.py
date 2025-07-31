from abc import ABC, abstractmethod


class PostsQueryInterface(ABC):

    @abstractmethod
    def get_post(self, post_id):
        pass

    @abstractmethod
    def get_posts(self, page, page_limit):
        pass
