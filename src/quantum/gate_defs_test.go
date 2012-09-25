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
	"testing"
)

func TestHadamardGate_2x2(t *testing.T) {
	matrix := NewHadamardGate(1)
	if matrix.width() != 2 {
		t.Errorf("Bad width in 2x2 Hadamard")
	}
	p := complex(1.0/math.Sqrt(2), 0)
	n := -p
	arr := []complex128{
		p, p,
		p, n,
	}
	for i := 0; i < 4; i++ {
		a := i / 2
		b := i % 2
		v := matrix.get(a, b)
		if v != arr[i] {
			t.Errorf("Bad value in 2x2 Hadamard matrix at index "+
				"%d, %d = %f; want %f",
				a, b, v, arr[i])
		}
	}
}

func TestHadamardGate_4x4(t *testing.T) {
	matrix := NewHadamardGate(2)
	if matrix.width() != 4 {
		t.Errorf("Bad width in 4x4 Hadamard")
	}
	p := complex(.5, 0)
	n := -p
	arr := []complex128{
		p, p, p, p,
		p, n, p, n,
		p, p, n, n,
		p, n, n, p,
	}
	for i := 0; i < 16; i++ {
		a := i / 4
		b := i % 4
		v := matrix.get(a, b)
		if v != arr[i] {
			t.Errorf("Bad value in 4x4 Hadamard matrix at index "+
				"%d, %d = %f; want %f",
				a, b, v, arr[i])
		}
	}
}

func TestDiffusionGate_2x2(t *testing.T) {
	matrix := NewDiffusionGate(1)
	if matrix.width() != 2 {
		t.Errorf("Bad width in 2x2 Diffusion")
	}
	p := complex(1.0, 0)
	n := complex(0.0, 0)
	arr := []complex128{
		n, p,
		p, n,
	}
	for i := 0; i < 4; i++ {
		a := i / 2
		b := i % 2
		v := matrix.get(a, b)
		if v != arr[i] {
			t.Errorf("Bad value in 2x2 Diffusion matrix at index "+
				"%d, %d = %f; want %f",
				a, b, v, arr[i])
		}
	}
}

func TestDiffusionGate_4x4(t *testing.T) {
	matrix := NewDiffusionGate(2)
	if matrix.width() != 4 {
		t.Errorf("Bad width in 4x4 Diffusion")
	}
	p := complex(0.5, 0)
	n := complex(-0.5, 0)
	arr := []complex128{
		n, p, p, p,
		p, n, p, p,
		p, p, n, p,
		p, p, p, n,
	}
	for i := 0; i < 16; i++ {
		a := i / 4
		b := i % 4
		v := matrix.get(a, b)
		if v != arr[i] {
			t.Errorf("Bad value in 4x4 Diffusion matrix at index "+
				"%d, %d = %f; want %f",
				a, b, v, arr[i])
		}
	}
}
