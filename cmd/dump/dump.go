// This is a very basic example of a program that implements rdb.decoder and
// outputs a human readable diffable dump of the rdb file.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cupcake/rdb"
	"github.com/cupcake/rdb/nopdecoder"

	"github.com/google/gops/agent"
)

type decoder struct {
	i int
	nopdecoder.NopDecoder
	w io.Writer
}

// var minSize = 1024 * 1024 // 1MiB
var minSize = 0

func (p *decoder) StartDatabase(n int) {
}

func (p *decoder) Set(key, value []byte, expiry int64) {
	if len(value) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q set %d\n", key, len(value))
}

func (p *decoder) Hset(key, field, value []byte) {
	if len(value) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q %q hset %d\n", key, field, len(value))
}

func (p *decoder) Sadd(key, member []byte) {
	if len(member) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q sadd %d\n", key, len(member))
}

func (p *decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *decoder) Rpush(key, value []byte) {
	if len(value) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q[%d] rpush %d\n", key, p.i, len(value))
	p.i++
}

func (p *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *decoder) Zadd(key []byte, score float64, member []byte) {
	if len(member) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q[%d] zadd %d\n", key, p.i, len(member))
	p.i++
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(os.Args[1])
	maybeFatal(err)
	err = rdb.Decode(f, &decoder{
		w: bufio.NewWriter(os.Stdout),
	})
	maybeFatal(err)
}
