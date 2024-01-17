#!/usr/bin/env python3

"""Verify qiskit's handling of qasm2 rz gate.

See: https://github.com/Qiskit/qiskit/issues/10790
"""

from qiskit.circuit import QuantumCircuit
from qiskit import quantum_info

s1 = """
OPENQASM 2.0;
include "qelib1.inc";
qreg q[1];
rz(pi/2) q[0];
"""
c1 = QuantumCircuit.from_qasm_str(s1)
print(quantum_info.Operator(c1).data)
print()

s2 = """
OPENQASM 2.0;
include "qelib1.inc";
qreg q[2];
rz(pi/4) q[0];
cx q[1],q[0];
rz(-pi/4) q[0];
cx q[1],q[0];
"""
c2 = QuantumCircuit.from_qasm_str(s2)
print(quantum_info.Operator(c2).data)
print()

s3 = """
OPENQASM 2.0;
include "qelib1.inc";
qreg q[1];
u1(pi/2) q[0];
"""
c3 = QuantumCircuit.from_qasm_str(s3)
print(quantum_info.Operator(c3).data)
print()

s4 = """
OPENQASM 2.0;
include "qelib1.inc";
qreg q[2];
u1(pi/4) q[0];
cx q[1],q[0];
u1(-pi/4) q[0];
cx q[1],q[0];
"""
c4 = QuantumCircuit.from_qasm_str(s4)
print(quantum_info.Operator(c4).data)
print()

print(c1.qasm())
print(c2.qasm())
print(c3.qasm())
print(c4.qasm())
