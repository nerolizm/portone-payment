# Portone Payment Service

## 기술스택

- go 1.21 이상
- make
- docker, docker-compose
- OpenAPI (swagger)

## 설치
### 의존성 설치
```bash
make dependencies
```

### 요구사항 체크
```bash
make check-all
```

## 환경 설정

`.env` 파일을 프로젝트 루트에 생성하고 다음 내용을 추가합니다: *(.env.sample 파일을 참고해주세요.)*

```bash
# 포트원 API 인증 정보
IMP_KEY=your_imp_key
IMP_SECRET=your_imp_secret
```

## 프로젝트 구조
```
.
├── api/                # OpenAPI Spec
├── cmd/                # 애플리케이션 진입점
├── internal/
│   ├── config/         # 환경 설정
│   ├── handler/        # HTTP 핸들러
│   ├── infrastructure/ # 외부 서비스 연동
│   ├── model/          # 도메인 모델
│   └── service/        # 비즈니스 로직
├── test/
│   └── integration/    # 통합 테스트
└── Makefile            # 빌드 및 개발 도구
```

## 실행
### 기본 실행
아래의 명령어를 실행하면 코드 검사, 테스트, 실행까지 순서대로 진행됩니다.
```bash
make all
```
*명령어에 대해 추가로 확인하고 싶으실 경우 `make help`를 입력하시면 됩니다.*

### 도커 실행
```bash
# make를 통한 도커 컨테이너 실행/종료
make docker-up
make docker-down

# docker-compose를 통한 도커 컨테이너 실행/종료
docker-compose up --build
docker-compose down -v --rmi local --remove-orphans 
```

> 기본적으로 `:8080` 포트로 동작하기 때문에 http://localhost:8080로 접속 가능합니다.

## 테스트
```bash
make test
```

## API 문서

API 문서는 OpenAPI(Swagger) 스펙으로 작성되어 있습니다.

Spec Path: `/api/openapi.json`

URL Path:
```
http://localhost:8080/swagger
```

### 결제 취소
```http
POST /cancel-payment
Content-Type: application/json

{
    "imp_uid": "포트원 결제 고유번호"
}
```

#### 성공 응답
```json
{
    "code": 0,
    "message": "success",
    "response": {
        "status": "cancelled"
    }
}
```

