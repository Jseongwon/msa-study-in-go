# Redis WebSocket 채팅 서버

이 프로젝트는 Redis와 WebSocket을 사용하여 실시간 채팅 서버를 구현한 Go 프로그램입니다.

## 주요 기능

1. **Redis 초기화**: 
   - `initRedis` 함수는 Redis 클라이언트를 초기화하고 연결을 확인합니다.
   - Redis 서버에 연결하여 채팅 메시지를 게시하고 구독합니다.

2. **WebSocket 연결 처리**: 
   - `handleConnections` 함수는 WebSocket 연결을 처리합니다.
   - 클라이언트로부터 수신된 메시지를 Redis에 게시합니다.

3. **Redis 구독자**: 
   - `redisSubscriber` 함수는 Redis 채널을 구독합니다.
   - 수신된 메시지를 WebSocket 클라이언트에 브로드캐스트합니다.

4. **메시지 브로드캐스터**: 
   - `broadcaster` 함수는 수신된 메시지를 모든 WebSocket 클라이언트에 전송합니다.

5. **메인 함수**: 
   - `main` 함수는 Redis 구독자와 브로드캐스터를 고루틴으로 실행합니다.
   - Gin 웹 서버를 설정하여 WebSocket 연결을 처리합니다.

## 실행 방법

1. Redis 서버를 실행합니다.
2. 이 프로그램을 실행하여 채팅 서버를 시작합니다.
3. 웹 브라우저에서 `http://localhost:8080/ws`에 접속하여 WebSocket 연결을 테스트합니다.

## TodoList

### 프로젝트 구조화

- [ ] `/cmd/redis-study` 디렉터리 생성
  - 메인 애플리케이션 코드를 이 디렉터리로 이동합니다.
  
- [ ] `/internal` 디렉터리 생성
  - WebSocket 및 Redis 관련 코드를 `/internal` 디렉터리로 이동합니다.
  - `handleConnections`, `redisSubscriber`, `broadcaster` 함수들을 `/internal/app` 디렉터리로 이동합니다.

- [ ] `/pkg` 디렉터리 생성
  - 외부에서 사용될 수 있는 라이브러리 코드가 있다면 이 디렉터리로 이동합니다.

### 추가적인 디렉터리

- [ ] `/configs` 디렉터리 생성
  - Redis 및 서버 설정 파일을 이곳에 저장합니다.

- [ ] `/scripts` 디렉터리 생성
  - 빌드 및 배포 스크립트를 이곳에 저장합니다.

- [ ] `/docs` 디렉터리 생성
  - 프로젝트 문서화 파일을 이곳에 저장합니다.

- [ ] `/test` 디렉터리 생성
  - 테스트 코드 및 데이터를 이곳에 저장합니다.

이 TodoList는 프로젝트를 표준 Go 프로젝트 레이아웃에 맞게 구조화하는 데 도움이 될 것입니다. 각 항목을 완료하면서 프로젝트의 유지보수성과 확장성을 높일 수 있습니다.

reference
- https://github.com/golang-tutorials/redis-chat-server
- https://github.com/golang-standards/project-layout/blob/master/README_ko.md
