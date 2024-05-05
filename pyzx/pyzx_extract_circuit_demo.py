import streamlit as st
import pyzx as zx

from qiskit.circuit import QuantumCircuit
from pyzx.extract import extract_circuit

with st.sidebar:
    with st.form("form1"):
        st.write("Generate a random circuit consisting of CNOT, HAD, and phase gates.")

        qubits1 = st.slider("qubits", 2, 16, 6)
        depth = st.slider("depth", 5, 100, 25)
        p_had = st.slider("prob(HAD)", 0.0, 0.5, 0.2)
        p_t = st.slider("prob(T)", 0.0, 0.5, 0.2)
        submitted1 = st.form_submit_button("Submit")

    with st.form("form2"):
        st.write("Generate a random CNOTs circuit.")
        qubits2 = st.slider("qubits", 2, 16, 6)
        depth2 = st.slider("depth", 5, 100, 25)
        submitted2 = st.form_submit_button("Submit")

st.title("PyZX reduce and extract demo")

if submitted1:
    st.header("Original circuit")

    orig_circ = zx.generate.CNOT_HAD_PHASE_circuit(qubits1, depth, p_had, p_t)
    orig_qasm = extract_circuit(orig_circ.to_graph()).to_basic_gates().to_qasm()
    orig_qiskit_circ = QuantumCircuit().from_qasm_str(orig_qasm)
    st.pyplot(orig_qiskit_circ.draw(output="mpl"))
    st.text(orig_circ.stats())

    st.header("Original circuit's graph")

    g = orig_circ.to_graph()
    st.pyplot(zx.draw(g))
    st.text(g)

    st.header("Reduced graph")

    zx.full_reduce(g)
    g.normalize()
    st.pyplot(zx.draw(g))
    st.text(g)

    st.header("Extracted circuit")

    new_circ = extract_circuit(g).to_basic_gates()
    new_qasm = new_circ.to_qasm()
    new_qiskit_circ = QuantumCircuit().from_qasm_str(new_qasm)
    st.pyplot(new_qiskit_circ.draw(output="mpl"))
    st.text(new_circ.stats())

if submitted2:
    st.header("CNOTs circuit")
    cnots_graph = zx.generate.cnots(qubits2, depth2)
    st.pyplot(zx.draw(cnots_graph))
    st.text(cnots_graph.stats())
