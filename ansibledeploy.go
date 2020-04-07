package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	//"runtime"
	"strings"
)

const (
	VERSION = "0.1.0"
)

const (
	ANSIBLE_CMD = "/usr/bin/ansible-playbook"
)

const (
	retOK = iota
	retFailed
	retInvaidArgs
)

func execCmd(cmd string, shell bool) (out []byte, err error) {
	fmt.Printf("run command: %s", cmd)
	if shell {
		out, err = exec.Command("bash", "-c", cmd).Output()
	} else {
		out, err = exec.Command(cmd).Output()
	}
	return out, err
}

func checkExistFiles(files ...string) bool {
	for _, file := range files {
		file = strings.TrimSpace(file)
		if _, err := isExists(file); err != nil {
			fmt.Printf("[ERROR] check %s with %s.\n", file, err)
			return false
		}
	}
	return true
}

func isExists(file string) (ret bool, err error) {
	// equivalent to Python's `if not os.path.exists(filename)`
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false, err
	} else {
		return true, nil
	}
}

func doUpdateAction(action string, inventory_file string, operation_file string,
	version string, concurrent int) {

	loginfo := "doUpdateAction"
	if action != "update" || inventory_file == "" || operation_file == "" {
		fmt.Printf("Error parameters in %s\n", loginfo)
		os.Exit(retFailed)
	}

}

func doDeployAction(action string, inventory_file string, operation_file string,
	singlemode bool, concurrent int, retry_file string, ext_vars string, section string) {

	loginfo := "doDeployAction"
	if action != "deploy" || inventory_file == "" || operation_file == "" {
		fmt.Printf("Error parameters in %s", loginfo)
		os.Exit(retFailed)
	}
}

func main() {
	var err error

	single_mode := flag.Bool("s", false, "Single mode in deploy one host for observation.")
	concurrent := flag.Int("c", 1, "Process nummber for run the command at same time.")
	program_version := flag.String("V", "", "Module program version for deploy.")
	extra_vars := flag.String("e", "", "Extra vars for ansible-playbook.")
	section := flag.String("S", "", "Inventory section for distinguish hosts or tags.")
	retry_file := flag.String("r", "", "Retry file for ansible redo failed hosts.")
	inventory_file := flag.String("i", "", "Specify inventory host file.")
	operation_file := flag.String("f", "", "File name for module configure(yml format).")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Printf("%s: %s\n", os.Args[0], VERSION)
		os.Exit(retOK)
	}

	var action string = flag.Arg(0)

	if *operation_file == "" || *inventory_file == "" {
		fmt.Printf("[Error] operation and inventory file must provide.\n\n")
		flag.Usage()
		os.Exit(retFailed)
	} else {
		ret := checkExistFiles(*operation_file, *inventory_file)
		if !ret {
			fmt.Printf("[Error] check exists of operation and inventory file.\n")
			os.Exit(retInvaidArgs)
		}
	}

	if *operation_file, err = filepath.Abs(*operation_file); err != nil {
		panic(err)
	}

	if *inventory_file, err = filepath.Abs(*inventory_file); err != nil {
		panic(err)
	}

	fmt.Printf("[%s] action on [%s]\n", action, *operation_file)
	switch action {
	case "check":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		fmt.Println("check configure file.")
	case "update":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		fmt.Println("update code.")
		doUpdateAction(action, *inventory_file, *operation_file, *program_version, *concurrent)
	case "deploy":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		fmt.Println("deploy code.")
		doDeployAction(action, *inventory_file, *operation_file, *single_mode, *concurrent, *retry_file, *extra_vars, *section)
	case "rollback":
		fmt.Printf("-------------Now doing in action: %s\n", action)
		fmt.Println("rollback code.")
	default:
		fmt.Println("Not supported action: %s\n", action)
		os.Exit(retFailed)
	}
}