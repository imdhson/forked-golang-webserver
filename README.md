golang-webserver
================

An example of a webserver in the Go programming language,
including a jQuery ajax request .
GO언어를 사용한 웹서버의 예시입니다. jQuery AJAX 요청을 포함하고 있습니다.

    $ go run webserver.go
    
... 브라우저로 이곳을 접속하세요: http://localhost:8080/home
home.html을 반환합니다.

http://localhost:8097/item/foo 이 웹페이지는 2nd GET 요청을 보냅니다.
JSON 응답을 해줍니다. {"name":"foo","what":"item"}, 
homepage의 span element에 포함된 일부와 동일합니다.


