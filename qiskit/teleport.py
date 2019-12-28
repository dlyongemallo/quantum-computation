#!/usr/bin/env python3

"""Quantum teleportation.
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)
from qiskit.providers.ibmq import least_busy
from qiskit.providers.ibmq.job.exceptions import IBMQJobFailureError
# from qiskit.visualization import plot_histogram, plot_bloch_multivector
from numpy import pi
import random

# Register with the API.
IBMQ.load_account()
provider = IBMQ.get_provider(hub='ibm-q')

# Set to true to use an actual device.
use_device = False
simulator = Aer.get_backend('statevector_simulator')
if use_device:
    device = least_busy(provider.backends(
        filters=lambda x: x.configuration().n_qubits <= 5 and
        not x.configuration().simulator and x.status().operational==True))
else:
    device = Aer.get_backend('qasm_simulator')

# Init 3 qubits and 2 classical bits.
q = QuantumRegister(3)
c = [ ClassicalRegister(1) for _ in range(3) ]
circuit = QuantumCircuit(q)
for creg in c:
    circuit.add_register(creg)

# Create the qubit to be teleported.
paramX = random.random() * pi
paramY = random.random() * pi
circuit.rx(paramX, q[0])
circuit.ry(paramY, q[0])
circuit.barrier()

# Record the input state using the simulator.
job = execute(circuit, simulator)
result = job.result()
inputstate = result.get_statevector(circuit, decimals=3)
print("\nState vector is:", inputstate)
# plot_bloch_multivector(inputstate)

# Create a Bell state (psi+) between q1 and q2.
circuit.h(q[1])
circuit.cx(q[1], q[2])

# Bell measurement of qubits on Alice's side.
circuit.cx(q[0], q[1])
circuit.h(q[0])
circuit.measure(range(2), range(2))
circuit.barrier()

# Send two classical bits to Bob to fix his qubit.
'''
TODO: This is a work-around for the fact that it is not possible to perform
any instruction on a measured qubit (error 7006), and also not possible to
classically condition on a single qubit
(https://github.com/Qiskit/qiskit-terra/issues/1160).
Otherwise, we would just do this:
circuit.cx(q[1], q[2])
circuit.cz(q[0], q[2])
'''
circuit.x(q[2]).c_if(c[1],1)
circuit.z(q[2]).c_if(c[0],1)

print("Circuit:")
print(circuit.draw())

# Execute the circuit on the simulator.
job = execute(circuit, simulator)

# Get the resulting state vector.
result = job.result()
outputstate = result.get_statevector(circuit, decimals=3)
print("\nState vector is:", outputstate)
# plot_bloch_multivector(outputstate)

# Undo the rotations on the output state.
circuit.ry(-paramY, q[2])
circuit.rx(-paramX, q[2])
circuit.measure(q[2], c[2])

# Run the job on the device.
job = execute(circuit, device, shots=1024)

# Get the result counts.
try:
    result = job.result()
    counts = result.get_counts(circuit)
    print("\nTotal counts are:", counts)
    # plot_histogram(counts)
except IBMQJobFailureError:
    print(job.error_message())
