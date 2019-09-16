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
	"math"
	"os"
	"quantum"
)

func main() {
	n := 3
	qreg := quantum.NewQReg(n+1, 0)
	quantum.HadamardRange(qreg, 1, n+1)
	h := quantum.NewHadamardGate(1)
	u_f := quantum.NewClassicalGate(func(x int) int {
		if x>>1 == 5 {
			return x ^ 1
		}
		return x
	},
		n+1)
	states := 1 << uint(n)
	d := quantum.NewDiffusionGate(n)
	iterations := int((math.Pi * math.Sqrt(float64(states))) / 4.0)
	for i := 0; i < iterations; i++ {
		qreg.BSet(0, 1)
		h.ApplyRange(qreg, 0)
		u_f.ApplyReg(qreg)
		d.ApplyRange(qreg, 1)
	}
	fmt.Printf("Found %d\n", qreg.Measure()>>1)
	os.Exit(0)
}
