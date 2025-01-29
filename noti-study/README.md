# FCM Notification Server

이 프로젝트는 Firebase Cloud Messaging(FCM)을 사용하여 푸시 알림을 전송하는 Go 프로그램입니다. PostgreSQL을 사용하여 디바이스 토큰을 저장하고 관리합니다.

## 주요 기능

1. **디바이스 토큰 등록**: 
   - `/register` 엔드포인트를 통해 사용자의 이메일과 디바이스 토큰을 데이터베이스에 저장합니다.
   - 이메일을 기준으로 토큰 정보를 식별할 수 있습니다.

2. **푸시 알림 전송**: 
   - `/send` 엔드포인트를 통해 발신자의 이메일과 수신자의 이메일을 받아 수신자의 토큰을 조회합니다.
   - 조회된 토큰을 사용하여 FCM을 통해 푸시 알림을 전송합니다.

## 실행 방법

1. PostgreSQL 서버를 실행합니다.
2. 필요한 Go 패키지를 설치합니다:
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/lib/pq
   go get firebase.google.com/go
   go get google.golang.org/api/option
   ```
3. 환경 변수 `FCM_API_KEY`를 설정하여 FCM API 키를 제공합니다.
4. 이 프로그램을 실행하여 알림 서버를 시작합니다.
5. `/register`와 `/send` 엔드포인트를 사용하여 디바이스 토큰을 등록하고 푸시 알림을 테스트합니다.

## TodoList

### 프로젝트 구조화

- [ ] `/cmd/noti-study` 디렉터리 생성
  - 메인 애플리케이션 코드를 이 디렉터리로 이동합니다.
  
- [ ] `/internal` 디렉터리 생성
  - FCM 및 데이터베이스 관련 코드를 `/internal` 디렉터리로 이동합니다.
  - `registerDeviceToken`, `sendPushNotification` 함수들을 `/internal/app` 디렉터리로 이동합니다.

- [ ] `/pkg` 디렉터리 생성
  - 외부에서 사용될 수 있는 라이브러리 코드가 있다면 이 디렉터리로 이동합니다.

### 추가적인 디렉터리

- [ ] `/configs` 디렉터리 생성
  - FCM 및 서버 설정 파일을 이곳에 저장합니다.

- [ ] `/scripts` 디렉터리 생성
  - 빌드 및 배포 스크립트를 이곳에 저장합니다.

- [ ] `/docs` 디렉터리 생성
  - 프로젝트 문서화 파일을 이곳에 저장합니다.

- [ ] `/test` 디렉터리 생성
  - 테스트 코드 및 데이터를 이곳에 저장합니다.

이 TodoList는 프로젝트를 표준 Go 프로젝트 레이아웃에 맞게 구조화하는 데 도움이 될 것입니다. 각 항목을 완료하면서 프로젝트의 유지보수성과 확장성을 높일 수 있습니다. 