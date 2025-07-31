from typing import TypeVar, Type, Dict, Callable, Any

T = TypeVar('T')

class Container:

    def __init__(self):
        self._instances: Dict[Type, Any] = {}
        self._factories: Dict[Type, Callable[[], Type]] = {}

    def register_singleton(self, interface: Type[T], instance: T) -> None:
        self._instances[interface] = instance

    def register_factory(self, interface: Type[T], factory: Callable[[], T]) -> None:
        self._factories[interface] = factory

    def get_instance(self, interface: Type[T]) -> T:
        instance = self._instances.get(interface)
        if instance is None:
            raise Exception(f"Interface {interface} not found")
        return instance

    def get_factory(self, interface: Type[T]) -> Callable[[], T]:
        factory = self._factories.get(interface)
        if factory is None:
            raise Exception(f"Factory {interface} not found")
        return factory


container = Container()

def get_container():
    return container