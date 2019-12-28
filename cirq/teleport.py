#!/usr/bin/env python3

"""Quantum teleportation.
"""

import cirq
import random

def main():
    # Init 3 qubits.
    circuit = cirq.Circuit()
    q0, q1, q2 = cirq.LineQubit.range(3)

    # Create a Bell state (psi+) between q1 and q2.
    circuit.append([
        cirq.H(q1),
        cirq.CX(q1, q2),
    ])

    # Create the qubit to be teleported.
    paramX = random.random()
    paramY = random.random()
    circuit.append([
        cirq.X(q0)**paramX,
        cirq.Y(q0)**paramY
    ])

    # Bell measurement of qubits on Alice's side.
    circuit.append([
        cirq.CNOT(q0, q1), cirq.H(q0)
    ])
    circuit.append([cirq.measure(q0, q1)])

    # Send two classical bits to Bob to fix his qubit.
    circuit.append([cirq.CNOT(q1, q2), cirq.CZ(q0, q2)])

    print("Circuit:")
    print(circuit)

if __name__ == '__main__':
    main()
