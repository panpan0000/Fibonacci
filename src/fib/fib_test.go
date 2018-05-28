// Usage: go test Fib_test.go  Fib.go -v
package fib
import (
	"testing"
	"reflect"
	"specs"
)
////////////////////////////
//
// UT for Fibnacci
//
///////////////////////////

func TestFibnacci(t *testing.T) {
	//Compare the expected result one by one
	for _, testCase := range specs.ShouldSuccess{
		ret := Fibonacci( testCase.Input );
		// ret is the what Fibonacci() returns
		// using  reflect.DeepEqual() to compare two slice
		if(  reflect.DeepEqual( ret, testCase.Expected ) ){
			t.Log("test case ", testCase.Input, " passed");
		}else{
			t.Fatal("Test Failure. test Case Input", testCase.Input)
			t.Log("ret=",ret);
			t.Log("expect=", testCase.Expected)
		}
	}
}

////////////////////////////
//
// Benchmark Testing
//
// Usage: go test fib.go  fib_test.go -test.bench=".*"
//////////////////////////
func BenchmarkFibonacciRestAPI(b *testing.B) {
	input:=100;
	b.StopTimer();
	b.StartTimer();
	for i := 0; i < b.N; i++ {
		Fibonacci( uint64(input) ) ;		
	}

}

