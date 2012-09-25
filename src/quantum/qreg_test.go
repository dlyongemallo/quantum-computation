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
	"testing"
)

// Helper function for testing. Returns true if the amplitude for the given
// basis state is set to 1, and all other amplitudes are set to 0.
func verifyBasisState(qreg *QReg, basis int) bool {
        for i, amplitude := range qreg.amplitudes {
                if amplitude != complex(0, 0) && i != basis {
                        return false
                }
        }
        return qreg.amplitudes[basis] == complex(1, 0)
}

// Test the various forms of the constructor.
func TestNewQReg(t *testing.T) {
        // Test constructor that takes in no initial values.
        qreg := NewQReg(4)
        if !verifyBasisState(qreg, 0) {
                t.Error("Expected |0000>.")
        }

        // Test constructor that takes in integer representation of basis state.
        qreg = NewQReg(8, 3)
        if !verifyBasisState(qreg, 3) {
                t.Error("Expected |00000011>.")
        }

        // Test constructor that takes in binary representation of basis state.
        qreg = NewQReg(5, 0, 1, 1, 0, 1)
        if !verifyBasisState(qreg, 13) {
                t.Error("Expected |01101>.")
        }
}

func TestQRegBSet_1BitCollapsed(t *testing.T) {
	qreg := NewQReg(1, 0)
	qreg.BSet(0, 1)
	if qreg.amplitudes[0] != complex(0, 0) {
		t.Errorf("Bad amplitude for state 0 = %+f, want 0",
			qreg.amplitudes[0])
	}
	if qreg.amplitudes[1] != complex(1, 0) {
		t.Errorf("Bad amplitude for state 1 = %+f, want 1",
			qreg.amplitudes[1])
	}
}

func TestQRegBSet_1BitEntangled(t *testing.T) {
	qreg := NewQReg(2, 0)
	qreg.amplitudes[0] = complex(1/math.Sqrt2, 0)
	qreg.amplitudes[1] = complex(-1/math.Sqrt2, 0)
	qreg.BSet(0, 1)
	if qreg.amplitudes[0] != complex(0, 0) {
		t.Errorf("Bad amplitude for state 0 = %+f, want 0",
			qreg.amplitudes[0])
	}
	if qreg.amplitudes[1] != complex(-1, 0) {
		t.Errorf("Bad amplitude for state 1 = %+f, want -1",
			qreg.amplitudes[1])
	}
}
