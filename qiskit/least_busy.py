#!/usr/bin/env python3

"""Find the least busy device backend on IBMQ.
"""

from qiskit_ibm_provider import IBMProvider, least_busy

# Register with the API.
provider = IBMProvider(instance='ibm-q/open/main')

# Get the least busy backend satisfying requirements.
backend = least_busy(provider.backends(
    filters=lambda x: x.configuration().n_qubits >= 5 and
    not x.configuration().simulator and x.status().operational==True))
print("least busy backend: ", backend)
