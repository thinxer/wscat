package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/mgutz/ansi"

	"golang.org/x/net/websocket"
)

var (
	flagPretty = flag.Bool("pretty", false, "Pretty print")
)

func main() {
	flag.Parse()

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	conn, err := websocket.Dial(u.String(), "", "http://"+u.Host)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var buf []byte
			if err := websocket.Message.Receive(conn, &buf); err != nil {
				return
			}
			if *flagPretty {
				color := ansi.ColorCode("green")
				reset := ansi.ColorCode("reset")
				fmt.Fprintf(os.Stdout, "%s%s%s\n", color, strings.Repeat("<", 10), reset)
				var v interface{}
				if err := json.Unmarshal(buf, &v); err != nil {
					os.Stdout.Write(buf)
				} else {
					buf, err = json.MarshalIndent(v, "", "  ")
					if err != nil {
						panic(err)
					}
					os.Stdout.Write(buf)
					fmt.Fprintln(os.Stdout)
				}
			} else {
				os.Stdout.Write(buf)
			}
		}
	}()
	io.Copy(conn, os.Stdin)
	conn.Close()
}
