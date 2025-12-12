# Go benchmark plots

## Install

```bash
python3 -m venv env
source env/bin/activate
pip3 install -r requirements.txt
```

# Plot
Finally, just call the python script.
```bash
python3 plot.py benchmark_gc_off.txt benchmark_gc_off.pdf
```

This will generate the graph as pdf (`result_<timestamp>.pdf`).

## Example

```bash
#GOGC=100 go test -bench='BenchmarkSign' -benchmem -count=30 -cpu=1,2,4,6,8,10,12,14,16,18,20,22,24,26,28,30,32,34,36,38,40,42,48,64 -run=^$ > benchmark_gc_100.txt
GOGC=off go test -bench='BenchmarkSign' -benchmem -count=30 -cpu=1,2,4,6,8,10,12,14,16,18,20,22,24,26,28,30,32,34,36,38,40,42,48,64 -run=^$ > benchmark_gc_off.txt
#GOGC=10000 go test -bench='BenchmarkSign' -benchmem -count=30 -cpu=1,2,4,6,8,10,12,14,16,18,20,22,24,26,28,30,32,34,36,38,40,42,48,64 -run=^$ > benchmark_gc_10000.txt

python3 plot.py benchmark_gc_off.txt benchmark_gc_off.pdf
```

Happy benchmarking!
