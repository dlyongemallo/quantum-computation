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

// Represents a quantum register
type QReg struct {
        // The width (number of qubits) of this quantum register.
	width      int

        // The complex amplitudes for each of the standard basis states.
        // There are math.Pow(2,width) of these.
	amplitudes []complex128
}

// Constructor for a QReg
func NewQReg(width int, values ...int) *QReg {
	qreg := &QReg{width, nil}
	qreg.Set(values...)
	return qreg
}

// Accessor for the width of a QReg
func (qreg *QReg) Width() int {
	return qreg.width
}

// Copy a quantum register
func (qreg *QReg) Copy() *QReg {
	new_qreg := &QReg{qreg.width, make([]complex128, len(qreg.amplitudes))}
	copy(new_qreg.amplitudes, qreg.amplitudes)
	return new_qreg
}

// Get the probability of observing a state
func (qreg *QReg) StateProb(state int) float64 {
	return cmplx.Abs(qreg.amplitudes[state] * qreg.amplitudes[state])
}

// Get the probability of observing a state for a specific bit
func (qreg *QReg) BProb(index int, value int) float64 {
	prob := float64(0.0)
	bit := 1 << uint(index)
	bitnot := (1 - value) << uint(index)
	// Iterate through all the amplitudes where this bit is 1
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
        hilbert_space_dim := 1<<uint(qreg.width)

        qreg.amplitudes = make([]complex128, hilbert_space_dim)
        if len(values) == 0 {
                // Set to |0...0>.
                qreg.amplitudes[0] = 1
        } else if len(values) == 1 {
                // Given an integer d, set to basis state |d>.
                if values[0] < 0 || values[0] >= hilbert_space_dim {
                        err_str := fmt.Sprintf("Value of %d is too large for " +
                                "QReg of width %d.", values[0], qreg.width)
                        panic(err_str)
                }
                qreg.amplitudes[values[0]] = 1
        } else if len(values) == qreg.width {
                // Given binary b_1, b_2, ..., b_k, set to basis state
                // |b_1 b_2 ... b_k>.
                basis_state_index := 0
                for _, value := range values {
                        basis_state_index <<= 1
                        if value < 0 || value > 1 {
                                panic("Expected 0 or 1 when setting value of " +
                                        "quantum register.")
                        }
                        basis_state_index += value
                }
                qreg.amplitudes[basis_state_index] = 1
        } else {
                panic("Bad values for quantum register.")
        }
}

// Set a particular bit in a QReg
func (qreg *QReg) BSet(index int, value int) {
	if value > 1 {
		err_str := fmt.Sprintf("Value %d should be either 0 or 1",
			value)
		panic(err_str)
	}
	bit := 1 << uint(index)
	bitval := value << uint(index)
	bitnot := (1 - value) << uint(index)
	bprob := qreg.BProb(index, value)
	if bprob > 0 {
		amp_factor := complex(1.0/math.Sqrt(bprob), 0)
		// Alter every state.  If it's the right qubit value, fix the
		// amplitude; otherwise, set the amplitude to 0.
		for state, amp := range qreg.amplitudes {
			if int(state)&bit == bitval {
				qreg.amplitudes[state] = amp * amp_factor
			} else {
				qreg.amplitudes[state] = complex(0, 0)
			}
		}
	} else {
		// Iterate through all the amplitudes where this bit is 1
		for state := int(0) | bit; state < int(len(qreg.amplitudes)); state = (state + 1) | bit {
			// Add the amplitude of the old state to the new state
			old_state := state - bitval
			new_state := state - bitnot
			qreg.amplitudes[new_state] += qreg.amplitudes[old_state]
			qreg.amplitudes[old_state] = complex(0, 0)
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
	for i, _ := range qreg.amplitudes {
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
	for i, _ := range qreg.amplitudes {
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
	for i, _ := range qreg.amplitudes {
		qreg.PrintStateln(i)
	}
}

func (qreg *QReg) PrintNonZero() {
	for i, state := range qreg.amplitudes {
		if state != 0 {
			qreg.PrintStateln(i)
		}
	}
}
