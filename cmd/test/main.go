package main

import (
	// "bufio"
	// "fmt"
	// "io"
	// "os"
	// "os/exec"
	// "strings"

	"flag"
	"fmt"
	"regexp"
	"time"

	"github.com/google/goexpect"
	"github.com/google/goterm/term"
)

const (
	timeout = 10 * time.Second
)

var (
	passExpect = flag.String("pass", "Password:", "pass expect")
)

func main() {
	// ./gaiad keys list --log_level info --keyring-backend file --output json

	// commandName := "./gaiad"
	// var params []string
	// params = []string{"keys", "list", "--log_level", "info", "--keyring-backend", "file", "--output", "json"}

	// cmd := exec.Command(commandName, params...)
	// fmt.Println("exec ", cmd.Args)

	// https://vimsky.com/examples/usage/golang_os_exec_Cmd_Output.html
	// https://blog.csdn.net/weixin_43223076/article/details/85083381

	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道&#xff0c;一般不需要显式的关闭该管道。
	// stdout, err := cmd.StdoutPipe()

	// if err != nil {
	// 	fmt.Println("cmd.StdoutPipe err: ", err)
	// 	return
	// }

	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	fmt.Println("cmd.StdinPipe err: ", err)
	// 	return
	// }

	// cmd.Stderr = os.Stderr
	// // cmd.Dir = dir
	// err = cmd.Start()
	// // err = cmd.Run()
	// if err != nil {
	// 	fmt.Println("cmd.Start err: ", err)
	// 	return
	// }

	// // stdoutStderr, err := cmd.Output()
	// // if err != nil {
	// // 	fmt.Println("cmd.Output err: ", err)
	// // }
	// // fmt.Printf("%s\n", stdoutStderr)

	// //创建一个流来读取管道内内容&#xff0c;这里逻辑是通过一行一行的读取的
	// reader := bufio.NewReader(stdout)

	// // reader := bufio.NewReader(stdin)

	// // for {
	// // 	line, err2 := cmd.Output()
	// // 	if err2 != nil || io.EOF == err2 {
	// // 		fmt.Println("Sstdout EOF!")
	// // 		break
	// // 	}
	// // 	l := string(line)
	// // 	if strings.EqualFold(l, "EOF") {
	// // 		fmt.Println("Need key: ", line)
	// // 	} else {
	// // 		fmt.Println("Readed: ", line)
	// // 	}
	// // }

	// //实时循环读取输出流中的一行内容
	// for {
	// 	stdin.Write([]byte("kerberos"))
	// 	stdin.Write([]byte("\r\n"))
	// 	line, err2 := reader.ReadString('\n')

	// 	if strings.EqualFold(line, "EOF") {
	// 		fmt.Println("Need key: ", line)
	// 	} else {
	// 		fmt.Println("Readed: ", line)
	// 	}

	// 	if err2 != nil || io.EOF == err2 {
	// 		fmt.Println("Sstdout EOF!")
	// 		break
	// 	}
	// }
	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Println("cmd.Wait err: ", err)
	// }

	// https://github.com/google/goexpect

	// go run main.go --pass asd
	// go run main.go --pass "Enter keyring passphrase:"

	flag.Parse()
	fmt.Println(term.Bluef("expect 1 example"))

	passRE := regexp.MustCompile(*passExpect) // "Enter keyring passphrase:")

	e, _, err := expect.Spawn("./gaiad keys list --log_level info --keyring-backend file --output json", -1)
	if err != nil {
		fmt.Println("expect.Spawn err: ", err)
	}
	defer e.Close()

	_, _, err = e.Expect(passRE, timeout)
	if err != nil {
		fmt.Println("expect.Expect err: ", err)
		return
	}
	e.Send("kerberos\n")

	promptRE := regexp.MustCompile("%")
	result, _, _ := e.Expect(promptRE, timeout)

	fmt.Println(term.Greenf("result: %s\n", result))

	return
}
