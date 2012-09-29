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
	"fmt"
	"math"
	"math/cmplx"
)

func closeEnough(a complex128, b complex128) bool {
	return math.Abs(cmplx.Abs(a)-cmplx.Abs(b)) < .0000000001
}

type Gate struct {
	// The matrix elements of the gate.
	get   func(row int, col int) complex128

	// The dimension of the Hilbert space.
	// This is equal to math.Pow(2, "bits").
	width func() int
	
	// The width of the gate (number of qubits).
	bits  func() int
}

func (gate *Gate) computeSquareElement(row int, col int, c chan bool) {
	sum := complex(0, 0)
	for i := 0; i < gate.width(); i++ {
		sum += gate.get(row, i) * gate.get(i, col)
	}
	if row == col {
		if closeEnough(sum, complex(1, 0)) {
			c <- false
			return
		}
	} else if closeEnough(sum, complex(0, 0)) {
		c <- false
		return
	}
	c <- true
}

// This tells us whether or not a gate is unitary (it should always be)
func (gate *Gate) IsUnitary() bool {
	c := make(chan bool)
	for row := 0; row < gate.width(); row++ {
		for col := 0; col < gate.width(); col++ {
			go gate.computeSquareElement(row, col, c)
		}
	}
	for i := 0; i < gate.width()*gate.width(); i++ {
		if <-c {
			return false
		}
	}
	return true
}

func NewFuncGateNoCheck(f func(row int, col int) complex128, bits int) *Gate {
	return &Gate{
		// get(row, col)
		f,
		// width()
		func() int {
			return 1 << uint(bits)
		},
		// bits()
		func() int {
			return bits
		}}
}

func NewFuncGate(f func(row int, col int) complex128, bits int) *Gate {
	gate := NewFuncGateNoCheck(f, bits)
	if !gate.IsUnitary() {
		panic("Gate is not unitary")
	}
	return gate
}

func NewArrayGate(arr []complex128) *Gate {
	width := int(math.Sqrt(float64(len(arr))))
	return NewFuncGate(
		// get(row, col)
		func(row int, col int) complex128 {
			return arr[row*width+col]
		},
		// bits()
		int(math.Log2(float64(width))))
}

func NewRealArrayGate(arr []float64) *Gate {
	newArr := make([]complex128, len(arr))
	for i, a := range arr {
		newArr[i] = complex(a, 0)
	}
	return NewArrayGate(newArr)
}

func NewClassicalGate(f func(x int) int, bits int) *Gate {
	return NewFuncGate(func(row int, col int) complex128 {
		if f(col) == row {
			return complex(1, 0)
		}
		return complex(0, 0)
	},
		bits)
}

func stateIndexForTarget(application int, targetValue int, size int, targets []int) int {
	stateVector := make([]int, size)
	for i := 0; i < size; i++ {
		stateVector[i] = 2
	}
	for i := 0; i < len(targets); i++ {
		stateVector[targets[i]] = (targetValue >> uint(i)) & 1
	}
	appPos := 0
	for i := 0; i < size; i++ {
		if stateVector[i] == 2 {
			stateVector[i] = (application >> uint(appPos)) & 1
			appPos++
		}
	}
	index := 0
	for i := 0; i < size; i++ {
		index += stateVector[i] << uint(i)
	}
	return index
}

type indexAmplitude struct {
	index     int
	amplitude complex128
}

// Compute one row of matrix multiplication
func (gate *Gate) computeRow(qreg *QReg, app int, row int, targets []int, c chan indexAmplitude) {
	sum := complex128(complex(0, 0))
	for col := 0; col < gate.width(); col++ {
		index := stateIndexForTarget(app, col, qreg.width, targets)
		sum += gate.get(row, col) * qreg.amplitudes[index]
	}
	index := stateIndexForTarget(app, row, qreg.width, targets)
	c <- indexAmplitude{index, sum}
}

// Apply an arbitrary matrix to a quantum register
// len(matrix) == 4 ** len(targets)
func (gate *Gate) Apply(qreg *QReg, targets []int) {
	// Verify that all the targets are valid
	for _, target := range targets {
		if target >= qreg.width {
			panic(fmt.Sprintf("%d is not a valid target", target))
		}
	}

	numApps := 1 << uint(qreg.width-len(targets))
	newAmplitudes := make([]complex128, len(qreg.amplitudes))
	// Each application of the matrix
	// app is the binary representation of the non-target states
	for app := 0; app < numApps; app++ {
		// Each row of the matrix
		c := make(chan indexAmplitude)
		for row := 0; row < gate.width(); row++ {
			go gate.computeRow(qreg, app, row, targets, c)
		}
		for row := 0; row < gate.width(); row++ {
			ia := <-c
			newAmplitudes[ia.index] = ia.amplitude
		}
	}
	qreg.amplitudes = newAmplitudes
}

func (gate *Gate) ApplyRange(qreg *QReg, targetRangeStart int) {
	targets := make([]int, gate.bits())
	for i := 0; i < gate.bits(); i++ {
		targets[i] = targetRangeStart + i
	}
	gate.Apply(qreg, targets)
}

func (gate *Gate) ApplyReg(qreg *QReg) {
	gate.ApplyRange(qreg, 0)
}

func (gate *Gate) Print() {
	// Get column sizes
	sizes := make([]int, gate.width())
	for col := 0; col < gate.width(); col++ {
		max := 0
		for row := 0; row < gate.width(); row++ {
			l := len(fmt.Sprintf("%+f", gate.get(row, col)))
			if l > max {
				max = l
			}
		}
		if col != 0 {
			max++
		}
		sizes[col] = max
	}
	// Print each row
	for row := 0; row < gate.width(); row++ {
		for col := 0; col < gate.width(); col++ {
			str := fmt.Sprintf("%+f", gate.get(row, col))
			for i := len(str); i < sizes[col]; i++ {
				fmt.Print(" ")
			}
			fmt.Print(str)
		}
		fmt.Println()
	}
}
