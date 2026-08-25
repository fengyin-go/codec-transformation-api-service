package service

import "codec/internal/store"

func ConvertBatchWithLeases(pool *store.LeasePool, inputs []string, convert func(string) (string, error)) ([]string, error) {
	outputs := make([]string, 0, len(inputs))
	for _, input := range inputs {
		// Wrap each iteration in its own function so the deferred lease
		// release runs at the end of the iteration, not at function return.
		// Otherwise acquired leases pile up and exhaust the pool before the
		// batch finishes, losing already-converted outputs.
		output, err := func(input string) (string, error) {
			lease, err := pool.Acquire()
			if err != nil {
				return "", err
			}
			defer lease.Close()
			output, err := convert(input)
			if err != nil {
				return "", err
			}
			lease.Commit()
			return output, nil
		}(input)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}
