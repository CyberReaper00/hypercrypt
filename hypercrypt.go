package main

import (
	"os"
    "fmt"
	"log"
	"flag"
	"path"
	"strings"
	"github.com/CyberReaper00/helper_utils/humain"
)

func to_binary(data string) []byte { return []byte(data) }

func convertor(data []byte, key string, action string) []byte {
	var shift int
	key_len := len(key)
	e_data  := make([]byte, len(data))

	if key_len % 8 == 0 { shift = 7
	} else { shift = key_len % 8 }

	var rotated byte
	for i, b := range data {
		switch action {
			case "e": rotated = (b << shift) | (b >> (8 - shift))
			case "d": rotated = (b >> shift) | (b << (8 - shift))
		}
		e_data[i] = rotated & 0xFF
	}
	return e_data
}

// func encryption(data []byte, key_len int) []byte {
// 	var shift int
// 	e_data := make([]byte, len(data))
//
// 	if key_len % 8 == 0 { shift = 7
// 	} else { shift = key_len % 8 }
//
// 	for i, b := range data {
// 		rotated	 := (b << shift) | (b >> (8 - shift))
// 		e_data[i] = rotated & 0xFF
// 	}
// 	return e_data
// }
//
// func decryption(data []byte, key_len int) []byte {
// 	var shift int
// 	d_data := make([]byte, len(data))
//
// 	if key_len % 8 == 0 { shift = 7
// 	} else { shift = key_len % 8 }
//
// 	for i, b := range data {
// 		rotated  := (b >> shift) | (b << (8 - shift))
// 		d_data[i] = rotated & 0xFF
// 	}
// 	return d_data
// }

func save_data_to_file(entry string, data []byte, keep bool, action string) {

	var new_file string

	dir  := path.Dir(entry)
	file := path.Base(entry)
	ext	 := path.Ext(file)
	filename, _, _ := strings.Cut(file, ext)

	switch action {
		case "e":
			_ = os.Mkdir(fmt.Sprintf("%s_e", dir), 0755)
			new_file = fmt.Sprintf("%s_e/%s_e%s", dir, filename, ext)
			action = "encrypted"

		case "d":
			_ = os.Mkdir(fmt.Sprintf("%s_d", dir), 0755)
			new_file = fmt.Sprintf("%s_d/%s_d%s", dir, filename, ext)
			action = "decrypted"
	}

	file_ptr, err := os.Create(new_file)
	humain.Err("couldnt create file: %s\n%s", err, new_file, err)

	_, err = file_ptr.Write(data)
	humain.Err("couldnt write to %s\n%s", err, new_file, err)
	defer file_ptr.Close()

	fmt.Printf("%s%s's data has been %s into %s'\n\n", filename, ext, action, new_file)
	if !keep {
		err = os.Remove(entry)
		humain.Err("couldnt remove %v\n%v", err, entry, err)
	}
}

func cryptify_dir(main_dir string, dir []os.DirEntry, i int, keep bool, key string, action string) {
	if (i >= len(dir)) { return }

	filepath := path.Join(main_dir, dir[i].Name())
	filedata, err := os.ReadFile(filepath)
	humain.Err("couldnt read file: %v\n%v", err, dir[i].Name(), err)

	var output []byte
	switch action {
		case "e":
			output = convertor(filedata, key, action)
		case "d":
			output = convertor(filedata, key, action)
	}

	save_data_to_file(filepath, output, keep, action)
	cryptify_dir(main_dir, dir, i+1, keep, key, action)
}

func main() {

	t := flag.Bool("t", false, "description for t")
	f := flag.Bool("f", false, "description for f")
	d := flag.Bool("d", false, "description for d")
	k := flag.Bool("k", false, "description for k")
	E := flag.Bool("E", false, "description for E")
	D := flag.Bool("D", false, "description for D")

	flag.Usage = func() {
		fmt.Println("usage: [-h] [-k] [-options ...]")

		fmt.Println("\n\033[7m SUMMARY \033[0m")
		fmt.Println("  Hypercrypt is designed to be a simple cryptographic utility" +
		"\n  that takes in data from a terminal, file or every file in a directory and" +
		"\n  puts the data through a cryptographic process to be converted" +

		"\n\n  When the data is gathered and is processed, the original file will be deleted" +
		"\n  and the new processed file will be put in its place - the new file will have" +
		"\n  either '_e' or '_d' appended to the original name based on the users usage")

		fmt.Println("\n\033[7m MISC OPTIONS \033[0m")
		fmt.Println("  -h\n\tdisplay this message and exit")

		fmt.Println("  -k\n\twhen active it will keep the original file instead of" +
		"\n\tdeleting it")

		fmt.Println("\n\033[7m DATA OPTIONS \033[0m")
		fmt.Println("  -t\n\twill ask for the message to be processed in the terminal" +
		"\n\tNOTE: this option is only here for demonstration purposes, using it" +
		"\n\tas an actual means of en/decrypting data would be discouraged")

		fmt.Println("  -f\n\twill ask for the file to be processed")

		fmt.Println("  -d\n\twill ask for the directory that holds all the files to be processed")

		fmt.Println("\n\033[7m MODE OPTIONS \033[0m")
		fmt.Println("  -E\n\tactivate encryption for the data provided")

		fmt.Println("  -D\n\tactivate decryption for the data provided")
	}

	flag.Parse()

	if !*t && !*f && !*d { log.Fatalln("Data flag not defined, try -h") }
	if !*E && !*D { log.Fatalln("Cryptographic flag not defined, try -h") }

	if *t {
		input	:= humain.Input("Enter message").(string)
		key		:= humain.Input("Enter key").(string)

		switch {
			case *E:
				bin_data := to_binary(input)
				encr := convertor(bin_data, key, "e")
				fmt.Printf("%s\n", encr)

			case *D:
				bin_data := to_binary(input)
				decr := convertor(bin_data, key, "d")
				fmt.Printf("%s\n", decr)
		}

	} else if *f {
		input	:= humain.Input("Enter name or path of file").(string)
		key		:= humain.Input("Enter key").(string)

		contents, err := os.ReadFile(input)
		humain.Err("couldnt locate file: %v\n%v", err, contents, err)

		switch {
			case *E:
				encr := convertor(contents, key, "e")
				save_data_to_file(input, encr, *k, "e")

			case *D:
				decr := convertor(contents, key, "d")
				save_data_to_file(input, decr, *k, "d")
		}

	} else if *d {
		dir_name:= humain.Input("Enter name or path of directory").(string)
		key		:= humain.Input("Enter key").(string)

		dir_contents, err := os.ReadDir(dir_name)
		humain.Err("couldnt locate directory: %v", err, dir_name)

		switch {
			case *E: cryptify_dir(dir_name, dir_contents, 0, *k, key, "e")
			case *D: cryptify_dir(dir_name, dir_contents, 0, *k, key, "d")
		}
	}
}
