
freqs <- table(read.delim("result.data")[, 1])

test <- chisq.test(freqs)
print(test)

if (test$p.value < 0.05) {
    stop("the distribution of samples is not uniform")
}

plot(freqs)
