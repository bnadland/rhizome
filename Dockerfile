FROM node:21.5.0 as frontend
WORKDIR /app
COPY ./ ./
RUN npm install && npm run css && npm run js

FROM golang:1.21.6 as backend
WORKDIR /app
ENV CGO_ENABLED 0
COPY ./ ./
COPY --from=frontend /app/internal/assets/static/ ./internal/assets/static/
RUN go run github.com/go-task/task/v3/cmd/task@latest templ && go build -o /app/rhizome ./cmd/rhizome

FROM scratch
WORKDIR /
COPY --from=backend /app/rhizome /app/rhizome
EXPOSE 3000
ENTRYPOINT ["/app/rhizome"]