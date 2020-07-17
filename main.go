package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			s := <-signalChannel
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				os.Exit(0)
			}
		}
	}()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/index.html")),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))
	e.Renderer = t
	e.Match([]string{"GET", "POST"}, "/", Home)

	e.Static("/chat", "public/chat")
	// Create a simple file server

	// Configure websocket route
	e.GET("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	e.Logger.Fatal(e.Start(":1234"))
}

func Home(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusOK, "index.html", nil)
	}

	inputText := c.FormValue("inputText")
	fmt.Println(inputText)
	result := ltcWalker(inputText)
	fmt.Println(result)

	return c.Render(http.StatusOK, "index.html", struct{ Result string }{result})
}

func handleConnections(c echo.Context) error {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			c.Logger().Errorf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		msg.Message = ltcWalker(msg.Message)
		broadcast <- msg
	}

	return nil
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func ltcWalker(input string) string {
	stringParts := strings.Split(input, " ")

	ltcString := bytes.NewBufferString("")
	for _, val := range stringParts {
		ltcString.WriteString(ltc(val) + " ")
	}

	return strings.Trim(ltcString.String(), " ")
}

func ltc(input string) string {
	if len(input) <= 2 {
		return input
	}

	return fmt.Sprintf("%c%d%c", input[0], len(input)-2, input[len(input)-1])
}
