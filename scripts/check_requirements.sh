#!/bin/bash

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 에러 카운트
errors=0

echo "Checking requirements..."

# .env 파일 체크
if [ -f .env ]; then
    echo -e "${GREEN}✓${NC} .env file exists"
else
    echo -e "${RED}✗${NC} .env file not found"
    errors=$((errors + 1))
fi

# make 명령어 체크
if command -v make >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} make is installed"
else
    echo -e "${RED}✗${NC} make is not installed"
    errors=$((errors + 1))
fi

# docker 체크
if command -v docker >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} docker is installed"
else
    echo -e "${RED}✗${NC} docker is not installed"
    errors=$((errors + 1))
fi

# docker-compose 체크
if command -v docker-compose >/dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} docker-compose is installed"
else
    echo -e "${RED}✗${NC} docker-compose is not installed"
    errors=$((errors + 1))
fi

# 결과 출력
if [ $errors -gt 0 ]; then
    echo -e "\n${RED}Found $errors error(s)${NC}"
    echo "Please install missing requirements and try again"
    exit 1
else
    echo -e "\n${GREEN}All requirements are satisfied!${NC}"
    exit 0
fi 