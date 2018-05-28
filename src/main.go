package main
import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"net/url"
	"os"
	"io"
	"encoding/json"
	"fib"// import means the path of library code.instead of package name
)

///////////////////////
// Http Response Code, ignore others for now.
//////////////////////
const(
	BadRequestErrorCode      = 400
	InternalServerErrorCode  = 500
)

//////////////////////
//  object to json
//////////////////////
func returnJson(value interface{}) string {
	json, _ := json.Marshal(value)
	return string(json)
}

////////////////////////
// error handler
//
// Param
// queryValue: the interpreted value
// w: http reponse
// r: http request
// err: the error message passing from upstream query processing routine
///////////////////////
func badRequestResp( queryValue int, w http.ResponseWriter, r *http.Request, err error ){

	w.WriteHeader( BadRequestErrorCode ); // 400 response code

	queryForm, _ := url.ParseQuery( r.URL.RawQuery );

	fmt.Fprintf(os.Stderr, "Bad Request: queryForm = %s\n",queryForm);
	fmt.Fprintf(w, "invalid num value, should be unsigned number larger than 0.\n");

	if( err == nil ){
		fmt.Fprintf(w, "but input was %d\n" , queryValue);
	}else{
		fmt.Fprintf(w, "input failure, error%s\n", err);
	}
	fmt.Fprintf(w, "Usage: $IP:$port/v1/fibonacci?num=33\n");
}

////////////////////////////////////////////
//
// process query from http request
//
// return the request num value and error object
//////////////////////////////////////////
func getQueryParam( r *http.Request ) ( int, error ) {
	// Retrieve the query by url.ParseQuery
	queryForm, _ := url.ParseQuery( r.URL.RawQuery )
	cnt := -1;
	query_key := "num";

	// the queryForm[query_key][0] should be the value like 123 in '?num=123'
	if( len( queryForm[query_key] )  > 0 ) {
		queryValue := queryForm[query_key][0];
		i, err := strconv.Atoi( queryValue );
		if( err == nil) {
			cnt = i;
		}else{
			return cnt, err;
		}

	}
	return cnt, nil;
}

///////////////////////////////////////
// Func:
//        standard net http handler
//        caculate fibonacci sequence based on total numbers from http request, and response result to http request
//
// Usage:
//        $IP:$port/v1/fibonacci/num?=5   or
//        $IP:$port/v1/fib/num?=5
//        and will response blank when if num?=0
//
//        it will response 400(Bad Request) if the query key invalid or the query value is negative
// 
// Param:
//        w: http.ResponseWriter
//        r: *http.Request
//////////////////////////////////
func GetFibonacci( w http.ResponseWriter, r *http.Request ){
	// Retrieve the query by url.ParseQuery
	// the queryForm[query_key][0] should be the value like 123 in '?num=123'
	cnt, err := getQueryParam( r );
	if( err != nil ) {
		badRequestResp( cnt, w, r, err ); // not elegant enough...
		return;
	}

	if( cnt < 0 ){
		// negative value, report error
		badRequestResp( cnt, w, r, nil );
		return;
	} else{
		fmt.Println("valid user input, total num =",cnt);
		// valid input, call the fib.Fibonacci() and response
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // specific header: json reponse
		//fib_string := fib.Array2String( fib.Fibonacci( uint64(cnt) ) )
		fib_slice := fib.Fibonacci( uint64(cnt))
		io.WriteString(w, returnJson( fib_slice ) )
	}
}

//////////////////////////////
// main
////////////////////////////
func main() {
	http.HandleFunc( "/v1/fib"      , GetFibonacci );
	http.HandleFunc( "/v1/fibonacci", GetFibonacci );
	port:=":8008";
	fmt.Println("Try to bind to port ", port);
	log.Fatal( http.ListenAndServe(port, nil) );
}
