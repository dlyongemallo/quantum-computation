#!/usr/bin/env python3

"""Create a wire.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister)
from qiskit.qasm3 import dumps

# Init 1 qubit and 1 classical bit.
q = QuantumRegister(1)
c = ClassicalRegister(1)
circuit = QuantumCircuit(q, c)

# Do a measurement.
circuit.measure(q, c)
print(circuit.draw())

# Output QASM
print(circuit.qasm())
print(dumps(circuit))
