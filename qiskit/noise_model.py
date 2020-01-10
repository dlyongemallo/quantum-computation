#!/usr/bin/env python3

"""Create a noise model for
"""

from qiskit import(
    QuantumCircuit,
    QuantumRegister,
    ClassicalRegister,
    execute, IBMQ, Aer)
from qiskit.providers.aer import noise

# Register with the API.
IBMQ.load_account()
provider = IBMQ.get_provider(hub='ibm-q')

# Get device properties.
device = provider.get_backend('ibmq_16_melbourne')
properties = device.properties()
coupling_map = device.configuration().coupling_map
noise_model = noise.device.basic_device_noise_model(properties)

# Init 2 qubits and 2 classical bits.
q = QuantumRegister(2)
c = ClassicalRegister(2)
circuit = QuantumCircuit(q, c)

# Create a Bell state (psi+).
circuit.h(q[0])
circuit.cx(q[0], q[1])
circuit.measure(q, c)
print(circuit.draw())

# Execute the circuit on the emulated device with noise.
emulator = Aer.get_backend('qasm_simulator')
job = execute(circuit, emulator, shots=1024, noise_model=noise_model,
    coupling_map=coupling_map,basis_gates=noise_model.basis_gates)

# Get the result counts.
result = job.result()
counts = result.get_counts(circuit)
print("\nTotal counts are:", dict(sorted(counts.items())))
