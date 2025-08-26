package main

import (
	"os"
    "fmt"
	"log"
	"flag"
	"strings"
	"github.com/CyberReaper00/helper_utils/humain"
)

type FSObject interface {
	[]os.DirEntry | string
}

func to_binary(data string) []byte { return []byte(data) }

func encryption(data []byte, key_len int) []byte {
	e_data := make([]byte, len(data))
	shift  := key_len % 8

	for i, b := range data {
		rotated	 := (b << shift) | (b >> (8 - shift))
		e_data[i] = rotated & 0xFF
	}
	return e_data
}

func decryption(data []byte, key_len int) []byte {
	d_data := make([]byte, len(data))
	shift  := key_len % 8

	for i, b := range data {
		rotated  := (b >> shift) | (b << (8 - shift))
		d_data[i] = rotated & 0xFF
	}
	return d_data
}

func save_data_to_file[T FSObject](entry T, data []byte, i int, main_dir string, mode string) {

	var filename string
	var new_file string
	var path	 string
	entry_type := entry.(type)

	if entry_type == []os.DirEntry {
		filename = fmt.Sprintf(entry[i].Name())
		new_file = fmt.Sprintf(entry[i].Name()+"_e")
		path	 = fmt.Sprintf(main_dir+dir[i].Name())

	} else if entry_type == string {
		filename = entry
		new_file = entry+"_e"
		path	 = entry
	}

	file, err := os.Create(new_file)
	if err != nil { humain.Err("couldnt create file: %s\n", new_file) }

	_, err = file.Write(data)
	if err != nil { log.Fatalf("couldnt write to %s\n", new_file) }
	defer file.Close()

	action := ""
	switch mode {
		case "e": action = "encrypted"
		case "d": action = "decrypted"
	}

	fmt.Println(filename, "has been", action ,"as", new_file)
	err = os.Remove(path)
	if err != nil { log.Fatalf("couldnt remove %v\n%v", path, err) }
}

func encr_dir(main_dir string, dir []os.DirEntry, i int, key_len int) {
	if (i >= len(dir)) { return count }

	filedata, err := os.ReadFile(fmt.Sprintf(main_dir+dir[i].Name()))
	if err != nil { log.Fatalf("couldnt read file: %v\n%v", dir[i].Name(), err) }
	
	bin_input := to_binary(filedata)
	encr := encryption(bin_input, key_len)
	save_data_to_file(dir, encr, i, main_dir, "e")

	return encr_dir(dir, i+1, main_dir)
}

func decr_dir(main_dir string, dir []os.DirEntry, i int) {
	if (i >= len(dir)) { return count }

	filedata, err := os.ReadFile(fmt.Sprintf(main_dir+dir[i].Name()))
	if err != nil { log.Fatalf("couldnt read file: %v\n%v", dir[i].Name(), err) }
	
	bin_input := to_binary(filedata)
	decr := decryption(bin_input, key_len)
	save_data_to_file(dir, decr, i, main_dir, "d")

	return encr_dir(dir, i+1, main_dir)
}

func main() {
	input	:= humain.Input("Enter text")
	key		:= humain.Input("Enter key")

	text 	:= input.(string)
	key_len	:= len(key.(string))

	t := flag.Bool("t", false, "description for t")
	f := flag.Bool("f", false, "description for f")
	d := flag.Bool("d", false, "description for d")
	E := flag.Bool("E", false, "description for E")
	D := flag.Bool("D", false, "description for D")

	if *t {
		input := humain.Input("Enter message")
		filedata, err := os.ReadFile(input)
		if err != nil { humain.Err("couldnt locate file: %v\n%v", filename, err) }
		encr := encryption(filedata, key_len)
		save_data_to_file()

	} else if *f {
		dir_name := humain.Input("Enter name or path of directory")
		dir_contents, err := os.ReadDir(dir_name)
		if err != nil { humain.Err("couldnt locate directory: %v", dir_name) }

		switch {
			case *E: 
			case *D: 
			default: humain.Err("Proper arguments were not defined, try -h")
		}

	} else if *d {
		dir_name := humain.Input("Enter name or path of directory")
		dir_contents, err := os.ReadDir(dir_name)
		if err != nil { humain.Err("couldnt locate directory: %v", dir_name) }

		switch {
			case *E: encr_dir(dir_name, dir_contents, 0, key_len)
			case *D: decr_dir(dir_name, dir_contents, 0, key_len)
			default: humain.Err("Proper arguments were not defined, try -h")
		}

	} else { humain.Err("Proper arguments were not defined, try -h") }
}
