#!/usr/bin/env python3

"""Find the least busy device backend on IBMQ.
"""

from qiskit import IBMQ
from qiskit.providers.ibmq import least_busy

# Register with the API.
IBMQ.load_account()
provider = IBMQ.get_provider(hub='ibm-q')

# Get the least busy backend satisfying requirements.
# provider.backends()
backend = least_busy(provider.backends(
    filters=lambda x: x.configuration().n_qubits <= 5 and
    not x.configuration().simulator and x.status().operational==True))
print("least busy backend: ", backend)
