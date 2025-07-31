from abc import ABC, abstractmethod

class VotesQueryInterface(ABC):
    @abstractmethod
    def get_vote(self, vote_id):
        pass

    @abstractmethod
    def get_votes(self, page, page_limit):
        pass