#!/usr/bin/env python3

"""Implementation of the Frauchiger-Renner thought experiment.
"""

import cirq
import numpy as np
import math

def main():
    # Init quantum registers: 4 agents and 2 systems.
    qc = cirq.Circuit()
    r, alice, s, bob, ursula, wigner = cirq.LineQubit.range(6)

    # Initial state of R is sqrt(1/3)|0> + sqrt(2/3)|1>.
    qc.append([
        cirq.SingleQubitMatrixGate(matrix=np.array(
            [[math.sqrt(1/3), -math.sqrt(2/3)],
            [math.sqrt(2/3), math.sqrt(1/3)]])).on(r),
    ])

    # Alice measures R in computational basis. She records the result in her
    # memory and if she obtained 1, she applies a Hadamard to S.
    qc.append([
        cirq.CX(r, alice),
        cirq.H.controlled().on(alice, s),
    ])

    # Bob measures S in computational basis.
    qc.append([
        cirq.CX(s, bob),
    ])

    # Ursula measures Alice's lab (R + A) in the basis
    # |ok> = sqrt(1/2)(|00> - |11>) and
    # |fail> = sqrt(1/2)(|00> + |11>).
    qc.append([
        cirq.CX(r, alice),
        cirq.H(r),
        cirq.CX(r, ursula),
        # cirq.H(r),
        # cirq.CX(r, alice),
    ])

    # Wigner measures Bob's lab (S + B) in the basis
    # |ok> = sqrt(1/2)(|00> - |11>) and
    # |fail> = sqrt(1/2)(|00> + |11>).
    qc.append([
        cirq.CX(s, bob),
        cirq.H(s),
        cirq.CX(s, wigner),
        # cirq.H(s),
        # cirq.CX(s, bob),
    ])

    # Measure Ursula and Wigner's qubits.
    qc.append([
        cirq.measure(ursula, wigner, key='result')
    ])

    print("Circuit:")
    print(qc)

    # Simulate the circuit.
    device = cirq.Simulator()
    result = device.run(qc, repetitions=1024)
    print(result.histogram(key='result'))

if __name__ == '__main__':
    main()
