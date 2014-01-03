% Return a Hadamard matrix acting on n qubits.
function [h, d] = Hadamard(nQubits)
    h1 = [1 1; 1 -1];
    d1 = sqrt(2);
    if (nQubits <= 1)
        h = h1;
        d = d1;
        return;
    end
    [hr, dr] = Hadamard(nQubits - 1);
    h = kron(h1, hr);
    d = d1 * dr;
end

% Return a controlled version of a gate.
function [gout, dout] = Controlled(gin, din)
    dout = sqrt(2) * din;
    gout = kron([1 0; 0 0], eye(size(gin))) + kron([0 0; 0 1], gin);
end

% Pauli X, Y, Z, and Hadamard gates.
X = [0 1; 1 0];
Y = [0 -i; i 0];
Z = [1 0; 0 -1];
[h, d] = Hadamard(1);

% H*I, followed by c-X, followed by H*I
h_eye = kron(h, eye(2));
[cx, cxd] = Controlled(X, sqrt(2));
circuit1 = h_eye * cx * h_eye

% H*H, followed by c-Z, followed by H*H
h2 = Hadamard(2);
[cz, czd] = Controlled(Z, sqrt(2));
h2 * cz * h2 / czd
