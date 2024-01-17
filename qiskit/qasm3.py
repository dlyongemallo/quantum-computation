#!/usr/bin/env python3

"""Exercise all gates in stdgates.inc.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister)
from qiskit.qasm3 import dumps
import math

# Init 3 qubits.
q = QuantumRegister(3)
circuit = QuantumCircuit(q)

# Do stuff.
circuit.p(math.pi, q[0])

circuit.x(q[0])
circuit.y(q[0])
circuit.z(q[0])

circuit.h(q[0])
circuit.s(q[0])
circuit.sdg(q[0])

circuit.t(q[0])
circuit.tdg(q[0])

circuit.sx(q[0])

circuit.rx(math.pi, q[0])
circuit.ry(math.pi, q[0])
circuit.rz(math.pi, q[0])

circuit.cx(q[0], q[1])
circuit.cy(q[0], q[1])
circuit.cz(q[0], q[1])
circuit.cp(math.pi, q[0], q[1])
circuit.crx(math.pi, q[0], q[1])
circuit.cry(math.pi, q[0], q[1])
circuit.crz(math.pi, q[0], q[1])
circuit.ch(q[0], q[1])

circuit.swap(q[0], q[1])

circuit.ccx(q[0], q[1], q[2])
circuit.cswap(q[0], q[1], q[2])

circuit.cu(math.pi, math.pi, math.pi, math.pi, q[0], q[1])

# circuit.CX(q[0], q[1])
# circuit.phase(math.pi, q[0])
# circuit.cphase(math.pi, q[0], q[1])
circuit.id(q[0])
# circuit.u1(math.pi, q[0])
# circuit.u2(math.pi, math.pi, q[0])
# circuit.u3(math.pi, math.pi, math.pi, q[0])

# Output QASM
print(circuit.qasm())
print(dumps(circuit))
