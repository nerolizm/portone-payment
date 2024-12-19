# Portone Payment API

## 기능
- 결제
- 결제 취소

## 시작하기

### 요구사항
- Go 1.21 이상

### 환경 변수 설정
`.env` 파일을 프로젝트 루트에 생성하고 다음 내용을 입력하세요:

```
PORT=:8080
IMP_KEY=your_imp_key
IMP_SECRET=your_imp_secret
```

### 실행

```bash
# 의존성 설치
go mod download

# 서버 실행
go run cmd/main.go
```

## API 문서

### 결제 취소
```
POST /cancel-payment
Content-Type: application/json

{
    "imp_uid": "포트원 결제 고유번호"
}
```

## 프로젝트 구조
```
.
├── cmd/                # 애플리케이션 진입점
├── internal/
│   ├── config/         # 환경 설정
│   ├── handler/        # HTTP 핸들러
│   ├── infrastructure/ # 외부 서비스 연동
│   ├── model/          # 도메인 모델
│   └── service/        # 비즈니스 로직
└── static/             # 정적 파일
```
