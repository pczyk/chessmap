package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

const TemplateFile = "chart.html"
const PassportUrlTemplate = "https://www.chess.com/awards/%s/passport"
const DataCodeRegEx = `data-code=\"(.*?)\"`

var pageTemplate *template.Template
var rgx *regexp.Regexp

// Holds content for filling the HTML template
type Data struct {
	Player    string
	Countries string
}

func main() {
	initialize()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initialize() {
	tmpl, err := template.ParseFiles(TemplateFile)
	if err != nil {
		log.Fatalf("Template file '%s' is missing!", TemplateFile)
	}
	pageTemplate = tmpl

	rgx = regexp.MustCompile(DataCodeRegEx)
}

func handler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.String()[1:]
	codes := retrieveCountryCodes(player)
	countriesString := createCountriesString(codes)

	data := Data{player, countriesString}

	err := pageTemplate.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func retrieveCountryCodes(player string) []string {
	countryCodes := make([]string, 0)

	url := fmt.Sprintf(PassportUrlTemplate, player)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return countryCodes
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return countryCodes
	}

	rs := rgx.FindAllStringSubmatch(string(body), -1)

	for _, v := range rs {
		countryCodes = append(countryCodes, v[1])
	}

	return countryCodes
}

func createCountriesString(codes []string) string {
	var countries strings.Builder
	countries.WriteString("['Country', 'Stamp']")
	for _, c := range codes {
		countries.WriteString(fmt.Sprintf(",['%s','1']", c))
	}

	return countries.String()
}
