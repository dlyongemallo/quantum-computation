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
// Authors: conleyo@google.com (Conley Owens),
//          davinci@google.com (David Yonge-Mallo)

package quantum

import (
	"math"
	"math/cmplx"
)

// Define gates for single-qubit operations.
func newOneQubitGate(arr [4]complex128) *Gate {
	var matrix [2][]complex128
	matrix[0] = arr[0:2]
	matrix[1] = arr[2:4]
	get := func(row, col int) complex128 {
		return matrix[row][col]
	}
	return &Gate{1, get}
}

// Define the gates coresponding to the Pauli matrices.
// The Pauli X gate or NOT gate.
func PauliX() *Gate {
	return newOneQubitGate([4]complex128{
		0, 1,
		1, 0})
}

// The Pauli Y gate.
func PauliY() *Gate {
	return newOneQubitGate([4]complex128{
		0, complex(0, -1),
		complex(0, 1), 0})
}

// The Pauli Z gate.
func PauliZ() *Gate {
	return newOneQubitGate([4]complex128{
		1,  0,
		0, -1})
}

// Define the arbitrary rotation gates.
// Rotation about the X-axis.
func RotationX(theta float64) *Gate {
	// R_x(theta) = cos(theta/2)I - i sin(theta/2)X
	t := theta/2
	cos := complex(math.Cos(t), 0)
	nisin := complex(0, -1 * math.Sin(t))
	return newOneQubitGate([4]complex128{
		cos,   nisin,
		nisin, cos})
}

// Rotation about the Y-axis.
func RotationY(theta float64) *Gate {
	// R_y(theta) = cos(theta/2)I - i sin(theta/2)Y
	t := theta/2
	cos := complex(math.Cos(t), 0)
	nsin := complex(-1 * math.Sin(t), 0)
	sin := complex(math.Sin(t), 0)
	return newOneQubitGate([4]complex128{
		cos, nsin,
		sin, cos})
}

// Rotation about the Z-axis.
func RotationZ(theta float64) *Gate {
	// R_z(theta) = cos(theta/2)I - i sin(theta/2)Z
	t := theta/2
	nexp := cmplx.Exp(complex(0, -1 * t))
	exp := cmplx.Exp(complex(0, t))
	return newOneQubitGate([4]complex128{
		nexp, 0,
		0,    exp})
}

// Hadamard Gate

func NewHadamardGate(width int) *Gate {
	d := float64(int(1 << uint(width>>1)))
	if width&1 == 1 {
		d *= math.Sqrt2
	}
	p := complex(1.0/d, 0)
	n := -p
	return NewFuncGateNoCheck(
		// get(row, col)
		func(row int, col int) complex128 {
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
		width)
}

func Hadamard(qreg *QReg, target int) {
	NewHadamardGate(1).Apply(qreg, []int{target})
}

func HadamardRange(qreg *QReg, targetRangeStart int, targetRangeEnd int) {
	targetRangeSize := targetRangeEnd - targetRangeStart
	gate := NewHadamardGate(targetRangeSize)
	gate.ApplyRange(qreg, targetRangeStart)
}

func HadamardReg(qreg *QReg) {
	HadamardRange(qreg, 0, qreg.width)
}

// Diffusion Gate

func NewDiffusionGate(width int) *Gate {
	a2 := complex(2.0/float64(int(1<<uint(width))), 0)
	a2m1 := a2 - complex(1.0, 0)
	return NewFuncGate(func(row int, col int) complex128 {
		if row == col {
			return a2m1
		}
		return a2
	},
		width)
}

func Diffusion(qreg *QReg, target int) {
	NewDiffusionGate(1).Apply(qreg, []int{target})
}

func DiffusionRange(qreg *QReg, targetRangeStart int, targetRangeEnd int) {
	targetRangeSize := targetRangeEnd - targetRangeStart
	gate := NewDiffusionGate(targetRangeSize)
	gate.ApplyRange(qreg, targetRangeStart)
}

func DiffusionReg(qreg *QReg) {
	DiffusionRange(qreg, 0, qreg.width)
}
