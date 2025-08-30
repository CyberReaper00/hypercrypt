# Hypercrypt
I have made this as a simple hobby project, it is not extensive or all encompassing in any way but i just wanted to see if i could build it
## Usage
`hc -h`
- Hypercrypt requires two flags at minimum to function
	- Misc Flags
		- `-h`: shows the help message and exits
		- `-k`: keeps the original file after creating the new processed one
	- Data Flags
		- These flags tell `hc` where the data is coming from
		- `-t`: takes in data directly from the terminal
		- `-f`: takes in data from a specified file
		- `-d`: recursively takes in data from all files in any directory
	- Cryptographic Flags
		- These flags tell `hc` what to do with the data
		- `-E`: encrypts the data
		- `-D`: decrypts the data
## Process
Hypercrypt takes in data and a key from the user through one of the pre-defined methods and then the data goes through the following process

**Encryption**
- It converts all the data into byte data for easier parsing
- The byte data is then scrambled using the `vigenere cypher` which uses the key given by the user to discern how much the data should be scrambled
- Once the data is scrambled, it is then converted to binary
- The binary data is then flipped, like this
	- original data - 01101011
	- new data - 10010100
- The data is then converted back into text or byte data
- The data is then written to a new file with a modified name that is dependent on the original file, like this
	- original file - file.txt
	- new file - file_e.txt or file_d.txt
- After this, the original file is deleted and in the case of en/decrypting an entire directory, the original directory and all its files will be deleted
	- This deletion is the default behavior and if it is not wanted then the user can use `-k` to keep the original directory and its contents or a file
	
**Decryption**
- Since the data was flipped that process must be reversed first
	- The data is taken in and converted into binary
- The binary data is then flipped again
- The flipped data is then converted back into text or byte data
- The byte data is then unscrambled using the `vigenere cypher` like before
- Now the data is back in its original form
- The data is then written to a new file
- After this, the original file or directory is deleted unless the user says otherwise