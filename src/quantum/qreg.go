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

// Set the QReg to a state in the standard basis. If no value is given, default
// to the all zero state. If one value is given, interpret it as the integer
// representation of a basis state. If a series of binary values are given,
// interpret them as the binary representation of a basis state.
func (qreg *QReg) Set(values ...int) {
	// The Hilbert space has dimension math.Pow(2,width).
	qreg.amplitudes = make([]complex128, 1<<uint(qreg.width))
	qreg.amplitudes[qreg.basisStateLabel(values...)] = 1
}

// Given a series of values, convert them into the label of a standard basis
// for the quantum register. If no values are given, return the label of the
// all-zero state |0...0>. If one value is given, interpret it as the integer
// (decimal) representation of a basis state label. If a series of binary
// values are given equal in number to width, interpret them as the binary
// representation of such a label.
func (qreg *QReg) basisStateLabel(values ...int) int {
	if len(values) == 0 {
		// The basis state is |0>.
		return 0
	} else if len(values) == 1 {
		// Given integer d, the basis state is |d>.
		if values[0] < 0 || values[0] >= len(qreg.amplitudes) {
			errStr := fmt.Sprintf("The state |%d> is not "+
				"possible for a register of width %d.",
				values[0], qreg.width)
			panic(errStr)
		}
		return values[0]
	} else if len(values) == qreg.width {
		// Given binary {b_1, b_2, ..., b_k}, the basis state is
		// |b_1 b_2 ... b_k>.
		label := 0
		for _, value := range values {
			label <<= 1
			if value < 0 || value > 1 {
				panic("Unexpected non-binary value in label " +
					"of quantum register.")
			}
			label += value
		}
		return label
	}
	panic("Bad label for quantum register.")
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

// Compute the probability of observing a basis state.
func (qreg *QReg) StateProb(values ...int) float64 {
	label := qreg.basisStateLabel(values...)
	// The probability of observing a state is the square of the magnitude
	// of the complex amplitude.
	magnitude := cmplx.Abs(qreg.amplitudes[label])
	return magnitude * magnitude
}

// Compute the probability of observing a specific bit for a basis state. The
// individual bits of the register are indexed starting with 0 on the left
// (i.e., most significant bit). A pair is returned corresponding to the
// probabilities of observing 0 and 1, respectively.
func (qreg *QReg) BProb(bitIndex int) [2]float64 {
	prob := float64(0.0)
	bitMask := 1 << uint(qreg.width-1-bitIndex)
	// Iterate through all the basis states where the indexed bit is 1, to
	// sum the probability of observing 1 for that bit.
	for label := 0 | bitMask; label < len(qreg.amplitudes); label = (label + 1) | bitMask {
		prob += qreg.StateProb(label)
	}
	return [2]float64{1.0 - prob, prob}
}

// Set a particular bit in a QReg. This should normally not be called directly,
// except for testing purposes, since it is not a physically realistic operation.
func (qreg *QReg) BSet(bitIndex int, value int) {
	if value < 0 || value > 1 {
		errStr := fmt.Sprintf("Value %d should be either 0 or 1.",
			value)
		panic(errStr)
	}
	bitMask := 1 << uint(qreg.width-1-bitIndex)
	valueMask := value << uint(qreg.width-1-bitIndex)
	bprob := qreg.BProb(bitIndex)[value]
	if bprob > 0 {
		// Go through the amplitudes associated with each basis state
		// label. If it has the bit set to the right value, scale the
		// amplitude appropriately; otherwise, set the amplitude to 0.
		newDenominator := complex(math.Sqrt(bprob), 0)
		for label := range qreg.amplitudes {
			amplitude := &qreg.amplitudes[label]
			if label&bitMask == valueMask {
				*amplitude /= newDenominator
			} else {
				*amplitude = complex(0, 0)
			}
		}
	} else {
		// This case should never happen. If the probability of reading the
		// given bit value is 0, this should result in an error.
		// TODO(davinci): Investigate why it is here, and if possible remove it.

		// Iterate through all the basis states where this bit is 1
		notValueMask := (1 - value) << uint(qreg.width-1-bitIndex)
		for label := 0 | bitMask; label < len(qreg.amplitudes); label = (label + 1) | bitMask {
			// Add the amplitude of the old label to the new label
			oldLabel := label - valueMask
			newLabel := label - notValueMask
			qreg.amplitudes[newLabel] += qreg.amplitudes[oldLabel]
			qreg.amplitudes[oldLabel] = complex(0, 0)
		}
	}
}

// Simulate a measurement on a bit, i.e., get the result of the measurement
// but without collapsing its quantum state.
func (qreg *QReg) BMeasurePreserve(bitIndex int) int {
	if rand.Float64() < qreg.BProb(bitIndex)[0] {
		return 0
	}
	return 1
}

// Measure a bit (the quantum state of this qubit will collapse).
func (qreg *QReg) BMeasure(bitIndex int) int {
	b := qreg.BMeasurePreserve(bitIndex)
	qreg.BSet(bitIndex, b)
	return b
}

// Simulate a measurement on a register, i.e., get the result of the measurement
// bit without collapsing its quantum state.
func (qreg *QReg) MeasurePreserve() int {
	r := rand.Float64()
	sum := float64(0.0)
	for label := range qreg.amplitudes {
		sum += qreg.StateProb(label)
		if r < sum {
			return label
		}
	}
	return len(qreg.amplitudes) - 1
}

// Measure a register.
func (qreg *QReg) Measure() int {
	outputLabel := qreg.MeasurePreserve()

	// We rescale the amplitude corresponding to the output basis state so
	// that its probability is 1. However, we preserve the (global) phase,
	// since the register may be a subsystem in a larger system.
	amplitude := qreg.amplitudes[outputLabel]
	amplitude /= complex(cmplx.Abs(amplitude), 0)
	for label := range qreg.amplitudes {
		qreg.amplitudes[label] = complex(0, 0)
	}
	qreg.amplitudes[outputLabel] = amplitude
	return outputLabel
}

func (qreg *QReg) PrintState(label int) {
	prob := qreg.StateProb(label)
	largest := (1 << uint(qreg.width)) - 1
	padding := int(math.Floor(math.Log10(float64(largest)))) + 1
	format := fmt.Sprintf("%%+f%%f|(%%%dd)%%0%db>", padding, qreg.width)
	fmt.Printf(format, qreg.amplitudes[label], prob, label, label)
}

func (qreg *QReg) PrintStateln(label int) {
	qreg.PrintState(label)
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
