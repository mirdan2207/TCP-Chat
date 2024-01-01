package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

const N = 3 // number of seconds per message

type client struct {
	conn        net.Conn
	nick        string
	room        *room
	commands    chan<- command
	lastMsgTime time.Time
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args:= strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command {
				id: CMD_NICK,
				client: c,
				args: args,
			}
		case "/join":
			c.commands <- command {
				id: CMD_JOIN,
				client: c,
				args: args,
			}
		case "/rooms":
			c.commands <- command {
				id: CMD_ROOMS,
				client: c,
				args: args,
			}
		case "/msg":
			if time.Since(c.lastMsgTime) >= N*time.Second {
				c.commands <- command {
					id: CMD_MSG,
					client: c,
					args: args,
				}
				c.lastMsgTime = time.Now()
			} else {
				c.err(errors.New("you can only send a message once every 3 seconds"))
			}
		case "/quit":
			c.commands <- command {
				id: CMD_QUIT,
				client: c,
				args: args,
			}
		case "/help":
			c.showHelp()
		default:
			c.err(fmt.Errorf("unknown command: %s\ntry the command /help to see all commands", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}

func (c *client) showHelp() {
	helpText := `
	/help - Show this help message
	/nick <nickname> - Set your nickname
	/join <room> - Join or create a room
	/rooms - List available rooms
	/msg <message> - Send a message (once every 3 seconds)
	/quit - Disconnect from the chat`

	c.msg("Available commands:\n" + helpText)
}