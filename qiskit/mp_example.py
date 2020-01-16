#!/usr/bin/env python3

"""Create circuit based on question in Qiskit Slack by Matan Parnas.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)
from qiskit.quantum_info.operators import Operator
import math

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


# Create subcircuit. The first two qubits are assumed to be |00>, and they
# will be turned into a Hardy state 1/sqrt(3)(|00> + |01> + |10>). Then the
# following operation is performed on the third qubit:
# |00> - identity
# |01> - H
# |10> - Rz(π/2)
sub = QuantumCircuit(3, name='{1, H, Rz(π/2)}')
op = Operator([
    [math.sqrt(1/3),  math.sqrt(2/3)],
    [math.sqrt(2/3), -math.sqrt(1/3)]
])
sub.unitary(op, [0])
sub.ch(0, 1)
sub.cx(1, 0)
sub.ch(0, 2)  # H if |01>
sub.crz(math.pi/2, 1, 2)  # Rz if |10>
rando = sub.to_instruction()

# Pick a random number from 1 to 3 and apply one of {1, H, Rz(π/2)}.
qc = QuantumCircuit(3, 2)
qc.append(rando, range(3))
qc.measure(range(2), range(2))
print(qc)

# Execute the circuit on the device.
job = execute(qc, device, shots=1024)

# Get the result counts.
result = job.result()
counts = result.get_counts(qc)
print("\nTotal counts are:", dict(sorted(counts.items())))
