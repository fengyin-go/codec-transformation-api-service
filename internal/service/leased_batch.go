package service

import "codec/internal/store"

func ConvertBatchWithLeases(pool *store.LeasePool, inputs []string, convert func(string) (string, error)) ([]string, error) {
	outputs := make([]string, 0, len(inputs))
	for _, input := range inputs {
		lease, err := pool.Acquire()
		if err != nil {
			return nil, err
		}
		defer lease.Close()
		output, err := convert(input)
		if err != nil {
			return nil, err
		}
		defer lease.Commit()
		outputs = append(outputs, output)
	}
	return outputs, nil
}
