# 몽고DB

> 몽고DB는 NoSQL로 **문서저장형 데이터베이스 엔진**이다.

## 몽고DB의 키워드 2개

1. NoSQL
2. 문서 저장소(Document store)

NoSQL : 관계형 데이터에 깊숙이 의존하지 않는다는 점을 암시적으로 나타낸다.<br>
Document store : collection, document는 몽고DB에서 중요한개념이다.<br>
    Document store를 사용하면 각 데이터는 고유한 id를 가진 분리된 문서(document)에 저장될 것이다.<br>
    몽고 DB는 collectioin과 Document로 구성되고 document들이 모여 collection이 형성된다.<br>
    이때, collection이 table이 되고 그 record가 document가 된다.<br>

## MyEvent의 event microservice

### 모델 작성

- model :  데이터베이스의 데이터와 매치되는 field를 담는 데이터 구조이다.<br>
    Go에서는 struct type을 활용해서 model을 구조화하여 생성한다.<br>

```
type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}

type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}


```

### 사용방법

- `go get gopkg.in/mgo.v2`<br>
    - mgo를 사용할 수 있다. bson을 import할 수 있다.<br>
- mongolayer를 통한 지속성 관리<br>
    1. collection의 이름을 나타내니는 몇가지 상수를 생성해보자.<br>
    ```
      const (
        DB     = "myevents"
        USERS  = "users"
        EVENTS = "events"
    )
    ```
  2. 외부 노출을 위한 session객체 활용<br>
    ```
    type MongoDBLayer struct {
    session *mgo.Session
    }
    ```
  3. mongoDB implementer 생성을 위한 constructor<br>
  생성과 동시에 초기화하는 코드를 실행해야한다.<br>
  원하는 데이터베이스로 연결해서 session handler를 얻는 것을 통해서 초기화를 진행한다.<br>
  ```
  func NewMongoDBLayer(connection string) (persistence.DatabaseHandler, error) {
    s, err := mgo.Dial(connection)
    return &MongoDBLayer{
        session: s,
    }, err
    }
  ```
    이때, Dial 메서드는 연결 문자열(데이터베이스와의 연결을 맺는데 필요한 정보)인 connection을 넘기면 세션을 반환해주는 함수이다.<br>
    ```
      func Dial(url string) (*Session, error) {
        session, err := DialWithTimeout(url, 10*time.Second)
        if err == nil {
            session.SetSyncTimeout(1 * time.Minute)
            session.SetSocketTimeout(1 * time.Minute)
        }
        return session, err
    }
    ```
  4. 구현체 생성<br>
  - receiver를 사용해서 session에 접근한다.<br>
    ```
    func(mgoLayer *MongoDBLayer) AddEvent(e persistence.Event) ([]byte, error){...}
    func (mgoLayer *MongoDBLayer) FindEvent(id []byte) (persistence.Event, error) {...}
    func (mgoLayer *MongoDBLayer) FindEventByName(name string) (persistence.Event, error) {...}
    func (mgoLayer *MongoDBLayer) FindAllAvailableEvents() ([]persistence.Event, error) {...}
    ```
  
  - connection pool에서의 session 획득 및 반환<br>
    connection pool에서 session을 획득하기 위해서 getFreshSession() method를 활용한다.<br>
    ```
    func (mgoLayer *MongoDBLayer) getFreshSession() *mgo.Session {
    return mgoLayer.session.Copy()
    }

    ```
    `mgoLayser.session.Copy()`를 통해서 연결 풀에서 **새로운 새션을 얻는 작업**을 진행한다.<br>
      이후, 해당 세션이 작업을 종료한 후에는 연결 풀로 반납되는 것을 보장해야하는데 이때, defer를 사용하면 된다.<br>
  - mongoDB로의 데이터 조회 및 저장 등 <br>
    session을 통해서 DB에 접근하고 원하는 collection에 접근해서 document를 저장, 조회, 수정 등의 작업을 실행하면 된다.<br>
    이때, `s.DB([database이름]).C([collection이름])`으로 collection에 session instance인 s를 통해서 접근이 가능하다.<br>
    Event처리를 위한 4가지 핸들러 구현체에서 사용되는 구문에 대해서 알아보자.<br>
    **새로운 이벤트 추가** => `s.DB(DB).C(Event).Insert(e)`<br>
    **id로이벤트 찾기** => `s.DB(DB).C(Event).FindId(bson.ObjectId(id)).One(&e)`<br>
    **이름으로 이벤트 찾기** => `s.DB(DB).C(Evnet).Find(bson.M{"name":name}).One(&e)` 이때 M은 map이다.<br>
    **모든 이벤트 찾기** => `s.DB(DB).C(Event).Find(nil).All(&events)` nil인수를 가지면 collection에서 찾아진 모든 것을 반환한다.<br>
  
## handler와 연결

handler 계층에서 DatabaseHandler를 지원하기 위해서 구조체를 생성해야한다.<br>
```
type eventServiceHandler struct {
	dbhandler persistence.DatabaseHandler
}
```
항상 구조체를 생성하면 그에 맞는 초기화 메서드를 지원해야한다.<br>
새로운 연결마다 새로운 eventServiceHandler로 접근하기 위해서 객체를 초기화하고 생성자를 작성할 필요가 있다.<br>
```
func NewEventHandler(databasehandler persistence.DatabaseHandler) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler: databasehandler,
	}
}
```

다음은 사용자의 요청에 맞게 응답을 하도록 handler를 작성하면 된다.<br>
handler.md에서 찾아보면 된다.<br>
----------

## 추가 개념

- bson : mgo adapter에서 찾을 수 있는 몽고 DB와 GO의 통신을 위해 사용되는 서드파티 프레임워크<br>
    bson.ObjectId타입은 해당 ID의 유효성을 검증하고자 추후 코드에서 사용할 수 있는 몇가지 유용한 메서드를 제공한다.<br>
    bson은 저장된 문서들에서 데이터를 나타내고자 몽고 DB에 사용되는 데이터 포맷이다.<br>
