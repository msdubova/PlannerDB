FROM golang

WORKDIR /planner

COPY . .

RUN go build .

CMD [ "/planner/planner" ] 