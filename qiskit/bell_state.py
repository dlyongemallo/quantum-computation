#!/usr/bin/env python3

"""Create a simple Bell state.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)
from qiskit.providers.ibmq import least_busy

# Register with the API.
IBMQ.load_account()
provider = IBMQ.get_provider(hub='ibm-q')

# Set to true to use an actual device.
use_device = False
if use_device:
    device = least_busy(provider.backends(
        filters=lambda x: x.configuration().n_qubits >= 2 and
        not x.configuration().simulator and x.status().operational==True))
    print("Using backend: ", device)
else:
    device = Aer.get_backend('qasm_simulator')

# Init 2 qubits and 2 classical bits.
q = QuantumRegister(2)
c = ClassicalRegister(2)
circuit = QuantumCircuit(q, c)

# Create a Bell state (psi+).
circuit.h(q[0])
circuit.cx(q[0], q[1])
circuit.measure(q, c)
print(circuit.draw())

# Execute the circuit on the device.
job = execute(circuit, device, shots=1000)

# Get the result counts.
result = job.result()
counts = result.get_counts(circuit)
print("\nTotal counts are:", counts)
