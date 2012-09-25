// Copyright 2011 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//      Unless required by applicable law or agreed to in writing, software
//      distributed under the License is distributed on an "AS IS" BASIS,
//      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//      See the License for the specific language governing permissions and
//      limitations under the License.
//
// Author: conleyo@google.com (Conley Owens)

package main

import (
	"fmt"
	/*"math"*/
	"os"
	/*"quantum"*/
	"math/rand"
)

func gcd(a int, b int) int {
	for a != 0 {
		a, b = b%a, a
	}
	return b
}

/*func isPrime(n int) bool {*/
/*if n % 2 == 0 {*/
/*return false*/
/*}*/
/*for r := 3; r < n; r++ {*/
/*if gcd(n, r) != 1 {*/
/*return false*/
/*}*/
/*if isPrime(r) {*/
/*let q be the largest factor of r-1*/
/*if (q > 4sqrt(r)log n) and (n(r-1)/q is not 1 (mod r)) then break*/
/*}*/
/*r := r+1*/
/*}*/

/*for a = 1 to 2sqrt(r)log n {*/
/*if ( (x-a)n is not (xn-a) (mod xr-1,n) ) then output COMPOSITE*/
/*}*/

/*output PRIME;*/
/*}*/

// TODO: Replace with a polynomial time algorithm
func isPowerOfPrime(n int) bool {
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return true
		}
	}
	return false
}

func newa(bign int) int {
	return rand.Intn(bign-1) + 1
}

// TODO: Make a real period function
func period(bign int, a int) int {
	return 0
}

func powermod(x int, y int, n int) int {
	ret := 1
	for i := 0; i < y; i++ {
		ret *= x
		ret %= n
	}
	return ret
}

func main() {
	/*bign := 6 // 2 * 3*/
	/*n := int(math.Ceil(math.Log2(float64(bign))))*/
	/*isPowerOfPrime(n)*/
	/*qreg := quantum.NewQReg(3*n, 0)*/
	/*quantum.HadamardRange(qreg, n, qreg.Size())*/
	/*qreg.Print()*/
	/*fmt.Println()*/

	/*var factor int*/
	/*for a := newa(bign); ; a = newa(bign) {*/
		/*factor = gcd(a, N)*/
		/*if factor != 1 {*/
			/*break*/
		/*}*/

		/*r := period(bign, a)*/
		/*if r%2 == 0 && powermod(a, r, bign) != bign-1 {*/
			/*pow := math.Pow(a, r/2)*/
			/*factor = gcd(pow+1, bign)*/
			/*break*/
		/*}*/
	/*}*/

	/*fmt.Printf("%d %d\n", common, N/common)*/
	fmt.Println("This doesn't work yet.  Be patient.")
	os.Exit(0)
}
