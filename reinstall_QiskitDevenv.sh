#!/bin/bash

# Create virtual environment.
cd ~/workspace/Qiskit
/usr/local/bin/python3 -m venv QiskitDevenv
source ~/workspace/Qiskit/QiskitDevenv/bin/activate
pip install --upgrade pip

cd ~/workspace/Qiskit/qiskit
pip install -r requirements-dev.txt
pip install -e .

# Terra
pip uninstall -y qiskit-terra
cd ~/workspace/Qiskit/qiskit-terra
pip install cython
pip install -r requirements-dev.txt
pip install -e .
# python examples/python/using_qiskit_terra_level_0.py

# Aer
pip uninstall -y qiskit-aer
cd ~/workspace/Qiskit/qiskit-aer
pip install cmake scikit-build cython
pip install -e .

# Ignis
pip uninstall -y qiskit-ignis
cd ~/workspace/Qiskit/qiskit-ignis
pip install -r requirements-dev.txt
pip install -e .

# Aqua
pip uninstall -y qiskit-aqua
cd ~/workspace/Qiskit/qiskit-aqua
pip install -r requirements-dev.txt
pip install -e .

# IBMQ Provider
pip uninstall -y qiskit-ibmq-provider
cd ~/workspace/Qiskit/qiskit-ibmq-provider
pip install -r requirements-dev.txt
pip install -e .

# Check installation versions.
python -c "import qiskit; print(qiskit.__qiskit_version__)"
