FROM rabbitmq:3.11.19-management

RUN apt-get update && apt-get install -y curl

RUN curl -OL https://github.com/noxdafox/rabbitmq-message-deduplication/releases/download/0.6.1/elixir-1.13.4.ez && \
mv elixir-1.13.4.ez plugins/

RUN curl -OL https://github.com/noxdafox/rabbitmq-message-deduplication/releases/download/0.6.1/rabbitmq_message_deduplication-0.6.1.ez && \
mv rabbitmq_message_deduplication-0.6.1.ez plugins/

RUN rabbitmq-plugins enable rabbitmq_message_deduplication