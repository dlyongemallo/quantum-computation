#!/usr/bin/env python3

"""Implementation of Deutsch algorithm (one qubit).

The Cirq version is here:
https://github.com/quantumlib/Cirq/blob/master/examples/deutsch.py
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, Aer)
from qiskit.qasm3 import dumps
from qiskit_ibm_provider import IBMProvider, least_busy
import random

# Set to true to use an actual device.
use_device = False
if use_device:
    provider = IBMProvider(instance='ibm-q/open/main')
    device = least_busy(provider.backends(
        filters=lambda x: x.configuration().n_qubits >= 2 and
        not x.configuration().simulator and x.status().operational==True))
    print("Using backend: ", device)
else:
    device = Aer.get_backend('qasm_simulator')

# Pick a secret function.
secret = [random.randint(0, 1) for _ in range(2)]
def append_oracle(ciruit, secret):
    if secret[0]:
        circuit.cx(0, 1)
        circuit.x(1)
    if secret[1]:
        circuit.cx(0, 1)

# Init 2 qubits and 1 classical bits.
circuit = QuantumCircuit(2, 1)

# Create the Deutsch algorithm circuit.
circuit.x(1)
circuit.barrier()
circuit.h(0)
circuit.h(1)
circuit.barrier()
append_oracle(circuit, secret)
circuit.barrier()
circuit.h(0)
circuit.measure(0, 0)

print("f(0) = {:d}, f(1) = {:d}".format(secret[0], secret[1]))
print(circuit)

# Execute the circuit on the device.
job = execute(circuit, device, shots=1024)

# Get the result counts.
result = job.result()
counts = result.get_counts(circuit)
print("\nTotal counts are:", dict(sorted(counts.items())))

# Output QASM
print(circuit.qasm())
print(dumps(circuit))
