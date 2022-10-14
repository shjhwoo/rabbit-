/*
what if we need to run a function on a remote computer and wait for the result? 
Well, that's a different story. This pattern is commonly known as Remote Procedure Call or RPC.
*/

원격으로 메세지를 주고받을 경우 RPC 개념을 사용한다

간혹 메세지를 보낸 주체는 메세지를 받는 컨슈머로부터 응답을 돌려받아야 할 때가 있는데, 이때 callback queue를 사용한다
즉, 컨슈머는 누구에게 응답을 돌려줄지 알아야 하기 때문에 메세지 보낼 경우에 누가 보낸건지 정보를 담아 보내는 것이다. 

mq에는 이를 설정하기 위해 reply_to 옵션을 사용한다. callback queue에 이름을 부여하기 위해 사용

correlation id?
모든 원격 요청마다 callback queue를 생성하는 것은 비효율적이기 때문에 클라이언트 하나당 하나의 콜백큐를 생성해줘야 한다. 
하지만 그렇게 할 경우 어느 메세지에 대해서 응답을 해줘야 할 지 모르는 상황이 생겨버린다고 함
그래서 이 속성을 사용하게 된다. 
즉, 모든 요청마다 고유의 값을 설정하게 되는데 이것이 correlation id이다.

* 콜백 큐에 있는 모르는 메세지를 무시하는 이유?
서버측의 레이스 컨디션 때문임
원격에 있는 서버가 ack 메시지를 보내기 전에 응답을 보내고 죽을 수도 있음. 이 경우 원격서버를 재시작하여 들어온 메세지를 다시 처리한다. 
즉 메세지를 보낸 입장에서는 결국 응답을 두번 받게 된다는 소리임. 
이 때문에 RPC는 멱등성을 가져야 하며, 응답받는 쪽은 중복되는 응답을 유연하게 처리해야 한다


go run rpc_server.go

go run rpc_client.go 30