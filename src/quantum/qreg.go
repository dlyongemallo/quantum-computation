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
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// A QReg represents a quantum register.
type QReg struct {
	// The width (number of qubits) of this quantum register.
	width int

	// The complex amplitudes for each of the standard basis states.
	// There are math.Pow(2,width) of these.
	amplitudes []complex128
}

// Constructor for a QReg of the given width. Optionally set its initial
// value to a basis state, specified either as an integer or a series of
// binary digits.
func NewQReg(width int, values ...int) *QReg {
	qreg := &QReg{width, nil}
	qreg.Set(values...)
	return qreg
}

// The eigenvectors of the Pauli matrices are defined here for convenience,
// with the eigenvalue +1 vector first, followed by the eigenvalue -1 vector.
// These are the eigenvectors of the Pauli Z matrix (and the standard basis
// states for one qubit).
func KetZero() *QReg {
	// |0> = [1; 0]
	return &QReg{1, []complex128{1, 0}}
}
func KetOne() *QReg {
	// |1> = [0; 1]
	return &QReg{1, []complex128{0, 1}}
}

// These are the eigenvectors of the Pauli X matrix.
func KetPlus() *QReg {
	// |+> = 1/sqrt{2}[1; 1]
	return &QReg{1, []complex128{1 / math.Sqrt2, 1 / math.Sqrt2}}
}
func KetMinus() *QReg {
	// |-> = 1/sqrt{2}[1; -1]
	return &QReg{1, []complex128{1 / math.Sqrt2, -1 / math.Sqrt2}}
}

// These are the eigenvectors of the Pauli Y matrix.
// There is a bug in go1.0.2 which prevents these from being initialised to
// the correct values if inlining is turned on.
// See http://code.google.com/p/go/issues/detail?id=4159
func KetPlusI() *QReg {
	// |+i> = 1/sqrt{2}[1; i]
	return &QReg{1, []complex128{1 / math.Sqrt2, complex(0, 1/math.Sqrt2)}}
}
func KetMinusI() *QReg {
	// |-i> = 1/sqrt{2}[1; -i]
	return &QReg{1, []complex128{1 / math.Sqrt2, complex(0, -1/math.Sqrt2)}}
}

// Convenience constructor for a qubit, specified by its spherical coordinates
// on the Bloch sphere.
func NewQubitWithBlochCoords(theta, phi float64) *QReg {
	// |psi> = cos(theta/2) + e^{i phi}sin(theta/2)
	t := complex(theta/2, 0)
	p := complex(phi, 0)
	qreg := &QReg{1, []complex128{cmplx.Cos(t),
		cmplx.Exp(complex(0, 1)*p) * cmplx.Sin(t)}}
	return qreg
}

// Accessor for the width of a QReg.
func (qreg *QReg) Width() int {
	return qreg.width
}

// Copy a QReg, for testing purposes only. The no-cloning theorem of course
// prevents copying an actual quantum register in an unknown arbitrary state.
func (qreg *QReg) Copy() *QReg {
	newQreg := &QReg{qreg.width, make([]complex128, len(qreg.amplitudes))}
	copy(newQreg.amplitudes, qreg.amplitudes)
	return newQreg
}

// Compute the probability of observing a state.
func (qreg *QReg) StateProb(state int) float64 {
        // TODO(davinci): Allow this to accept a series of binary values for
        // specifying the state.
        // The probability of observing a state is the square of the magnitude
        // of the complex amplitude.
        magnitude := cmplx.Abs(qreg.amplitudes[state])
	return magnitude * magnitude
}

// Compute the probability of observing a state for a specific bit.
func (qreg *QReg) BProb(index int, value int) float64 {
	prob := float64(0.0)
	bit := 1 << uint(index)
	bitnot := (1 - value) << uint(index)
	// Iterate through all the basis states where this bit is 1
	for state := 0 | bit; state < len(qreg.amplitudes); state = (state + 1) | bit {
		prob += qreg.StateProb(state - bitnot)
	}
	return prob
}

// Set the QReg to a state in the standard basis. If no value is given, default
// to the all zero state. If one value is given, interpret it as the integer
// representation of a basis state. If a series of binary values are given,
// interpret them as the binary representation of a basis state.
func (qreg *QReg) Set(values ...int) {
	// The Hilbert space has dimension math.Pow(2,width).
	hilbertSpaceDim := 1 << uint(qreg.width)

	qreg.amplitudes = make([]complex128, hilbertSpaceDim)
	if len(values) == 0 {
		// Set to |0...0>.
		qreg.amplitudes[0] = 1
	} else if len(values) == 1 {
		// Given an integer d, set to basis state |d>.
		if values[0] < 0 || values[0] >= hilbertSpaceDim {
			errStr := fmt.Sprintf("Value of %d is too large for "+
				"QReg of width %d.", values[0], qreg.width)
			panic(errStr)
		}
		qreg.amplitudes[values[0]] = 1
	} else if len(values) == qreg.width {
		// Given binary b_1, b_2, ..., b_k, set to basis state
		// |b_1 b_2 ... b_k>.
		basisStateIndex := 0
		for _, value := range values {
			basisStateIndex <<= 1
			if value < 0 || value > 1 {
				panic("Expected 0 or 1 when setting value of " +
					"quantum register.")
			}
			basisStateIndex += value
		}
		qreg.amplitudes[basisStateIndex] = 1
	} else {
		panic("Bad values for quantum register.")
	}
}

// Set a particular bit in a QReg
func (qreg *QReg) BSet(index int, value int) {
	if value > 1 {
		errStr := fmt.Sprintf("Value %d should be either 0 or 1",
			value)
		panic(errStr)
	}
	bit := 1 << uint(index)
	bitval := value << uint(index)
	bitnot := (1 - value) << uint(index)
	bprob := qreg.BProb(index, value)
	if bprob > 0 {
		ampFactor := complex(1.0/math.Sqrt(bprob), 0)
		// Alter every state.  If it's the right qubit value, fix the
		// amplitude; otherwise, set the amplitude to 0.
		for state, amp := range qreg.amplitudes {
			if int(state)&bit == bitval {
				qreg.amplitudes[state] = amp * ampFactor
			} else {
				qreg.amplitudes[state] = complex(0, 0)
			}
		}
	} else {
		// Iterate through all the amplitudes where this bit is 1
		for state := int(0) | bit; state < int(len(qreg.amplitudes)); state = (state + 1) | bit {
			// Add the amplitude of the old state to the new state
			oldState := state - bitval
			newState := state - bitnot
			qreg.amplitudes[newState] += qreg.amplitudes[oldState]
			qreg.amplitudes[oldState] = complex(0, 0)
		}
	}
}

// Measure a bit without collapsing its quantum state
func (qreg *QReg) BMeasurePreserve(index int) int {
	if rand.Float64() < qreg.BProb(index, 0) {
		return 0
	}
	return 1
}

// Measure a bit (the quantum state of this qubit will collapse)
func (qreg *QReg) BMeasure(index int) int {
	b := qreg.BMeasurePreserve(index)
	qreg.BSet(index, b)
	return b
}

// Measure a register without collapsing its quantum state
func (qreg *QReg) MeasurePreserve() int {
	r := rand.Float64()
	sum := float64(0.0)
	for i := range qreg.amplitudes {
		sum += qreg.StateProb(i)
		if r < sum {
			return i
		}
	}
	return len(qreg.amplitudes) - 1
}

// Measure a register
func (qreg *QReg) Measure() int {
	value := qreg.MeasurePreserve()
	var amp complex128
	if real(qreg.amplitudes[value]) > 0 {
		amp = complex(1, 0)
	} else {
		amp = complex(-1, 0)
	}
	for i := range qreg.amplitudes {
		qreg.amplitudes[i] = complex(0, 0)
	}
	qreg.amplitudes[value] = amp
	return value
}

func (qreg *QReg) PrintState(index int) {
	prob := qreg.StateProb(index)
	largest := (1 << uint(qreg.width)) - 1
	padding := int(math.Floor(math.Log10(float64(largest)))) + 1
	format := fmt.Sprintf("%%+f%%f|(%%%dd)%%0%db>", padding, qreg.width)
	fmt.Printf(format, qreg.amplitudes[index], prob, index, index)
}

func (qreg *QReg) PrintStateln(index int) {
	qreg.PrintState(index)
	fmt.Println()
}

func (qreg *QReg) Print() {
	for i := range qreg.amplitudes {
		qreg.PrintStateln(i)
	}
}

func (qreg *QReg) PrintNonZero() {
	for i, amplitude := range qreg.amplitudes {
		if amplitude != 0 {
			qreg.PrintStateln(i)
		}
	}
}
