package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	tag *string
	dir *string
	csv string
)

const (
	outputFolder = "locales"
)

func checkFileExt(_path string) error {
	if ext := path.Ext(_path); ext == ".tmpl" || ext == ".tpl" {
		return nil
	}

	return fmt.Errorf("%s is not a template file", _path)
}

func parseFile(_path string) ([]string, error) {
	var (
		transaltedRegexp   = "\\{\\{" + *tag + "\\s\\S+\\s\\S+\\}\\}"
		unTransaltedRegexp = "\\{\\{" + *tag + "\\s\\S+\\}\\}"
	)

	content, err := ioutil.ReadFile(_path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	re, _ := regexp.Compile(transaltedRegexp)
	result := re.FindAllString(string(content), -1)

	re, _ = regexp.Compile(unTransaltedRegexp)
	result = append(result, re.FindAllString(string(content), -1)...)

	return result, nil
}

func findTranslations(line string) (string, string) {
	result :=
		strings.Split(
			strings.TrimSpace(
				strings.TrimSuffix(
					strings.TrimPrefix(
						line,
						"{{"+*tag,
					),
					"}}",
				),
			),
			"\" \"",
		)

	if len(result) == 1 {
		result = append(result, "")
	}

	return strings.Replace(result[0], "\"", "", -1), strings.Replace(result[1], "\"", "", -1)
}

func writeFile(result map[string]string) error {
	for key, val := range result {
		csv += key + "," + val + "\n"
	}

	return nil
}

func translator(_path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	err = checkFileExt(_path)
	if err == nil {

		lines, err := parseFile(_path)
		if err != nil {
			return err
		}

		var result = make(map[string]string)

		for _, line := range lines {
			key, value := findTranslations(line)

			if !strings.HasPrefix(key, ".") {
				result[key] = value
			}
		}

		err = writeFile(result)
		if err != nil {
			return err
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func main() {
	tag = flag.String("t", "t", "Translation tag name.")
	dir = flag.String("d", "", "Root directory from where to start parsing.")
	flag.Parse()
	fmt.Println(*dir)

	if *dir == "" {
		fmt.Println("Please set `-d` flag pointing to directory from where to start parsing!")
		flag.PrintDefaults()
	}

	if err := filepath.Walk(*dir, translator); err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile("main.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(csv)
}
