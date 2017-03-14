package main

/*
* UDC, USETYPE, MACINF, MACCT, USER_NOTES
*/

import(
	"bufio"
	"fmt"
	"os"
	"strings")

const dl string = "::" // delimiter; separates UDC from machine info
const nl string = "|" // newline; gets converted to \n when reading with pmac

const sroot string = "C:/udp/c/" // local source root

func prt_str(udc string, raw_str string){
	fmt.Printf("UDCMT: %s: %s\n", strings.ToUpper(strings.Split(raw_str, "::")[0]), strings.Split(raw_str, "::")[1])
}

// get machine type info
func get_mactype(fp string, udc string){
	uStr := strings.Split(udc, "-")[0] // split udc at `-` character, to get the mactype code
	f, err := os.Open(fp)
	if err != nil{
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f) // create new scanner to read lines from file

	for scanner.Scan(){
		if strings.Contains(scanner.Text(), "##") == false{ // if not a comment in the .udc file
			if strings.Contains(scanner.Text(), uStr){
				prt_str(uStr, scanner.Text())
				f.Close()
				break
			}
		}
	}

	if err := scanner.Err(); err != nil{ // final error check
		panic(err)
	}
	defer f.Close()

	f.Close()
}

func main(){
	fmt.Printf("\n") // for space
	udcF := sroot + "src2/udc/test.txt" // local.udp (currently using test file until program is working correctly)
	oinF := sroot + "src2/udc/oinf.udc" // other info file

	list := [5]string {"UDC", "USETYPE", "MACINF", "MACCT", "USER_NOTES"} // categories

	reader := bufio.NewReader(os.Stdin) // get user input

	m := make(map[string]string) // init mapping

	// iterate over list and get info for each element
	for i := 0; i < len(list); i++{
		for{
			// essentially a `while true` loop, which will run until usable data is input for m[list[i]]
			fmt.Printf("%s:\t", list[i])
			m[list[i]], _ = reader.ReadString('\n')
			m[list[i]] = strings.TrimSpace(m[list[i]])
			if m[list[i]] != ""{
				break
			}else{
				continue
			}
		}
	}

	fmt.Printf("\n") // spacing
	get_mactype(oinF, m[list[0]])
	
	for i := 0; i < len(list); i++{
		fmt.Printf("%s\n", m[list[i]]) // print all data gathered data
	}

	// ask user if the provided info is correct
	for{
		fmt.Print("Is the above data correct? [n] ")
		c, _ := reader.ReadString('\n')
		c = strings.TrimSpace(c)
		if strings.Compare(strings.ToLower(c), "y") == 0{
			break
		}else if strings.Compare(strings.ToLower(c), "n") == 0{
			fmt.Println("Tool Aborting!")
			os.Exit(0)
		}else{
			continue
		}
	}

	// format information -- adding two \n characters effectively newlines when writing to file ... nevermind
	var mStr string = m[list[0]] + dl
	mStr = mStr + m[list[1]] + nl
	mStr = mStr + m[list[2]] + nl
	mStr = mStr + m[list[3]] + nl
	mStr = mStr + m[list[4]] + "\n"

	file, err := os.OpenFile(udcF, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(mStr); err != nil{
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString("\n"); err != nil{ // this doesn't work, either :/
		panic(err)
	}
	fmt.Println("File successfully updated!")
}