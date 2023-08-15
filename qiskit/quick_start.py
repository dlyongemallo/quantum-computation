# https://qiskit.org/ecosystem/ibm-provider/tutorials/Migration_Guide_from_qiskit-ibmq-provider.html
# Quick check that IBMQ credentials are working.

from qiskit import QuantumCircuit, transpile
from qiskit_ibm_provider import IBMProvider

# Save account credentials.
# IBMProvider.save_account(token=MY_API_TOKEN)

# Load previously saved account credentials.
provider = IBMProvider()

# Create a circuit
qc = QuantumCircuit(2)
qc.h(0)
qc.cx(0, 1)
qc.measure_all()

# Select a backend.
backend = provider.get_backend("ibmq_qasm_simulator")

# Transpile the circuit
transpiled = transpile(qc, backend=backend)

# Submit a job.
job = backend.run(transpiled)
# Get results.
print(job.result().get_counts())
