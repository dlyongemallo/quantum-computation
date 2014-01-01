n = 4;
N = 2**n;

% Hamming weights.
kvec = 0:1:N/2-1;

% Counts of hidden strings with that Hamming weight.
cvec = arrayfun(@(k) nchoosek(N,k), kvec);

% Probability of a specific hidden string of the Hamming weight.
pvec = arrayfun(@(k) (1 - 2*k/N)**2, kvec);

% Probability of any hidden string of the Hamming weight.
total = 2**N / N;
mvec = cvec .* pvec / (total / 2);

% Cumulative sum of probabilities.
svec = cumsum(mvec);

% Comparison to classical case.
plot(kvec, svec, kvec, cumsum(cvec)/sum(cvec));
