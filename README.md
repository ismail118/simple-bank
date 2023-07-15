# simple-bank
backend master class

### External Library
- Migration database: 
https://github.com/golang-migrate/migrate
  - installation golang-migrate cli https://github.com/golang-migrate/migrate/tree/master/cmd/migrate (update to migrate from golang code).
  - migrate using golang code https://pkg.go.dev/github.com/golang-migrate/migrate/v4#readme-use-in-your-go-project
- Postgres Database driver: https://github.com/jackc/pgx
  - example database pool connection using px: https://github.com/jackc/pgx/wiki/Getting-started-with-pgx 
- Framework gin: https://github.com/gin-gonic/gin
- Viper for load config app: https://github.com/spf13/viper
- UUID : https://github.com/google/uuid
- JWT token : https://github.com/golang-jwt/jwt
- PASETO token: https://github.com/o1egl/paseto
- GRPC: https://github.com/grpc/grpc-go
- Log JSON format: https://github.com/rs/zerolog
- Queue task and async processing (used Redis as message broker): https://github.com/hibiken/asynq
- email: https://github.com/jordan-wright/email

### Development Tools
- Docker: https://docs.docker.com/get-docker/
- Github actions for ci Continuous Integration: https://docs.github.com/en/actions/learn-github-actions/understanding-github-actions
- Minikube for for localhost kubernates CD: https://kubernetes.io/docs/tutorials/hello-minikube/
  - Download and instalation: https://minikube.sigs.k8s.io/docs/start/
  - Kubectl cli: https://kubernetes.io/docs/tasks/tools/#kubectl
  - Ingress web service kubernates: https://kubernetes.io/docs/concepts/services-networking/ingress/
- gRPC: https://grpc.io/docs/languages/go/quickstart/
  - proto doc: https://protobuf.dev/programming-guides/proto3/
- For call gRPC server tools: https://github.com/ktr0731/evans
- Grpc gateway for auto http to grpc request: https://github.com/grpc-ecosystem/grpc-gateway
- Docker Redis for queue message broker:
  ```
  $ docker pull redis:7.0.12-alpine
  ```
  to start redis server locally
  ```
  $ make redis
  ```
- Next