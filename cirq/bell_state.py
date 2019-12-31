#!/usr/bin/env python3

"""Creates a simple Bell state.
"""

import cirq

def main():
    # Init 2 qubits.
    circuit = cirq.Circuit()
    q0, q1 = cirq.LineQubit.range(2)

    # Create a Bell state (psi+).
    circuit.append([
        cirq.H(q0),
        cirq.CX(q0, q1),
        cirq.measure(q0, q1, key='result')
    ])
    print("Circuit:")
    print(circuit)

    # Simulate the circuit.
    device = cirq.Simulator()
    result = device.run(circuit, repetitions=1000)
    print(result.histogram(key='result'))

if __name__ == '__main__':
    main()
