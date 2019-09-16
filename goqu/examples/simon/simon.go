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

func getParity(x int) int {
	p := 0
	for i := x; i != 0; i >>= 1 {
		p ^= i & 1
	}
	return p
}

func solveLinearSystem(system []int) int {
	bits := len(system)
	for pos := 0; pos < bits; pos++ {
		mask := 1 << uint(pos)
		// Find the first element with this mask and swap with the
		// current position
		for i := pos; i < bits; i++ {
			if system[i]&mask != 0 {
				temp := system[i]
				system[i] = system[pos]
				system[pos] = temp
				break
			}
		}
		for i := pos + 1; i < bits; i++ {
			if system[i]&mask != 0 {
				system[i] ^= system[pos]
			}
		}
	}
	c := 0
	for pos := bits - 1; pos < bits; pos-- {
		bit := getParity(system[pos]^(1<<uint(pos))) << uint(pos)
		c |= bit
		if bit == 0 {
			for i := 0; i < pos; i++ {
				system[i] &= math.MaxInt32 ^ (1 << uint(pos))
			}
		}
	}
	return c
}

func main() {
	bits := 3
	h := quantum.NewHadamardGate(bits)
	u_f := quantum.NewClassicalGate(func(x int) int {
		secret := 5 // 101
		value := -1
		y := x >> 3
		smaller := y ^ secret
		if y < smaller {
			smaller = y
		}
		for i := 0; i <= smaller; i++ {
			if i < i^secret {
				value++
			}
		}
		return x ^ value
	},
		2*bits) // This function maps {0, 1}^6 -> {0, 1}^6
	system := make([]int, bits) // This holds the values in our linear
	// system
	system_map := make(map[int]bool)
	// Fill our linear system
	for len(system_map) < bits {
		qreg := quantum.NewQReg(2*bits, 0)
		h.ApplyRange(qreg, bits)
		u_f.ApplyReg(qreg)
		h.ApplyRange(qreg, bits)
		value := qreg.Measure() >> uint(bits)
		if value != 0 && system_map[value] == false {
			system[len(system_map)] = value
			system_map[value] = true
		}
	}
	// Solve the linear system
	secret := solveLinearSystem(system)
	fmt.Printf("Secret is %d\n", secret)
	os.Exit(0)
}
