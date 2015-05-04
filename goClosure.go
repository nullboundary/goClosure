package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
)

type closureRes struct {
	CompiledCode string        `json:"compiledCode,omitempty"`
	Errors       []codeError   `json:"errors,omitempty"`
	Warnings     []codeWarn    `json:"warnings,omitempty"`
	ServerErrors []serverError `json:"serverErrors,omitempty"`
	Statistics   codeStat      `json:"statistics,omitempty"`
}

type codeError struct {
	Charno    int    `json:"charno,omitempty"`
	ErrorMsg  string `json:"error,omitempty"`
	Lineno    int    `json:"lineno,omitempty"`
	File      string `json:"file,omitempty"`
	ErrorType string `json:"type,omitempty"`
	Line      string `json:"line,omitempty"`
}

type codeWarn struct {
	Charno   int    `json:"charno,omitempty"`
	Lineno   int    `json:"lineno,omitempty"`
	File     string `json:"file,omitempty"`
	WarnType string `json:"type,omitempty"`
	WarnMsg  string `json:"warning,omitempty"`
	Line     string `json:"line,omitempty"`
}

type serverError struct {
	Code     int    `json:"code,omitempty"`
	ErrorMsg string `json:"error,omitempty"`
}

type codeStat struct {
	OriginalSize   int `json:"originalSize,omitempty"`
	CompressedSize int `json:"compressedSize,omitempty"`
	CompileTime    int `json:"compileTime,omitempty"`
}

func main() {
	app := cli.NewApp()
	app.Name = "goClosure"
	app.Version = "0.1"
	app.Authors = []cli.Author{cli.Author{Name: "Noah Shibley", Email: "dev@pass.ninja"}}
	app.Usage = "Concat and minify via Google Closure Compiler API"
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Commands = []cli.Command{
		{
			Name:      "concat",
			ShortName: "c",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "none",
					Usage: "Changes the path root for the js files listed in the html input file. Syntax <oldPath>:<newPath>",
				},
				cli.StringFlag{
					Name:  "modify, m",
					Value: "none",
					Usage: "Changes input html file <script> tags replacing the many old js files with one concated new file",
				},
			},
			Usage:       "Read an html file and concat the js files in order",
			Description: "goClosure concat <input html file> <output js file>",
			Action:      concatCommand,
		},
		{
			Name:      "minify",
			ShortName: "m",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "advanced,a",
					Usage: "sets compilation level to advanced",
				},
				cli.BoolFlag{
					Name:  "whitespace,w",
					Usage: "sets compilation level to white space only",
				},
			},
			Usage:       "Minify one js file via Google Closure Compiler API",
			Description: "goClosure minify <inputfile> <outputfile>",
			Action:      minifyCommand,
		},
		{
			Name:      "all",
			ShortName: "a",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "none",
					Usage: "Changes the path root for the js files listed in the html input file. Syntax <oldPath>:<newPath>",
				},
				cli.StringFlag{
					Name:  "modify, m",
					Value: "none",
					Usage: "Changes input html file <script> tags replacing the many old js files with one concated new file",
				},
				cli.BoolFlag{
					Name:  "advanced,a",
					Usage: "sets compilation level to advanced",
				},
				cli.BoolFlag{
					Name:  "whitespace,w",
					Usage: "sets compilation level to white space only",
				},
			},
			Usage:       "Both Concat and Minify. Input is same as concat",
			Description: "goClosure all <input html file> <output js file>",
			Action:      allCommand,
		},
	}

	app.Run(os.Args)
}

//concatCommand reads an html file looking for js src and concats the js files in order found in the file
func concatCommand(c *cli.Context) {

	//require 1+ arguements
	if len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, "concat")
		return
	}

	//concat together various js files
	fmt.Println("Concating files...")
	jsFileData := concat(c.Args().First(), c.String("path"))

	//writeout to one big file
	err := ioutil.WriteFile(c.Args().Get(1), []byte(strings.Join(jsFileData, "\n")), 0644)
	if err != nil {
		println("Error writing file: ", err.Error())
		return
	}

	if c.String("modify") != "none" {
		modifyHtml(c.Args().First(), c.Args().Get(1), c.String("modify"))
	}

}

//minifyCommand minifies one js file using the Google Closure Compiler API
func minifyCommand(c *cli.Context) {

	//require 1 argument only
	if len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, "minify")
		return
	}

	//slurp whole js file
	jsFileData, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		println("Error loading file: ", err.Error())
		return
	}

	//send to be minified
	fmt.Println("Minifying file...")
	if c.Bool("advanced") && c.Bool("whitespace") {
		cli.ShowCommandHelp(c, "minify")
		return
	}
	compileLevel := "SIMPLE_OPTIMIZATIONS" //default
	switch {
	case c.Bool("advanced"):
		compileLevel = "ADVANCED_OPTIMIZATIONS"
	case c.Bool("whitespace"):
		compileLevel = "WHITESPACE_ONLY"
	}
	respData := minify(string(jsFileData), compileLevel)

	//save results
	err = ioutil.WriteFile(c.Args().Get(1), []byte(respData.CompiledCode), 0644)
	if err != nil {
		println("Error writing file: ", err.Error())
		return
	}

}

//allCommand performs both Concat and Minify. Input parameters are the same as concat
func allCommand(c *cli.Context) {

	if len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, "all")
		return
	}

	fmt.Println("Concating files...")
	jsFileData := concat(c.Args().First(), c.String("path"))

	//send to be minified
	fmt.Println("Minifying files...")
	if c.Bool("advanced") && c.Bool("whitespace") {
		cli.ShowCommandHelp(c, "all")
		return
	}
	compileLevel := "SIMPLE_OPTIMIZATIONS" //default
	switch {
	case c.Bool("advanced"):
		compileLevel = "ADVANCED_OPTIMIZATIONS"
	case c.Bool("whitespace"):
		compileLevel = "WHITESPACE_ONLY"
	}
	respData := minify(strings.Join(jsFileData, "\n"), compileLevel)

	//save results
	err := ioutil.WriteFile(c.Args().Get(1), []byte(respData.CompiledCode), 0644)
	if err != nil {
		println("Error writing file: ", err.Error())
		return
	}

	if c.String("modify") != "none" {
		modifyHtml(c.Args().First(), c.Args().Get(1), c.String("modify"))
	}

}

//minify connects to the closure compiler api and returns the minified result
func minify(jsData string, compileLevel string) closureRes {

	//build url
	apiUrl := "http://closure-compiler.appspot.com"
	resource := "/compile"
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("Error parsing URI: ", err.Error())
		os.Exit(1)
	}
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)
	fmt.Printf("connecting %s \n", urlStr)

	//encode url parameters
	data := url.Values{}
	data.Set("js_code", jsData)
	data.Add("compilation_level", compileLevel)
	data.Add("language", "ECMASCRIPT5_STRICT")
	data.Add("output_format", "json")
	data.Add("output_info", "compiled_code")
	data.Add("output_info", "errors")
	data.Add("output_info", "warnings")
	data.Add("output_info", "statistics")

	//make request
	client := http.Client{}
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	println(resp.Status)
	defer resp.Body.Close()
	result := closureRes{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("jsonRead Error: %v", err)
		os.Exit(1)
	}

	printResult(result)

	return result
	//return body

}

func concat(htmlFile string, pathFlag string) []string {

	jsDataList := make([]string, 10)

	//load html file
	htmlHandle, err := os.Open(htmlFile)
	if err != nil {
		fmt.Printf("Error loading html file: ", err.Error())
		os.Exit(1)
	}
	defer htmlHandle.Close()

	jsFiles := getJSFileList(htmlHandle)

	//if path flag used create new root path for each file
	var pathOld, pathNew string
	if pathFlag != "none" {
		pathReplace := pathFlag //get prefix replacement path old:new
		pathOldNew := strings.Split(pathReplace, ":")

		//make sure prefix replacement formatted old:new
		if len(pathOldNew) < 2 {
			//cli.ShowCommandHelp(c, "concat")
			fmt.Println("prefix replacement formatted old:new")
			os.Exit(1)
		}

		pathOld = pathOldNew[0]
		pathNew = pathOldNew[1]
	}

	//trim and join new path prefix
	for i, filename := range jsFiles {
		//make sure filename extension is js
		if path.Ext(filename) != ".js" {
			continue
		}

		//if path flag trim file names paths and replace with new
		if pathFlag != "none" {
			trimmedPath := strings.TrimPrefix(filename, pathOld)
			jsFiles[i] = path.Join(pathNew, trimmedPath)
		}

		println(jsFiles[i])
		fileData, err := ioutil.ReadFile(jsFiles[i]) //slurp whole file
		if err != nil {
			fmt.Printf("Error loading file: %s", err.Error())
			os.Exit(1)
		}

		jsDataList = append(jsDataList, string(fileData)) //concat files together
	}

	return jsDataList
}

func printResult(result closureRes) {

	if len(result.ServerErrors) > 0 {
		for i, _ := range result.ServerErrors {
			fmt.Println("-------------------------------------------")
			fmt.Printf("Server Error: %d \n", i)
			fmt.Printf("Error code: %d \n", result.ServerErrors[i].Code)
			fmt.Printf("Error msg: %s \n", result.ServerErrors[i].ErrorMsg)
		}
		fmt.Printf("closure server errors, process aborted")
		os.Exit(1)
	}

	if len(result.Errors) > 0 {
		for i, _ := range result.Errors {
			errorMsg := result.Errors[i]
			fmt.Println("-------------------------------------------")
			fmt.Printf("Error: %d \n", i)
			fmt.Printf("%s: %s at line %d character %d \n %s \n", errorMsg.ErrorType, errorMsg.ErrorMsg, errorMsg.Lineno, errorMsg.Charno, errorMsg.Line)
		}
		fmt.Printf("minify errors, process aborted")
		os.Exit(1)
	}

	for i, _ := range result.Warnings {
		warning := result.Warnings[i]
		fmt.Println("-------------------------------------------")
		fmt.Printf("Warning %d \n", i)
		fmt.Printf("%s: %s at line %d character %d \n %s \n", warning.WarnType, warning.WarnMsg, warning.Lineno, warning.Charno, warning.Line)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println("statistics")
	fmt.Printf("original size: %d KB\n", (result.Statistics.OriginalSize / 1024))
	fmt.Printf("compressed size: %d KB\n", (result.Statistics.CompressedSize / 1024))
	fmt.Printf("compile time: %ds \n", result.Statistics.CompileTime)

}

func modifyHtml(htmlFile string, jsFileName string, scriptPath string) {

	//load html file
	htmlHandle, err := os.Open(htmlFile)
	if err != nil {
		fmt.Printf("Error loading html file: ", err.Error())
		os.Exit(1)
	}
	defer htmlHandle.Close()

	doc, err := goquery.NewDocumentFromReader(htmlHandle)
	if err != nil {
		fmt.Printf("html parse error: ", err.Error())
		os.Exit(1)
	}

	var finalElm *goquery.Selection
	sel := doc.Find("script")

	for i := range sel.Nodes {
		single := sel.Eq(i)

		jsPath, _ := single.Attr("src")
		if strings.HasPrefix(jsPath, "http") || strings.HasPrefix(jsPath, "//") {
			finalElm = single.Parent() //get parent to append to
			continue
		}
		single.Remove()
	}

	jsFilePath := path.Join(scriptPath, jsFileName)

	//add a script tag. The dumb way. Can't get the right way to work
	finalElm.AppendHtml(`<script type="text/javascript" src="` + jsFilePath + `"></script>`)

	//writeout to one big file
	docString, err := doc.Html()
	if err != nil {
		println("Error rendering file: ", err.Error())
		os.Exit(1)
	}
	ext := path.Ext(htmlFile)
	outHtml := strings.TrimSuffix(htmlFile, ext) + ".min" + ext
	err = ioutil.WriteFile(outHtml, []byte(docString), 0644)
	if err != nil {
		println("Error writing file: ", err.Error())
		os.Exit(1)
	}

}

func getJSFileList(htmlFile *os.File) []string {

	jsList := make([]string, 5)

	doc, err := goquery.NewDocumentFromReader(htmlFile)
	if err != nil {
		fmt.Printf("html parse error: ", err.Error())
		os.Exit(1)
	}

	sel := doc.Find("script")
	for i := range sel.Nodes {
		single := sel.Eq(i)
		path, _ := single.Attr("src")
		if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "//") {
			continue
		}

		jsList = append(jsList, path)
	}
	return jsList
}
