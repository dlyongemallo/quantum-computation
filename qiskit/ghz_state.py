#!/usr/bin/env python3

"""Create a GHZ state.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)

# Set to true to use an actual device.
use_device = False
if use_device:
    IBMQ.load_account()
    provider = IBMQ.get_provider(hub='ibm-q')
    device = least_busy(provider.backends(
        filters=lambda x: x.configuration().n_qubits >= 3 and
        not x.configuration().simulator and x.status().operational==True))
    print("Using backend: ", device)
else:
    device = Aer.get_backend('qasm_simulator')

# Init 2 qubits and 2 classical bits.
q = QuantumRegister(3)
c = ClassicalRegister(3)
circuit = QuantumCircuit(q, c)

# Create a GHZ state.
circuit.h(q[0])
circuit.cx(q[0], q[1])
circuit.cx(q[0], q[2])
circuit.measure(q, c)
print(circuit.draw())

# Execute the circuit on the device.
job = execute(circuit, device, shots=1024)

# Get the result counts.
result = job.result()
counts = result.get_counts(circuit)
print("\nTotal counts are:", dict(sorted(counts.items())))
