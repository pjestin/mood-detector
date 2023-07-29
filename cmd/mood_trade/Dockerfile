FROM docker.io/library/golang:1.20.0

ADD . /app
WORKDIR /app

RUN apt-get update && apt-get -y upgrade
RUN apt-get -y install cron

RUN go build ./cmd/mood_trade

ADD cmd/mood_trade/crontab /etc/cron.d/mood-detector-cron
RUN chmod 0644 /etc/cron.d/mood-detector-cron
RUN touch /var/log/cron.log

CMD env >> /etc/environment && cron && tail -f /var/log/cron.log
