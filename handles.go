package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

//go:embed log.html
var logPage []byte

//go:embed message.html
var messageTepmlate string

type Message struct {
	Time    string
	Method  string
	Host    string
	Url     string
	Headers []string
	Content string
}

func NewMessage(r *http.Request) Message {

	m := Message{}

	m.Time = time.Now().Format("2006-01-02 15:04:05")
	m.Method = r.Method
	m.Host = r.Host
	m.Url = r.URL.String()

	for k, v := range r.Header {
		header := fmt.Sprintf("%s: %s", k, strings.Join(v, " "))
		m.Headers = append(m.Headers, header)
	}

	defer r.Body.Close()
	content, _ := io.ReadAll(r.Body)
	if len(content) == 0 {
		m.Content = "*empty*"
	} else {
		m.Content = string(content)
	}

	return m
}

func (m Message) String() string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s %s%s\n\n", m.Time, m.Method, m.Host, m.Url))

	sb.WriteString("Headers:\n")
	for _, s := range m.Headers {
		sb.WriteString(s)
		sb.WriteString("\n")
	}

	sb.WriteString("\nContent:\n")
	sb.WriteString(m.Content)

	return sb.String()
}

func (m Message) HTML() string {
	tmpl, err := template.New("message").Parse(messageTepmlate)
	if err != nil {
		return err.Error()
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, m)
	if err != nil {
		return err.Error()
	}

	return sb.String()
}

func HandleLog(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.String() == "/" {
		w.Write(logPage)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	message := NewMessage(r)
	messageString := message.String()

	fmt.Println(messageString)
	fmt.Fprint(w, messageString)
	hub.BroadcastMessage(message.HTML())
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	wsconn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	connection := hub.AddConnection(wsconn)

	go func() {
		for {
			m := <-connection.send
			err := wsutil.WriteServerMessage(wsconn, ws.OpText, []byte(m))
			if err != nil {
				hub.DeleteConnection(connection)
				return
			}
		}
	}()
}
