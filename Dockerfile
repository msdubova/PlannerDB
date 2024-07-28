FROM golang

WORKDIR /planner

COPY . .

RUN go build -o planner .

CMD ["/planner/planner"]

