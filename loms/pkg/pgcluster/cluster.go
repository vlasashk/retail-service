package pgcluster

import (
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cluster struct {
	write         *pgxpool.Pool
	read          []*pgxpool.Pool
	roundRobinIdx int
	mu            *sync.Mutex
}

func New() *Cluster {
	return &Cluster{
		roundRobinIdx: 0,
		mu:            &sync.Mutex{},
	}
}

func (c *Cluster) SetWriter(master *pgxpool.Pool) *Cluster {
	c.write = master
	return c
}

func (c *Cluster) AddReader(readPool ...*pgxpool.Pool) *Cluster {
	c.read = append(c.read, readPool...)
	return c
}

func (c *Cluster) GetWriter() (*pgxpool.Pool, error) {
	if c.write == nil {
		return nil, errors.New("writer is not set")
	}
	return c.write, nil
}

func (c *Cluster) GetReader() (*pgxpool.Pool, error) {
	if len(c.read) == 0 {
		return nil, errors.New("reader is not set")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	reader := c.read[c.roundRobinIdx]
	c.roundRobinIdx = (c.roundRobinIdx + 1) % len(c.read)

	return reader, nil
}

func (c *Cluster) Close() {
	c.write.Close()
	for _, pool := range c.read {
		pool.Close()
	}
}
