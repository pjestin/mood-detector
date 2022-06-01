FROM golang:1.18

ADD . /app
WORKDIR /app

RUN apt-get update
RUN apt-get -y install cron

RUN go build ./cmd/gather_mood

ADD cmd/gather_mood/docker/crontab /etc/cron.d/mood-detector-cron
RUN chmod 0644 /etc/cron.d/mood-detector-cron
RUN touch /var/log/cron.log

CMD env >> /etc/environment && cron && tail -f /var/log/cron.log
