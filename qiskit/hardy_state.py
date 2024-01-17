#!/usr/bin/env python3

"""Create the Hardy state 1/sqrt(3)(|00> + |01> + |10>).
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer, transpile)
from qiskit.qasm3 import dumps
from qiskit.quantum_info.operators import Operator
from qiskit.circuit.library.standard_gates import get_standard_gate_name_mapping
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

# Output QASM2
print(circuit.qasm())

# Output QASM3
# Note: This requires a workaround because of qiskit issue #11558.
# print(dumps(circuit))
basis_gates = list(get_standard_gate_name_mapping())
print(dumps(transpile(circuit, basis_gates=basis_gates, optimization_level=0)))
