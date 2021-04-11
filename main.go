package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fsnotify/fsnotify"
)

func executeCommand(file *os.File, date string) []byte {
	cmd := "grep -a '" + date + "' /var/log/auth.log"
	//fmt.Println("BURASI CMD ->>>>", cmd)
	// source $HOME/.profile
	//out, err := exec.Command("bash", "-c", "grep -a 'Apr 11 02:11:06' /var/log/auth.log").Output()
	out, err := exec.Command("bash", "-c", cmd).Output()
	log.Println("OUT ->", string(out))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("\n\n\n GREP: %s ", out)
	//file.WriteString(string(out))
	return out
}

func main() {
	//current := time.Now().Format("Jan 02 3:4:5")
	file, err := os.Create("/home/baris/log.txt")
	//file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {

					//fmt.Println("YARATILAN TIME -> ", currentDateAndTime)
					exactTime := time.Now().Format("Jan 02 15:04:05")
					//fmt.Println("\n HEMEN YAZDIRILAN->", time.Now().Format("Jan 02 15:04:05"))
					log.Println("modified file: ", event.Name)
					//log.Println("event type", event.Op)
					//mt.Println("EXACT TIME->>", exactTime)
					fmt.Println("EXECUTE COMMAND USTU")
					output := executeCommand(file, exactTime)
					file.WriteString(string(output))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/var/log/auth.log")
	if err != nil {
		log.Println("hata basilan yer burasi")
		log.Fatal(err)
	}
	<-done
}
