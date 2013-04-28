// Copyright (c) 2012, André Simon
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL ANDRE SIMON BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.net/websocket"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
)

type Message struct {
	Type     string
	Template string
	Snippet  string
}

func chooseTemplate(tpl string) (tmpl *template.Template, err error) {
	if len(tpl) > 0 {
		return template.New("test").ParseFiles(tpl)
	}
	return template.New("test").Parse("{{.Snippet}}")
}

func webSocketHandler(c *websocket.Conn) {
	defer c.Close()
	var s string
	var m Message
	var buffer bytes.Buffer
	var tmpl *template.Template
	var err error
	var out []byte

	r := bufio.NewReader(c)
	s, _ = r.ReadString('\n')

	json.Unmarshal([]byte(s), &m)

	var n int32
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	tmpFileName := fmt.Sprint(n)
	tmpFileName += "."
	tmpFileName += m.Type

	// go requires .go file suffix
	file, err := os.Create(fmt.Sprintf("%s%c%s", os.TempDir(), os.PathSeparator, tmpFileName))
	if err != nil {
		fmt.Fprint(c, err)
		return
	}
	tmpl, err = chooseTemplate(m.Template)
	if err != nil {
		fmt.Fprint(c, "template error!")
		return
	}
	if len(m.Template) > 0 {
		err = tmpl.ExecuteTemplate(&buffer, m.Template, m)
	} else {
		err = tmpl.Execute(&buffer, m)
	}
	if err != nil {
		fmt.Fprint(c, "template execution error!")
		return
	}
	file.WriteString(buffer.String())
	fName := file.Name()
	file.Close()

	if m.Type == "go" {
		out, err = exec.Command(m.Type, "run", fName).Output()
	} else {
		out, err = exec.Command(m.Type, fName).Output()
	}
	if err != nil {
		fmt.Print(err)
		fmt.Fprint(c, "execution/compilation error!")
		return
	}

	fmt.Fprint(c, string(out))
}

func main() {
	var myPort = flag.Int("port", 12345, "set listening port (adjust static/run.js accordingly)")
	var myBaseDir = flag.String("dir", ".", "set base directory")
	flag.Parse()
	fmt.Printf("listening on port %v...\n", *myPort)
	http.Handle("/", http.FileServer(http.Dir(*myBaseDir)))
	http.Handle("/ws", websocket.Handler(webSocketHandler))
	err := http.ListenAndServe(":"+fmt.Sprint(*myPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
