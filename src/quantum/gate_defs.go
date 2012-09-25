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

package quantum

import (
	"math"
)

// Hadamard Gate

func NewHadamardGate(bits int) *Gate {
	d := float64(int(1 << uint(bits>>1)))
	if bits&1 == 1 {
		d *= math.Sqrt2
	}
	p := complex(1.0/d, 0)
	n := -p
	return NewFuncGateNoCheck(func(row int, col int) complex128 {
		// Calculate (-1)**<i,j> / sqrt(2**n)
		par := 0
		for anded := row & col; anded > 0; anded >>= 1 {
			par ^= anded & 1
		}
		if par == 1 {
			return n
		}
		return p
	},
		bits)
}

func Hadamard(qreg *QReg, target int) {
	NewHadamardGate(1).Apply(qreg, []int{target})
}

func HadamardRange(qreg *QReg, target_range_start int, target_range_end int) {
	target_range_size := target_range_end - target_range_start
	gate := NewHadamardGate(target_range_size)
	gate.ApplyRange(qreg, target_range_start)
}

func HadamardReg(qreg *QReg) {
	HadamardRange(qreg, 0, qreg.width)
}

// Diffusion Gate

func NewDiffusionGate(bits int) *Gate {
	a2 := complex(2.0/float64(int(1<<uint(bits))), 0)
	a2m1 := a2 - complex(1.0, 0)
	return NewFuncGate(func(row int, col int) complex128 {
		if row == col {
			return a2m1
		}
		return a2
	},
		bits)
}

func Diffusion(qreg *QReg, target int) {
	NewDiffusionGate(1).Apply(qreg, []int{target})
}

func DiffusionRange(qreg *QReg, target_range_start int, target_range_end int) {
	target_range_size := target_range_end - target_range_start
	gate := NewDiffusionGate(target_range_size)
	gate.ApplyRange(qreg, target_range_start)
}

func DiffusionReg(qreg *QReg) {
	DiffusionRange(qreg, 0, qreg.width)
}
