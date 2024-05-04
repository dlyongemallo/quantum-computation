from shiny import reactive, experimental
from shiny.express import input, ui, render
from shinyswatch import theme
import pyzx as zx
from pyzx import Circuit

from qiskit.circuit import QuantumCircuit
from pyzx.extract import extract_circuit

theme.darkly()
# zx.settings.mode = "notebook"
# zx.settings.drawing_backend = "matplotlib"

ui.page_opts(title="PyZX reduce and extract demo")

with ui.sidebar(title = "Generate a random circuit consisting of CNOT, HAD, and phase gates."):
    ui.input_slider("qubits", "qubits", 2, 16, 6)
    ui.input_slider("depth", "depth", 5, 100, 25)
    ui.input_slider("p_had", "prob(HAD)", 0.0, 0.5, 0.2)
    ui.input_slider("p_t", "prob(T)", 0.0, 0.5, 0.2)
    ui.input_action_button("submit", "Submit")

orig_circ = reactive.value()
orig_graph = reactive.value()


@reactive.effect
@reactive.event(input.submit)
def _():
    orig_circ.set(zx.generate.CNOT_HAD_PHASE_circuit(input.qubits(), input.depth(), input.p_had(), input.p_t()))


@reactive.calc
def orig_graph():
    return orig_circ().to_graph()


def reduced_graph():
    g = orig_graph()
    zx.full_reduce(g)
    g.normalize()
    return g


def reduced_circ():
    return extract_circuit(reduced_graph()).to_basic_gates()


def reduced_qiskit_circ():
    new_qasm = reduced_circ().to_qasm()
    new_qiskit_circ = QuantumCircuit().from_qasm_str(new_qasm)
    return new_qiskit_circ


with ui.card():
    experimental.ui.card_title("Original circuit")

    @render.plot
    @reactive.calc
    def plot_orig_circ():
        orig_qasm = extract_circuit(orig_circ().to_graph()).to_basic_gates().to_qasm()
        orig_qiskit_circ = QuantumCircuit().from_qasm_str(orig_qasm)
        return orig_qiskit_circ.draw(output="mpl")

    @render.text
    @reactive.calc
    def stats():
        return orig_circ().stats()

with ui.card():
    experimental.ui.card_title("Original circuit's graph")

    @render.plot
    @reactive.calc
    def plot_orig_graph():
        return zx.draw(orig_graph())

    @render.text
    @reactive.calc
    def orig_graph_stats():
        return orig_graph()


with ui.card():
    experimental.ui.card_title("Reduced graph")

    @render.plot
    @reactive.calc
    def plot_reduced_graph():
        return zx.draw(reduced_graph())

    @render.text
    @reactive.calc
    def reduced_graph_stats():
        return reduced_graph()


with ui.card():
    experimental.ui.card_title("Extracted circuit")

    @render.plot
    @reactive.calc
    def plot_extracted_circ():
        return reduced_qiskit_circ().draw(output="mpl")

    @render.text
    @reactive.calc
    def reduced_circ_stats():
        return reduced_circ().stats()
