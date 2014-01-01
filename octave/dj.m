% Deutsch-Jozsa algorithm

% Number of bits in input to f(x_{1}...x_{n}).
n = 3;

% Number of bits needed for each definition of f(x).
N = 2**n;

% Declare the hidden promise function.
function u = Uf(p, w)
    % The vector of the promise function must be defined such that
    % p[x+1] = f(x) where x is interpreted as an integer index.
    % Note that this vector must be of dimensional (N,1).

    % Convert the {0,1}-vector to a corresponding {+1,-1} vector.
    % (The operation should really be "q = exp(i*pi*p)", but this is faster.)
    q = -2*p+1;

    % Apply the function and emulate "phase kickback".
    u = q.*w;
end

% Determine whether a vector of bits is a generalized Morse sequence.
function m = ismorse(bvec)
    len = size(bvec,1);

    % A vector of length 1 is always Morse.
    if (len == 1)
        m = true;
        return
    end

    % If the two halves are identical or complementary, then the vector is
    % Morse is half of it is Morse.
    lefthalf = bvec(1:len/2);
    righthalf = bvec(len/2+1:end);
    if (lefthalf == righthalf || lefthalf == 1 - righthalf)
        m = ismorse(lefthalf);
        return
    end

    % Otherwise, it isn't.
    m = false;
end

% Determine whether the bits are constant, morse, balanced, or none of the above.
% Format bits for display, and also indicate if constant or balanced (or neither).
function [bvec, type, formatted] = nicebits(number, N)
    bits = dec2bin(number,N);
    bvec = (bits-'0')';
    s = sum(bvec);
    if (s == N/2)
        if (ismorse(bvec))
            type = 'bm';
            formatted = sprintf('%s (bm)', bits);
        else
            type = 'b';
            formatted = sprintf('%s (b) ', bits);
        end
    elseif (s == 0 || s == N)
        type = 'c';
        formatted = sprintf('%s (c) ', bits);
    else
        type = '';
        formatted = sprintf('%s     ', bits);
    end
end

% Hadamard gate of size N*N.
H = hadamard(N);

% Denominator in probabilities of measuring each output.
denom = zeros(1,N);

% Keep count of each balanced and Morse strings.
nBalanced = 0;
nMorse = 0;

% Keep track of sizes of orthogonal sets.
orthoSetSize = zeros(1,2**(N-1));

% Cycle through all possible promise functions which begin with '0'.
% (The binary complements result in the same output but with -1 global phase.)
for i = 0:(2**(N-1))-1
    % Create the |+>^{n} vector.
    v = ones(N,1);

    % Get the promise function as a {0,1}-vector.
    [p, type, formatted] = nicebits(i,N);

    % Count whether morse or balanced.
    if (isempty(type))
        % continue
    elseif (type == 'bm')
        nMorse = nMorse + 1;
        nBalanced = nBalanced +1;
    elseif (type == 'b')
        nBalanced = nBalanced +1;
    end

    % Apply the hidden promise function.
    v = Uf(p, v);

    % Apply the Hadamard transform.
    v = H*v;

    % Finally, normalize the vector.
    v = v/N;

    % Print the output.
    % disp(sprintf('%s: %s', formatted,
    %     strrep(['(' sprintf(' % 2.3f,', v) ')'], ',)', ' )')))

    % Tally the probability of measuring |0...0>.
    denom = denom + abs(v').**2;

    % Find sets of orthogonal vectors.
    vs(i+1,:) = v;
    for j = 1:i
        if (dot(vs(j,:),v) == 0)
            % disp(sprintf('%s : %s', dec2bin(j-1,N), dec2bin(i,N)));
            orthoSetSize(j)++;

            orthoSetSize(i+1)++;
            % break;
        end
    end
end

% Output the sizes of the orthogonal sets.
for i = 0:(2**(N-1))-1
    if (orthoSetSize(i+1) != 0)
        disp(sprintf('%s : %d', dec2bin(i,N), orthoSetSize(i+1)));
    end
end

% Output the probability denominator.
% disp(sprintf('Denominator in p(|0...0>): %d', denom));
disp(sprintf('%s  sum: %s',
    repmat(' ', 1, N),
    strrep(['(' sprintf(' % 2.3f,', denom) ')'], ',)', ' )')));
disp(sprintf('#entries (beginning with 0): %d', 2**(N-1)));
disp(sprintf('#balanced: %d (%d balanced and Morse)',
    nBalanced, nMorse));

