package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8801", "http service address")
var outdir = flag.String("outdir", ".", "file out dir")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

var fileServer *FileServer

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	fileSession := fileServer.CreateFileSession()
	defer fileServer.DestoryFileSession(fileSession)

	remoteAddr := strings.ReplaceAll(r.RemoteAddr, ":", "_")
	now := time.Now().UnixNano() / 1000000

	baseDir := remoteAddr + "_" + strconv.FormatInt(now, 10)
	fileSession.SetDir(*outdir, baseDir)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		pck := NewLogPck()
		if pck.ParseMsg(message) == false {
			log.Println("error pck")
			break
		}

		fileSession.Write(pck)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	fileServer = NewFileServer()
	http.HandleFunc("/", serve)
	http.HandleFunc("/savedata", serve)
	log.Println("Start Listen Addr : ", *addr)
	fileServer.Start()
	log.Fatal(http.ListenAndServe(*addr, nil))

	fileServer.Stop()
}
