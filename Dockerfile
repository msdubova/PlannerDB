<<<<<<< HEAD
FROM golang

WORKDIR /planner

COPY . .

RUN go build .

CMD [ "/planner/planner" ] 
=======
FROM golang

WORKDIR /planner

COPY . .

RUN go build -o planner .

CMD ["/planner/planner"]
>>>>>>> 42c3a27 (Created Dockerfile, docker-compose)
