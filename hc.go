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

func v_enc(plaintext []byte, key string) []byte {
	ciphertext := make([]byte, len(plaintext))
	key			= strings.ToUpper(key)
	keyIndex   := 0

	for i, b := range plaintext {
		if b >= 'A' && b <= 'Z' {
			shift 			:= int(key[keyIndex%len(key)]) - 'A'
			shiftedChar 	:= (b-'A'+byte(shift))%26 + 'A'
			ciphertext[i]	 = shiftedChar
			keyIndex++

		} else if b >= 'a' && b <= 'z' {
			shift 			:= int(key[keyIndex%len(key)]) - 'A'
			shiftedChar 	:= (b-'a'+byte(shift))%26 + 'a'
			ciphertext[i] 	 = shiftedChar
			keyIndex++

		} else { ciphertext[i] = b }
	}
	return ciphertext
}

func v_dec(ciphertext []byte, key string) []byte {
	key		   = strings.ToUpper(key)
	keyIndex  := 0
	plaintext := make([]byte, len(ciphertext))

	for i, b := range ciphertext {
		if b >= 'A' && b <= 'Z' {
			shift 			:= int(key[keyIndex%len(key)]) - 'A'
			shiftedChar 	:= (b-'A'-byte(shift)+26)%26 + 'A'
			plaintext[i] 	 = shiftedChar
			keyIndex++

		} else if b >= 'a' && b <= 'z' {
			shift 			:= int(key[keyIndex%len(key)]) - 'A'
			shiftedChar 	:= (b-'a'-byte(shift)+26)%26 + 'a'
			plaintext[i] 	 = shiftedChar
			keyIndex++

		} else { plaintext[i] = b }
	}
	return plaintext
}

func convertor(data []byte, key string, action string) []byte {
	var shift int
	key_len := len(key)
	e_data  := make([]byte, len(data))

	if key_len % 8 == 0 { shift = 7
	} else { shift = key_len % 8 }

	var rotated byte
	switch action {
		case "e":
			for i, b := range data {
				rotated 	= (b << shift) | (b >> (8 - shift))
				e_data[i] 	= rotated & 0xFF
			}

		case "d":
			for i, b := range data {
				rotated 	= (b >> shift) | (b << (8 - shift))
				e_data[i] 	= rotated & 0xFF
			}
	}
	return e_data
}

func save_data_to_file(entry string, data []byte, action string) {

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
}

func cryptify_dir(main_dir string, dir []os.DirEntry, i int, keep bool, key string, action string) {

	if (i >= len(dir) && !keep) {
		err := os.RemoveAll(main_dir)
		humain.Err("couldnt remove %v\n%v", err, dir[i-1].Name(), err)
		return

	} else if (i >= len(dir)) { return }

	filepath 		:= path.Join(main_dir, dir[i].Name())
	filedata, err 	:= os.ReadFile(filepath)
	humain.Err("couldnt read file: %v\n%v", err, dir[i].Name(), err)

	var new_data []byte
	var output   []byte
	switch action {
		case "e":
			new_data = v_enc(filedata, key)
			output   = convertor(new_data, key, action)
		case "d":
			new_data = convertor(filedata, key, action)
			output   = v_dec(new_data, key)
	}

	save_data_to_file(filepath, output, action)
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

		fmt.Println("" +
		"\nHypercrypt is designed to be a simple cryptographic utility that" +
		"\ntakes in data from a terminal, file or every file in a directory" +
		"\nand puts the data through a cryptographic process to be converted" +

		"\n\nWhen the data is gathered and is processed, the original file" +
		"\nwill be deleted and the new processed file will be put in its place" +
		"\nThe new file will have either '_e' or '_d' appended to the original" +
		"\nname, based on the users usage" +

		"\n\n\033[7m MISC OPTIONS \033[0m" +
		"\n  -h\n\tdisplay this message and exit" +
		"\n  -k\n\twhen active it will keep the original file instead of" +
		"\n\tdeleting it" +

		"\n\n\033[7m DATA OPTIONS \033[0m" +
		"\n  -t\n\twill ask for the message to be processed in the terminal" +
		"\n\tNOTE: this option is only here for demonstration purposes" +
		"\n\tusing it as an actual means of en/decrypting data would be" +
		"\n\tdiscouraged" +
		"\n  -f\n\twill ask for the file to be processed" +
		"\n  -d\n\twill ask for the directory that holds all the files to be" +
		"\n\tprocessed" +

		"\n\n\033[7m MODE OPTIONS \033[0m" +
		"\n  -E\n\tactivate encryption for the data provided" +
		"\n  -D\n\tactivate decryption for the data provided")
	}
	flag.Parse()

	if !*t && !*f && !*d { log.Fatalln("Data flag not defined, try -h") }
	if !*E && !*D { log.Fatalln("Cryptographic flag not defined, try -h") }

	if *t {
		input:= humain.Input("Enter message").(string)
		key	 := humain.Input("Enter key").(string)

		switch {
			case *E:
				bin_data 	:= to_binary(input)
				v_data		:= v_enc(bin_data, key)
				encr 		:= convertor(v_data, key, "e")
				fmt.Printf("%s\n", encr)

			case *D:
				bin_data 	:= to_binary(input)
				decr 		:= convertor(bin_data, key, "d")
				v_data 		:= v_dec(decr, key)
				fmt.Printf("%s\n", v_data)
		}

	} else if *f {
		input := humain.Input("Enter name or path of file").(string)
		key	  := humain.Input("Enter key").(string)

		contents, err := os.ReadFile(input)
		humain.Err("couldnt locate file: %v\n%v", err, contents, err)

		switch {
			case *E:
				v_data := v_enc(contents, key)
				encr   := convertor(v_data, key, "e")
				save_data_to_file(input, encr, "e")

			case *D:
				decr   := convertor(contents, key, "d")
				v_data := v_dec(decr, key)
				save_data_to_file(input, v_data, "d")
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
