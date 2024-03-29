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
// Author: dlyongemallo@gmail.com (David Yonge-Mallo)

package main

import (
	"fmt"
	"quantum"
)

// Implement the Deutsch problem.
func main() {
	qreg := quantum.NewQReg(2, 1)
	fmt.Println(qreg)

	quantum.HadamardReg(qreg)
	quantum.NewRealArrayGate([]float64{
		0, 1, 0, 0,
		1, 0, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}).ApplyReg(qreg)
	quantum.Hadamard(qreg, 1)
	if qreg.BMeasure(1) == 0 {
		fmt.Println("constant")
	} else {
		fmt.Println("balanced")
	}
}
