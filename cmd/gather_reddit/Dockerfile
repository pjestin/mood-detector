FROM golang:1.19.1

ADD . /app
WORKDIR /app

RUN apt-get update && apt-get -y upgrade
RUN apt-get -y install cron

RUN go build ./cmd/gather_reddit

ADD cmd/gather_reddit/crontab /etc/cron.d/mood-detector-cron
RUN chmod 0644 /etc/cron.d/mood-detector-cron
RUN touch /var/log/cron.log

CMD env >> /etc/environment && cron && tail -f /var/log/cron.log
