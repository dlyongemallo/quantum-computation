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
	"testing"
)

// Threshold for how close two probabilities or complex amplitudes have to be
// before they're considered equal.
const threshold = 0.0000000001

// Helper function for testing. Returns true if the amplitude for the given
// basis state is set to 1, and all other amplitudes are set to 0.
func isBasisState(qreg *QReg, basis int) bool {
	for i, amplitude := range qreg.amplitudes {
		if amplitude != complex(0, 0) && i != basis {
			return false
		}
	}
	return qreg.amplitudes[basis] == complex(1, 0)
}

// Helper function to test that two probabilities are "close enough".
func verifyProb(expected, actual float64) bool {
	return math.Abs(actual-expected) < threshold
}

// Helper function to test that two complex amplitudes are "close enough".
func verifyAmplitude(expected, actual complex128) bool {
	return math.Abs(cmplx.Abs(actual)-cmplx.Abs(expected)) < threshold
}

// Test the various forms of the constructor.
func TestNewQReg(t *testing.T) {
	// Test constructor that takes in no initial values.
	qreg := NewQReg(4)
	if !isBasisState(qreg, 0) {
		t.Error("Expected |0000>.")
	}

	// Test constructor that takes in integer representation of basis state.
	qreg = NewQReg(8, 3)
	if !isBasisState(qreg, 3) {
		t.Error("Expected |00000011>.")
	}

	// Test constructor that takes in binary representation of basis state.
	qreg = NewQReg(5, 0, 1, 1, 0, 1)
	if !isBasisState(qreg, 13) {
		t.Error("Expected |01101>.")
	}
}

// Test that the correct values are computed for the probability of observing
// a basis state.
func TestQRegStateProb(t *testing.T) {
	// TODO(davinci): Add more tests.

	// Test the |+> state.
	qreg := KetPlus()
	if !verifyProb(float64(0.5), qreg.StateProb(0)) {
		t.Errorf("Bad probability for state |+> = %+f, expected 0.5.",
			qreg.StateProb(0))
	}
	if !verifyProb(float64(0.5), qreg.StateProb(1)) {
		t.Errorf("Bad probability for state |-> = %+f, expected 0.5.",
			qreg.StateProb(1))
	}
}

// Test that the correct values are computed for the probability of observing
// a given bit.
func TestQRegBProb(t *testing.T) {
	qreg := NewQReg(2)
	qreg.amplitudes = []complex128{
		cmplx.Sqrt(0.1), // |00>
		cmplx.Sqrt(0.2), // |01>
		cmplx.Sqrt(0.3), // |10>
		cmplx.Sqrt(0.4), // |10>
	}

	// |00> and |01>
	if !verifyProb(qreg.BProb(0)[0], 0.3) {
		t.Errorf("Bad probability for |0?>, expected 0.3.")
	}

	// |10> and |11>
	if !verifyProb(qreg.BProb(0)[1], 0.7) {
		t.Errorf("Bad probability for |1?>, expected 0.7.")
	}

	// |00> and |10>
	if !verifyProb(qreg.BProb(1)[0], 0.4) {
		t.Errorf("Bad probability for |?0>, expected 0.4.")
	}

	// |01> and |11>
	if !verifyProb(qreg.BProb(1)[1], 0.6) {
		t.Errorf("Bad probability for |?1>, expected 0.6.")
	}
}

func TestQRegBSet_1BitCollapsed(t *testing.T) {
	qreg := NewQReg(1, 0)
	qreg.BSet(0, 1)
	if qreg.amplitudes[0] != complex(0, 0) {
		t.Errorf("Bad amplitude for state 0 = %+f, expected 0.",
			qreg.amplitudes[0])
	}
	if qreg.amplitudes[1] != complex(1, 0) {
		t.Errorf("Bad amplitude for state 1 = %+f, expected 1.",
			qreg.amplitudes[1])
	}
}

func TestQRegBSet_1BitEntangled(t *testing.T) {
	qreg := NewQReg(2, 0)
	qreg.amplitudes[0] = complex(1/math.Sqrt2, 0)
	qreg.amplitudes[1] = complex(-1/math.Sqrt2, 0)
	qreg.BSet(1, 1)
	if !verifyAmplitude(complex(0, 0), qreg.amplitudes[0]) {
		t.Errorf("Bad amplitude for state 0 = %+f, expected 0.",
			qreg.amplitudes[0])
	}
	if !verifyAmplitude(complex(-1, 0), qreg.amplitudes[1]) {
		t.Errorf("Bad amplitude for state 1 = %+f, expected -1.",
			qreg.amplitudes[1])
	}
}
