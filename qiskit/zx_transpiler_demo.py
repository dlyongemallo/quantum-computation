from qiskit.circuit.random import random_circuit
from qiskit.transpiler import PassManager
from zxpass import ZXPass
from qiskit import transpile

qc = random_circuit(4, 4)
print("Before:")
print(qc)

transpiled_qc = transpile(qc, optimization_level=3)
print("Default method:")
print(transpiled_qc)

pass_manager = PassManager(ZXPass())
zx_qc = pass_manager.run(qc)
# zx_qc = transpile(qc, optimization_method="zxpass", optimization_level=3)
print("zxpass:")
print(zx_qc)
