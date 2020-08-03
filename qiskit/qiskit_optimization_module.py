# Copied from: https://gist.github.com/ismaelfaro/2f9351259bb8107db0d06d774109b2b8
# example to solve a quadratic problem using a Quantum optimization algorithm

from qiskit.optimization import QuadraticProgram
from qiskit.optimization.algorithms import GroverOptimizer

from qiskit import Aer

from docplex.mp.model import Model

# construct problem with DOcplex
model = Model('docplex_model')
x, y = model.binary_var_list(2)
model.maximize(x*y + x + y)

# convert DOcplex model to Qiskit Quadratic Program
qp = QuadraticProgram()
qp.from_docplex(model)

# define a Quantum backend on which to run the optimization
# a Qiskit simulator in this case
backend = Aer.get_backend('statevector_simulator')

# use the Grover Adaptive Search (GAS) to solve the optimization problem
grover = GroverOptimizer(num_value_qubits=3, quantum_instance=backend)
result = grover.solve(qp)

# printing results: x=[1.0,1.0], fval=3
print(result)
