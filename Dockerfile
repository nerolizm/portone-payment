FROM golang:1.21-alpine

WORKDIR /app

# 기본 도구 및 빌드 의존성 설치
RUN apk add --no-cache \
    make \
    gcc \
    musl-dev

# CGO 활성화
ENV CGO_ENABLED=1

# 소스 코드 복사
COPY . .

# 의존성 설치
RUN go mod download

# Make 실행을 위한 기본 진입점
CMD ["make", "all"] 