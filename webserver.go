//
// webserver.go
//
// GoLang을 통한 웹서버의 간단한 구현 예제입니다.
//
// 사용 예:
//
//   # 백그라운드에서 사용하기
//   $ go run webserver &
//
//   실행중에 브라우저로 페이지를 방문하세요.
//   It responds in one of several ways : 몇 가지 방법으로 응답합니다.
//
//  (0)  /home으로는 home HTML 페이지를 보내줍니다. 이것은 AJAX secondary GET을 수행합니다.
//
//   (1) /generic URL은 조금의 text/plain 진단을 보내줍니다.
//       URL: http://localhost:8097/generic/page?color=purple
//       browser (text/plain) :
//           FooWebHandler says ...
//             request.Method      'GET'
//             request.RequestURI  '/generic/page?color=purple'
//             request.URL.Path    '/generic/page'
//             request.Form        'map[color:[purple]]'
//             request.Cookies()   '[testcookiename=testcookievalue]'
//
//   (2) /item/텍스트스트링 으로 URL을 입력하면 간단한 JSON 응답을 해줍니다.
// 	  실제 앱에서는 textstring은 item의 이름이 될 수 있습니다. 그리고 응답은 그에 대한 설명이 되곤합니다.
//
//       URL: http://localhost:8097/item/yellow
//       browser (application/json) :
//           {"name":"yellow", "what":"item"}
//
//   (3) 다른 페이지는 에러페이지를 출력해줍니다.
//
//       URL: http://localhost:8097/other/path
//       browser :
//           404 page not found
//
// 매 방문은 간단한 쿠키를 설정해줍니다. 첫번 째 방문 이후로는 요청을 할 수 있습니다.
//
// AJAX 설정을 하려면, 너는 정보와 submission의
// 데이터를 URL에 넣기 위한 AJAX 인코드 요청을 결정해야할 것이다.
// REST API는 여기 있는 /item/name 예제와 같이 요청하거나 전송한
// 정보를 경로에 입력하는 URL과 함께 GET or PUT을 사용합니다.
// 또는 URL의 ?dll=value 부분에 전달된 양식이나 데이터를 사용할 수도 있지만
// 그다지 깨끗하지는 않다고 생각합니다.
// 그런 다음 클라이언트의 Javascript에 데이터를 다시 전달하려면
//	item/name 예제와 같이 JSON을 사용하는 것이 좋습니다.

// For a discussion of REST see
// en.wikipedia.org/wiki/Representational_state_transfer#Central_principle
//
// GO는 또한 서드파티 라이브러리로 gorilla/mux 라는 흥미로운 것이 있습니다.
// URL에서 더 간지나는 방법으로 정보를 추출하거나
// 어떤 함수가 요청에 응할지 정할 수 있다.
// http://www.gorillatoolkit.org/pkg/mux 여기서 참고하도록
//
// 참고할 만한 자료들
//   http://golang.org/pkg/net/http      particularly #Request
//   http://golang.org/pkg/net/url/#URL  what's in request.URL
//   https://devcharm.com/pages/8-golang-net-http-handlers
//   http://www.alexedwards.net/blog/a-recap-of-request-handling
//   http://blog.golang.org/json-and-go
//
// Jim Mahoney | cs.marlboro.edu | MIT License | March 2014
// KOR (한국어) 번역 - git.imdhs.one
//					git.imdhson.com
//					github.com/imdhson

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func SetMyCookie(response http.ResponseWriter) {
	// 응답에 간단한 쿠키를 추가합니다.
	cookie := http.Cookie{Name: "testcookiename", Value: "testcookievalue"}
	http.SetCookie(response, &cookie)
}

// /generic URL 형식에 대한 응답
func GenericHandler(response http.ResponseWriter, request *http.Request) {

	// 쿠키를 설정하고 MIME type을 http 헤더에 설정
	SetMyCookie(response)
	response.Header().Set("Content-type", "text/plain")

	//URL을 Parse하고 POST 데이터를 요청에 포함합니다.
	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	//text 진단 결과를 클라이언트에게 전달
	fmt.Fprint(response, "FooWebHandler says ... \n")
	fmt.Fprintf(response, " request.Method     '%v'\n", request.Method)
	fmt.Fprintf(response, " request.RequestURI '%v'\n", request.RequestURI)
	fmt.Fprintf(response, " request.URL.Path   '%v'\n", request.URL.Path)
	fmt.Fprintf(response, " request.Form       '%v'\n", request.Form)
	fmt.Fprintf(response, " request.Cookies()  '%v'\n", request.Cookies())
}

// /home에 대한 응답으로 html home page를 응답해줌
func HomeHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-type", "text/html; charset=utf-8") //imdhson 수정함
	webpage, err := ioutil.ReadFile("home.html")
	if err != nil {
		http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
	}
	fmt.Fprint(response, string(webpage))
}

// /item/...에 대한 응답
func ItemHandler(response http.ResponseWriter, request *http.Request) {

	// 쿠키를 설정하고 MIME type을 http 헤더에 설정
	SetMyCookie(response)
	response.Header().Set("Content-type", "application/json")

	// 클라이언트에게 재전송할 일부 샘플 데이터
	data := map[string]string{"what": "item", "name": ""}

	// URL 형식이 /item/name이 맞는가?
	var itemURL = regexp.MustCompile(`^/item/(\w+)$`)
	var itemMatches = itemURL.FindStringSubmatch(request.URL.Path)
	// itemMatches는 regex 매치로 다음과 같이 작동  ["/item/which", "which"]
	if len(itemMatches) > 0 {
		// 참일 경우 JSON을 클라이언트에게 전송
		data["name"] = itemMatches[1]
		dataall := itemMatches
		dataall = append(dataall, "This is long JSON data for calculation for bytes.")
		json_bytes, _ := json.Marshal(data)
		json_all, _ := json.Marshal(dataall)
		fmt.Fprintf(response, "%s\n", json_bytes)
		fmt.Fprintf(response, "%s\n", json_all)
	} else {
		// 거짓일 경우 오류 전달
		http.Error(response, "404 page not found", 404)
	}
}

// imdhson이 연습용으로 추가한 함수.
// 디버그 목적으로 사용하기 적합
func splitHangeul(in string) []rune {
	var out []rune
	for _, v := range in {
		out = append(out, rune(v))
	}
	return out
}

func main() {
	port := 8080
	portstring := strconv.Itoa(port)

	// 요청 핸들러를 두가지의 URL 패턴에 대응하게 생성함
	//  문서는 pattern 에 대해 확정성이 부족해보임. 그럼 gorilla/mux를 쓰는것이 좋다
	mux := http.NewServeMux()

	mux.Handle("/home", http.HandlerFunc(HomeHandler))
	mux.Handle("/item/", http.HandlerFunc(ItemHandler))
	mux.Handle("/generic/", http.HandlerFunc(GenericHandler))

	//  지정된 포트로 서버를 가동하여 listen 시작
	// (개인적으로 생각하길 서버 이름도 여기서 설정가능 할 것이다.)
	log.Print("Listening on port " + portstring + " ... ")
	err := http.ListenAndServe(":"+portstring, mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
