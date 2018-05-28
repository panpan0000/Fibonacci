package main
import (  
	"testing"
	"net/http/httptest"
	"net/http"
	"specs"
	"strconv"
	//"fmt"
)

//////////////////////////////
//
//
//////////////////////////////
type APITestCase struct{
	input   string    // query string
	expected []uint64 // expected body
	httpCode int      // expected http code
}

//////////////////////////////
//
//  Test Case Array
//
//////////////////////////////
func generateCases( ) []APITestCase{
	s := []APITestCase{
		{ "num=-1"    , []uint64{}, 400 },
		{ "num=-11111", []uint64{}, 400 },
		{ "num=abc"   , []uint64{}, 400 },
	};
	// copy the test case from specs folder
	for _, e := range specs.ShouldSuccess{
		var c APITestCase;
		c.input = "num=" + strconv.FormatUint( e.Input, 10 ) ; // expecting no error in strconv
		c.expected = e.Expected;
		c.httpCode = http.StatusOK;
		s = append( s, c );
	}
	return s;
}
///////////////////////////
// helper: []uint64 to string
// NOTE: the formatting is ad-hock to this code
//////////////////////////
func uint64ArrayToString( arr []uint64) string{
	if( len(arr) == 0 ){
		return "[]";
	}
	s:="[";
	for _,e := range arr[:len(arr)-1] {
		s += strconv.FormatUint( e, 10 ) + "," ;
	}
	s += strconv.FormatUint( arr[len(arr)-1], 10);
	return s +  "]";
}

//////////////////////////////
//
// Unit Test Function
//
//////////////////////////////
func TestFibonacciRestAPI(t *testing.T) {

	cases := generateCases();

	// walk thru all test cases
	for _, c := range cases {
		// create a new request
		req, err := http.NewRequest("GET", "/v1/fib?"+c.input, nil);
		if err != nil {
			t.Fatal(err);
		}
		// create ResponseRecorder
		recorder := httptest.NewRecorder();

		// Call the Http Handler
		GetFibonacci(recorder, req);

		// Check Response code
		status := recorder.Code;
		if( status != c.httpCode) {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK);
		}
		if( recorder.Code == http.StatusOK ){
			// Check response body
			expected := uint64ArrayToString( c.expected ) ; // convert to string
			if( recorder.Body.String() != expected ) {
				t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected);
			}
		}

	} // end of for
}

