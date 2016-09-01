// Licensed under GPL, 2016
// Refer to LICENSE for more details
// Refer to README for structural information

package main

import (
	"fmt"
	"log"
	"io"
	"net"
	"time"
	"math/rand"
)

func getRndWord() string {
	return words[rand.Intn(len(words))]
}

func getRndTask() (task Task, _ error) {
	for task.isAcceptable() != nil {
		task = Task{}

		for i := 1 + rand.Intn(4); i > 0; i-- {
			task.Match = append(task.Match, getRndWord())
		}

		for i := 1 + rand.Intn(4); i > 0; i-- {
			task.Dmatch = append(task.Dmatch, getRndWord())
		}
	}
	return
}

func gameServer() {
	listener, err := net.Listen("tcp", ":25921")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		go func(c io.ReadWriteCloser) {
			fmt.Fprintln(c, "@ ReGeX test server with words")
			defer c.Close()

			var (
				lvl  uint
				task Task
				sol  string
				err  error
			)

			for {
				// read level
				fmt.Fprint(c, ": ")
				_, err = fmt.Fscanf(c, "%d", &lvl)
				if err == io.EOF {
					fmt.Fprintln(c, "\n@ session quit")
					break
				} else if err != nil {
					fmt.Fprintf(c, "! %s\n", err.Error())
					break
				}

				if task, err = getRndTask(); err != nil{
					fmt.Fprintf(c, "! %s\n", err.Error())
					break
				}

				// print information
				task.calcLevel()
				fmt.Fprintf(c, "@ lvl: %d\n", task.Level)

				// print information
				for _, w := range task.Match {
					fmt.Fprintf(c, "+ %s\n", w)
				}
				for _, w := range task.Dmatch {
					fmt.Fprintf(c, "- %s\n", w)
				}

				// read suggestion
				fmt.Fprint(c, "> ")
				_, err = fmt.Fscanf(c, "%s", &sol)
				if err != nil {
					fmt.Fprintf(c, "! %s\n", err.Error())
					break
				}

				// respond adequately (wait if wrong)
				if task.matches(sol) {
					fmt.Fprintln(c, "@ correct")
					task.submit(sol)
				} else {
					fmt.Fprintln(c, "@ invalid")
					time.Sleep(time.Second * time.Duration(len(sol)))
				}

				// clear task object
				task = Task{}
			}
		}(conn)
	}
}
