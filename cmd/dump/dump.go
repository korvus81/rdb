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
	fmt.Fprintf(p.w, "%q\tset\t%d\n", key, len(value))
}

func (p *decoder) Hset(key, field, value []byte) {
	if len(value) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q\thset\t%d\t%q\n", key, len(value), field)
}

func (p *decoder) Sadd(key, member []byte) {
	if len(member) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q\tsadd\t%d\n", key, len(member))
}

func (p *decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *decoder) Rpush(key, value []byte) {
	if len(value) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q\trpush\t%d\t%d\n", key, len(value), p.i)
	p.i++
}

func (p *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *decoder) Zadd(key []byte, score float64, member []byte) {
	if len(member) < minSize {
		return
	}
	fmt.Fprintf(p.w, "%q\tzadd\t%d\t%d\n", key, len(member), p.i)
	p.i++
}

func maybeFatal(err error) {
	if err != nil {
		log.Fatal(err)
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
