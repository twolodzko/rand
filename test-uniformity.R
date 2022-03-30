
freqs <- table(read.delim("result.data")[, 1])

print(round(freqs / sum(freqs) * 100, 3))

test <- chisq.test(freqs)
print(test)

if (test$p.value < 0.05) {
    stop("the distribution of samples is not uniform")
}

plot(freqs)
