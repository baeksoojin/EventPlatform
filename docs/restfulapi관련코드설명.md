# Restful API 구현

> Go로 구현하는 클라우드 네이티브 애플리케이션 [저자 : 미나안드라오스] 책을 기반으로 작성하였습니다.

## golia/mux

Go가 가진 웹 서드파티 툴킷을 강력한 웹 앱 구축이 빠르고 효율적으로 되도록 서로 함께 도와주는 Go 패키지들의 Collection으로 구성되어있다.<br>

- golia/mux<br>
    고릴라 웹 툴키시의 핵심 패키지이다.<br>
    요청 라우터와 디스패처로 기술된다.<br>
- 사용<br>
  - 설치 : ``import "github.com/golia/mux``
  - mux 패키지를 사용한 라우터 생성
    ``r = mux.NewRouter()``<br>
  - 예시<br>
    공통 상대 url인 /events에 대해서는 subrouter를 생성할 수 있다.<br>
    ``eventsrouter := r.PathPrefix("events").Subrouter()``
    /events 로 시작하는 어떤 url 경로도 잡을 수 있도록 사용되는 PathPrefix 메서드를 사용한다.<br>

----


## EventPlatform

1. 예약 프로세스 : 행사에 이용가능한 좌석 확인, 동일한 이름으로 미리 예약된 것이 없는지 확인, 예약을 저장하는 작업<br>
2. 이벤트 처리 : 콘서트, 연국 등의 모든 종류의 지원돼야하는 행사를 인지해야함, 행사 주소, 전체 좌석수, 행사 기간등<br>
3. 검색 처리 : 예약과 이벤트 조회를 위한 효율적인 검색 필요<br>

### Event Microservice

#### 3가지 Task
1. 항목을 검색 => get / parameter
2. 동시에 모든 이벤트를 조회 => get / parameter
3. 새로운 이벤트 생성 => post / json

#### rest.go
- `Methods()`와 `Path()`를 활용한 HTTP 메서드 정의 및 상대 URL 경로를 정의<br>
method, url을 설정한 router를 handler에 넘긴다.<br>
handler는 HTTP요청을 하나의 동작과 어떻게 연결시킬지에 대한 것이다.<br>
핸들러는 `HandlerFunc()`은 func(http.ResponseWriter, *http.Request)의 형태의 인수를 취한다.<br>
이때, 구조화하기 위해서 handler method를 handlers.go에서 처리한다.<br>

- ex) new event 생성
    ``eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)"``

- 라우터, 경로, 핸들러 정의를 마친 이후)<br>
`return http.ListenAndServe(endpoint, r)`를 net/http 패키지를 활용해서 요청을 리슨하는 TCP 주소를 지정한다.<br>

#### handlers.go

- router의 HandlerFunc()안에는 2가지 중요한 인수를 가져야한다.<br>
  HandleFunc()로 전달되는 함수는(3가지 서비스가 있으니까 3가지가 존재함) 모두 func(http.ResponseWriter, *http.Request) Signature를 지원해야한다.<br>

    이때, handler는 이들을 전부 관리하고자 생성된 Go 구조체 객체이다.<br>

  - eventServiceHandler
      3가지 업무에 대해서 핸들러를 지원하기 위해서 eventServiceHandler라는 custom 구조체 type을 정의하여 각 업무에서 포인터를 활용한 리시버로 persistence를 지킬 수 있다.<br>

      ```
      type eventServiceHandler struct(){}
  
      func(eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request){}
      func(eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request){}
      func(eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request){}
  
     ```
    
    자세한 구현은 persistence(MSA의 지속성 계층)에서 살펴볼 예정이다!<br>
  