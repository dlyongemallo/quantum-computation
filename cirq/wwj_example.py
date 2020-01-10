#!/usr/bin/env python3

"""Creates circuit posted to following GitHub issue:
https://github.com/Strilanc/Quirk/issues/450
https://algassert.com/quirk#circuit={%22cols%22:[[%22Y%22],[%22inputA1%22,%22Y^(A/2^n)%22],[%22Y^(A/2^n)%22,%22inputA1%22]]}
"""

import cirq

a, b = cirq.LineQubit.range(2)
c = cirq.Circuit(
    cirq.Y(a),
    cirq.Y(b).controlled_by(a)**0.5,
    cirq.Y(a).controlled_by(b)**0.5)
c2 = cirq.google.optimized_for_xmon(c)
print(c2)
