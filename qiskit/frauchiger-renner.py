#!/usr/bin/env python3

"""Implementation of the Frauchiger-Renner thought experiment.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, Aer)
from qiskit_ibm_provider import IBMProvider, least_busy
import math
import numpy as np

# Set to true to use an actual device.
use_device = False
simulator = Aer.get_backend('statevector_simulator')
if use_device:
    provider = IBMProvider(instance='ibm-q/open/main')
    device = least_busy(provider.backends(
        filters=lambda x: x.configuration().n_qubits >= 6 and
        not x.configuration().simulator and x.status().operational==True))
    print("Using backend: ", device)
else:
    device = Aer.get_backend('qasm_simulator')

# Init quantum registers: 4 agents and 2 systems.
alice = QuantumRegister(1, 'alice')
bob = QuantumRegister(1, 'bob')
ursula = QuantumRegister(1, 'ursula')
wigner = QuantumRegister(1, 'wigner')
r = QuantumRegister(1, 'r')
s = QuantumRegister(1, 's')
c = ClassicalRegister(2)
qc = QuantumCircuit(r, alice, s, bob, ursula, wigner, c)

# Initial state of R is sqrt(1/3)|0> + sqrt(2/3)|1>.
desired_r = [ math.sqrt(1/3), math.sqrt(2/3) ]
qc.initialize(desired_r, [r])

# Alice measures R in computational basis. She records the result in her memory
# and if she obtained 1, she applies a Hadamard to S.
qc.cx(r, alice)
qc.ch(alice, s)
qc.barrier()

# Bob measures S in computational basis.
qc.cx(s, bob)
qc.barrier()

# Simulate the state vector.
statevector = execute(qc, simulator).result().get_statevector(qc, decimals=3)
labeled_statevector = list(zip([format(i, '04b')[::-1] for i in range(16)], statevector))
print("\nState vector (RASB) is:\n", np.vstack(labeled_statevector))

# Ursula measures Alice's lab (R + A) in the basis
# |ok> = sqrt(1/2)(|00> - |11>) and
# |fail> = sqrt(1/2)(|00> + |11>).
qc.cx(r, alice)
qc.h(r)
qc.cx(r, ursula)
# qc.h(r)          # for symmetry
# qc.cx(r, alice)  # for symmetry
qc.barrier()

# Wigner measures Bob's lab (S + B) in the basis
# |ok> = sqrt(1/2)(|00> - |11>) and
# |fail> = sqrt(1/2)(|00> + |11>).
qc.cx(s, bob)
qc.h(s)
qc.cx(s, wigner)
# qc.h(s)        # for symmetry
# qc.cx(s, bob)  # for symmetry
qc.barrier()

# Measure Ursula and Wigner's qubits.
qc.measure(ursula, c[0])
qc.measure(wigner, c[1])

print(qc)

# Simulate the state vector.
# statevector = execute(qc, simulator).result().get_statevector(qc, decimals=3)
# labeled_statevector = list(zip([format(i, '06b')[::-1] for i in range(64)], statevector))
# print("\nState vector (RASBUW) is:\n", np.vstack(labeled_statevector))

# Execute the circuit on the device.
job = execute(qc, device, shots=1024)
result = job.result()
counts = result.get_counts(qc)
print("\nTotal counts are:", dict(sorted(counts.items())))
