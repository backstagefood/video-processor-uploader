FROM golang:1.24.4-alpine AS base

ARG VERSION
ARG PROJECT_NAME

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'github.com/backstagefood/video-processor-uploader/internal/controller/handlers.Version=${VERSION}' -X 'github.com/backstagefood/video-processor-uploader/internal/controller/handlers.ProjectName=${PROJECT_NAME}'" -o video-processor-uploader ./cmd/app/.

FROM alpine

ENV GIN_MODE=release

# Criar diretório de trabalho
WORKDIR /app

COPY --from=base /app/video-processor-uploader .

CMD ["/app/video-processor-uploader"]