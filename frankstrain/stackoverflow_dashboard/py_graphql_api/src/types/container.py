from abc import ABC, abstractmethod
from typing import TypeVar, Type, Callable

T = TypeVar('T')


class ContainerInterface(ABC):

    @abstractmethod
    def register_singleton(self, interface: Type[T], instance: T) -> None:
        pass

    @abstractmethod
    def register_factory(self, interface: Type[T], factory: Callable[[], T]) -> None:
        pass

    @abstractmethod
    def get_instance(self, interface: Type[T]) -> T:
        pass

    @abstractmethod
    def get_factory(self, interface: Type[T]) -> Callable[[], T]:
        pass