#!/usr/bin/env python3

"""Create a Z-spider with 1 input and 2 outputs (i.e., a Z-copy spider).
"""

import sys, os, math
from fractions import Fraction

import pyzx as zx
from pyzx import print_matrix
from pyzx.basicrules import *

Z = zx.VertexType.Z
X = zx.VertexType.X
B = zx.VertexType.BOUNDARY
SE = zx.EdgeType.SIMPLE
HE = zx.EdgeType.HADAMARD

zcopy = zx.Graph()

in1  = zcopy.add_vertex(B, qubit=1, row=0)
z1   = zcopy.add_vertex(Z, qubit=1, row=1)
out1 = zcopy.add_vertex(B, qubit=0, row=2)
out2 = zcopy.add_vertex(B, qubit=2, row=2)

zcopy.add_edge((in1,z1))
zcopy.add_edge((z1,out1))
zcopy.add_edge((z1,out2))

# Optional: specify inputs and outputs, only matters for sequential composition.
# Alternatively: zcopy.auto_detect_io()
zcopy.set_inputs((in1,))
zcopy.set_outputs((out1,out2))

# This won't work in the terminal, obviously.
# zx.draw(zcopy)
print(zcopy)
