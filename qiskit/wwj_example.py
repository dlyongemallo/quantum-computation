#!/usr/bin/env python3

"""Create circuit posted to Qiskit Slack.
https://algassert.com/quirk#circuit={%22cols%22:[[%22Y%22],[%22inputA1%22,%22Y^(A/2^n)%22],[%22Y^(A/2^n)%22,%22inputA1%22]]}
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)
from qiskit.extensions import CyGate
import math
import numpy as np

# Register with the API.
IBMQ.load_account()
provider = IBMQ.get_provider(hub='ibm-q')

device = Aer.get_backend('unitary_simulator')

# Create subcircuit.
sub_q = QuantumRegister(2)
sub = QuantumCircuit(sub_q, name='ctrl-Y^{A/2}')
sub.cry(math.pi/2, sub_q[0], sub_q[1])
sub.rz(math.pi/4, sub_q[0])
c_sqrt_y_half = sub.to_instruction()

# Init 2 qubits.
q = QuantumRegister(2)
qc = QuantumCircuit(q)

# Construct circuit.
qc.y(q[0])
qc.append(c_sqrt_y_half, [q[0], q[1]])
qc.append(c_sqrt_y_half, [q[1], q[0]])
print(qc)

job = execute(qc, device)
result = job.result()
print(np.around(result.get_unitary(qc),3))
