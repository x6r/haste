package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	r "github.com/parnurzeal/gorequest"
	"github.com/x6r/haste/cmd"
	"golang.design/x/clipboard"
)

type Data struct {
	Key string `json:"key,omitempty"`
}

const (
	version     = "1.1.4"
	apiEndpoint = "documents"
)

var (
	instanceUrl  = flag.String("i", "https://p.x4.pm", "hastebin instance where the content will be uploaded")
	filePath     = flag.String("f", "", "upload a file")
	returnRaw    = flag.Bool("r", false, "returns raw url")
	printVersion = flag.Bool("v", false, "prints program version")
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	var content, ext string
	if *filePath == "" {
		content = flag.Arg(0)
	} else {
		file, err := os.ReadFile(*filePath)
		if err != nil {
			log.Fatalln(err)
		}
		content = string(file)
		slice := strings.Split(*filePath, ".")
		ext = slice[len(slice)-1]
	}

	if content == "" {
		if isPipe() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			content = scanner.Text()
		} else {
			content, *instanceUrl, *returnRaw = cmd.Execute()
		}
	}

	res := upload(*instanceUrl, content)

	if ext != "" {
		res += "." + ext
	}

	clipboard.Write(clipboard.FmtText, []byte(res))
	fmt.Println("copied to clipboard:", res)
}

func upload(url string, s string) (res string) {
	req := r.New()
	resp, body, errs := req.Post(fmt.Sprintf("%s/%s", url, apiEndpoint)).Type("text").Send(s).End()

	if errs != nil {
		log.Fatalln(errs)
	}

	if resp.StatusCode != 200 {
		log.Fatalln(resp.Status)
	}

	var data Data
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Fatalln(err)
	}

	if *returnRaw {
		res = fmt.Sprintf("%s/raw/%s", url, data.Key)
	} else {
		res = fmt.Sprintf("%s/%s", url, data.Key)
	}

	return
}

func isPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}
