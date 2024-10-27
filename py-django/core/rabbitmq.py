from kombu import Connection


def create_rabbitmq_connection() -> Connection:
    return Connection("amqp://guest:guest@rabbitmq:5672//")