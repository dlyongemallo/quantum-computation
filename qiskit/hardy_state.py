#!/usr/bin/env python3

"""Create the Hardy state 1/sqrt(3)(|00> + |01> + |10>).
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)
from qiskit.quantum_info.operators import Operator
from math import sqrt

# Set to true to use an actual device.
use_device = False
if use_device:
    IBMQ.load_account()
    provider = IBMQ.get_provider(hub='ibm-q')
    device = least_busy(provider.backends(
        filters=lambda x: x.configuration().n_qubits >= 2 and
        not x.configuration().simulator and x.status().operational==True))
    print("Using backend: ", device)
else:
    device = Aer.get_backend('qasm_simulator')

# Init 2 qubits and 2 classical bits.
circuit = QuantumCircuit(2, 2)
op = Operator([
    [sqrt(1/3),  sqrt(2/3)],
    [sqrt(2/3), -sqrt(1/3)]
])
circuit.unitary(op, [0])
circuit.ch(0, 1)
circuit.cx(1, 0)
circuit.measure([0, 1], [0, 1])
print(circuit.draw())

# Execute the circuit on the device.
job = execute(circuit, device, shots=1024)

# Get the result counts.
result = job.result()
counts = result.get_counts(circuit)
print("\nTotal counts are:", dict(sorted(counts.items())))
