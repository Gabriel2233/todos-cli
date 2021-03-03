package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

/*
	App - Cli that creates a todo list inside the user')s current project
	Some commands - New list, Add item to list, Remove item from list,
		Update items(in progress, done, not done), Print the list items
*/

func getPath(name string) string {
	return fmt.Sprintf("./todo-cli/%v.txt", name)
}

const (
	not_done_todo_header = "Not Done\n"
	doing_todo_header    = "Doing\n"
	done_todo_header     = "Done\n"
)

func writeTodoStructureToFile(f io.Writer) {
	f.Write([]byte(not_done_todo_header))
	f.Write([]byte(" \n"))
	f.Write([]byte(doing_todo_header))
	f.Write([]byte(" \n"))
	f.Write([]byte(done_todo_header))
	f.Write([]byte(" \n"))
}

func file2Lines(name string) ([]string, error) {
	f, err := os.Open(getPath(name))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func findStatusLines(lines []string) []int {
	var allStatus []int

	for i, status := range lines {
		if status == "Not Done" || status == "Doing" || status == "Done" {
			allStatus = append(allStatus, i)
		}
	}

	return allStatus
}

func writeTodoInIndex(path, name string, index int, lines []string) error {
	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += name
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

func writeTodoByStatus(lines []string, status, name, path string) error {

	statusLines := findStatusLines(lines)
	realPath := getPath(path)

	switch status {
	case "not done":
		writeTodoInIndex(realPath, name, statusLines[0]+1, lines)
	case "doing":
		writeTodoInIndex(realPath, name, statusLines[1]+1, lines)
	case "done":
		writeTodoInIndex(realPath, name, statusLines[2]+1, lines)
	default:
		log.Fatal("you must provide a valid status")
	}

	fmt.Println("Gotcha ;)")
	return nil
}

func deleteFile(name string) {
	err := os.Remove(getPath(name))

	if err != nil {
		log.Fatal(err)
	}
}

func createStructuredFile(name string) error {

	path := getPath(name)

	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return err
	}

	writeTodoStructureToFile(f)

	return nil
}

func createDir() {
	err := os.Mkdir("todo-cli", 0700)

	if err != nil {
		log.Fatal(err)
	}
}

func dirExists() bool {
	if _, err := os.Stat("./todo-cli"); os.IsNotExist(err) {
		return false
	}

	return true
}

func main() {
	create_cmd := flag.NewFlagSet("new", flag.ExitOnError)
	add_name := create_cmd.String("n", "default", "the created list name")

	delete_cmd := flag.NewFlagSet("del", flag.ExitOnError)
	del_name := delete_cmd.String("n", "default", "the deleted list name")

	insert_cmd := flag.NewFlagSet("todo", flag.ExitOnError)
	todo_list_name := insert_cmd.String("l", "default", "specify the list in which the todo will be stored")
	todo_name := insert_cmd.String("n", "default todo", "new todo name")
	todo_status := insert_cmd.String("s", "not done", "specifies the todo's current status (doing, done, not done)")

	switch os.Args[1] {
	case "new":
		create_cmd.Parse(os.Args[2:])
		exists := dirExists()

		if exists {
			createStructuredFile(*add_name)
		} else {
			createDir()
			createStructuredFile(*add_name)
		}

		fmt.Printf("created list %v", *add_name)

	case "del":
		delete_cmd.Parse(os.Args[2:])
		deleteFile(*del_name)

		fmt.Printf("%v was successfully deleted :)", *del_name)
	case "todo":
		insert_cmd.Parse(os.Args[2:])
		lines, err := file2Lines(*todo_list_name)

		if err != nil {
			log.Fatal(err)
		}

		escaped_todo_name := *todo_name + "\n"

		writeTodoByStatus(lines, *todo_status, escaped_todo_name, *todo_list_name)
	}
}
