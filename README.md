# GGNet(Golang Game Network)

Golang 기반의 게임 서버 구현

## 인프라

### 개발단계에서는 Docker 이용하여 구축

- Redis
- Mongodb

### AWS를 이용한 인프라 구축

- AWS EC2
- AWS ElasticCache
- AWS RDB or NoSQL

## 진행 순서

### 게임 서버 일반 구현

#### 1. 게임 서버 실행

- [] 게임 서버는 설정된 IP 주소와 포트 번호를 사용해 네트워크 상에서 클라이언트의 접속을 대기.

#### 2. 클라이언트 접속 대기

- [] 서버는 연결 요청을 대기하며, 동시에 여러 클라이언트의 접속을 처리할 수 있도록 구성

#### 3. 클라이언트 접속 수락

- [] 클라이언트가 접속을 시도하면, 서버는 접속을 수락하고 세션을 생성하여 클라이언트와 연결하여 관리.

#### 4. 세션 정보 저장

- [] 생성된 세션 정보는 Redis 에 저장하여 관리

#### 5. 인증 및 로그인 처리

- 클라이언트는 아이디와 비밀번호를 서버에 전송하고, 서버는 받은 인증 정보를 검증하여 인증토큰을 발급.

#### 6. 사용자 데이터 로드 및 초기화

- 로그인이 완료되면, 서버는 사용자의 정보를 Redis 에 업데이트하고 사용자에게 제공.

#### 7. 접속 종료 및 세션 삭제

- 클라이언트가 접속을 종료하면, 서버는 해당 클라이언트의 세션을 종료하고, Redis 에서 세션 정보를 삭제.

### 채널 시스템 구현

#### 1. 채널 생성

- 서버는 설정한 채널의 수만큼 채널을 생성하여 Redis 에 저장.

#### 2. 채널 리스트 요청

- 클라이언트가 서버에 채널 리스트를 요청. 서버는 채널의 목록을 클라이언트에게 전송.

#### 3. 채널 입장 요청

- 클라이언트는 특정 채널에 입장하기 위해 인증 토큰과 입장하고자 하는 채널의 식별 정보가 포함하여 요청.

#### 4. 채널 분배 및 관리

- 게임 서버는 클라이언트의 요청을 받고, 사용자의 레벨, 경험치등을 고려하여 채널을 할당하고, 사용자의 세션 정보를 업데이트.

#### 5. 채널 입장 승인

- 요청한 채널에 입장할 수 있도록 서버에서 승인하고 클라이언트에게 응답.

#### 5. 채널 이동 및 로그아웃

- 클라이언트가 다른 채널로 이동하길 원하면, 새로운 채널 입장 요청을 서버에 전송하고, 서버는 새로운 채널에 입장할 수 있도록 처리.
- 클라이언트가 접속을 종료할 경우, 서버는 사용자를 현재 채널에서 제거하고, 사용자의 세션 정보를 정리.

### 채팅 시스템 구현

#### 1. 공지 기능 구현

- 서버는 모든 클라이언트에게 공지 메시지를 전송할 수 있는 기능을 구현.

#### 2. 채널 채팅 기능

- 같은 채널에 있는 클라이언트들 간의 채팅 기능 구현.

#### 3. 귓속말 기능

- 특정 클라이언트에게만 메시지를 전송할 수 있는 귓속말 기능 구현.

### 채널 밸런싱 기능 구현

- 서버는 각 채널의 부하를 균등하게 분배하기 위한 로직을 구현.
- 특정 채널의 인원이 많아질 경우 새로운 채널로 유도하는 기능 구현.