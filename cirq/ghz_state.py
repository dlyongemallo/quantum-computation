#!/usr/bin/env python3

"""Creates a GHZ state.
"""

import cirq

def main():
    # Init 3 qubits.
    circuit = cirq.Circuit()
    q0, q1, q2 = cirq.LineQubit.range(3)

    # Create a GHZ state.
    circuit.append([
        cirq.H(q0),
        cirq.CX(q0, q1),
        cirq.CX(q0, q2),
        cirq.measure(q0, q1, q2, key='result')
    ])
    print("Circuit:")
    print(circuit)

    # Simulate the circuit.
    device = cirq.Simulator()
    result = device.run(circuit, repetitions=1024)
    print(result.histogram(key='result'))

if __name__ == '__main__':
    main()
