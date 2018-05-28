package fib
import (
	"fmt"
    "strconv"
)
func PrintArray( s []uint64 ){ 
    for _,e := range s{
        fmt.Print(e, "\t");  
    }
}
func Array2String( s[]uint64) string{
    str:="";
	for _,e :=range s{
		str += strconv.FormatUint(e,10) + " " ;
	}
	return str;
}
/////////////////////////////////////////////
// Func:
//   Fibonacci: F(n)=F(n-1)+F(n-2) (n>=2)
//   Return the first n Fibonacci array. example, if n == 5, return [0,1,1,2,3]
//   Note: if n == 0, return an blank array instead.
// Param:
//   n: the length of the returned Fibonacci numbers array
// Output:
//   the first n Fibonacci array.
/////////////////////////////////////////////
func Fibonacci( n uint64 ) []uint64 {
    if( 0 == n ){
	    return ( []uint64{} );
    }
    if( 1 == n ){
	    return ( []uint64{ 0 } );
    }
	// cache is the staging array for DP(Dynamic Processing)
	cache :=make( []uint64, n, n);
	ret   :=make( []uint64, n, n);
	cache[0] = 0;
	cache[1] = 1;
	copy( ret, cache ); // ret <- cache
	for i:=uint64(2); i<n; i++ {
        ret[i] = cache[i-1] + cache[i-2];
		cache[i] = ret[i];
	}
    return ret;
}

