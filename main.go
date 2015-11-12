package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/mgutz/ansi"
)

var (
	flagPretty = flag.Bool("pretty", false, "Pretty print")
	flagDebug  = flag.Bool("debug", false, "Print message type")
)

func main() {
	flag.Parse()

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	h := http.Header{"Origin": {"http://" + u.Host}}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), h)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			t, buf, err := conn.ReadMessage()
			if *flagDebug {
				log.Println("type:", t, len(buf), err)
				continue
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

	b := bufio.NewReader(os.Stdin)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			break
		}
	}

	conn.Close()
}
