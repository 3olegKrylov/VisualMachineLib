package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var Count int
var Token string

//struct for our vm, you can check here https://mholt.github.io/json-to-go/ for convert json to struct of golang
type VMs []struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	RAM  int     `json:"ram"`
	CPU  float64 `json:"cpu"`
	Ssd  int     `json:"ssd"`
	Sata string  `json:"sata"`
}
func preparation(){
	client := &http.Client{}
	URL := "https://t3.linxdatacenter.com/api/v1/auth/" //http of auth
	resp, _ := http.NewRequest("GET", URL, nil)
	resp.SetBasicAuth("test3", "test88") //basic auth
	res, _ := client.Do(resp)

	//take token
	body, err := ioutil.ReadAll(res.Body)
	Token = string(body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	fmt.Printf("%s\n", body)
	//--------------------------------------
	//take information(json) about virtuals mashines
	URLvm := "https://t3.linxdatacenter.com/api/v1/vm/"
	client2 := &http.Client{}
	site, _ := http.NewRequest("GET", URLvm, nil)
	site.Header.Add("x-auth", string(body)) //add token for API
	siteIN, _ := client2.Do(site)

	body2, err := ioutil.ReadAll(siteIN.Body)

	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	fmt.Printf("string that have after Request: %s\n", body2)

	json_parse(string(body2), &ParseVM)
	fmt.Println("ParseVM без указателя:", ParseVM)
	fmt.Println("ParseVM c указателя:", *ParseVM)
	fmt.Println("ParseVM если пройтись range:")
	for k, v := range *ParseVM {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}
}

type SimpleVM struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	RAM  int     `json:"ram"`
	CPU  float64 `json:"cpu"`
	Ssd  int     `json:"ssd"`
	Sata string  `json:"sata"`
}

//Просто ф-я для более удобного парсинга
func json_parse(Data string, obj interface{}) {
	var b_Data = []byte(Data)
	err := json.Unmarshal(b_Data, obj)
	if err != nil {
		log.Println("error:", err)
	}
}

/*func CreateVMs(Data string) *VMs{
	var var_VMs = &VMs{}
	json_parse(Data, var_VMs)
	fmt.Println(var_VMs)
	return var_VMs
}*/

var ParseVM = &VMs{}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", *ParseVM)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "write", nil)
}
func StrToInt(str string) (int, error) { //
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

func SaveInAPI() {
	URLvm := "https://t3.linxdatacenter.com/api/v1/vm/"
	client2 := &http.Client{}

	j, _ := json.Marshal(*ParseVM)
	site, _ := http.NewRequest("POST", URLvm, bytes.NewBuffer(j))
	site.Header.Add("x-auth", string(Token)) //add token for API
	_, _ = client2.Do(site)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	//fmt.Println(id," is id that save")
	name := r.FormValue("name")
	ram, _ := StrToInt(r.FormValue("ram"))
	cpu, _ := strconv.ParseFloat(r.FormValue("cpu"), 64)
	ssd, _ := StrToInt(r.FormValue("ssd"))
	sata := r.FormValue("sata")

	//fmt.Println(name," ", ram," ", cpu," ", ssd," ", sata)
	if id != "" {
		fmt.Println("here")
		idi, _ := StrToInt(id)

		NewVM := &VMs{{
			ID:   idi,
			Name: name,
			RAM:  ram,
			CPU:  cpu,
			Ssd:  ssd,
			Sata: sata,
		}}


		URLvm := "https://t3.linxdatacenter.com/api/v1/vm/"
		client2 := &http.Client{}

		j, _ := json.Marshal(*NewVM)
fmt.Println("bytes.NewBuffer(j", bytes.NewBuffer(j))
		site, _ := http.NewRequest("POST", URLvm, bytes.NewBuffer(j))
		site.Header.Add("x-auth", string(Token)) //add token for API
		_, _ = client2.Do(site)

		preparation()

	} else {
		Count++
		NewVM := SimpleVM{
			ID:   Count,
			Name: name,
			RAM:  ram,
			CPU:  cpu,
			Ssd:  ssd,
			Sata: sata,
		}
		fmt.Println(NewVM)
		m := &VMs{}
		*m = append(*ParseVM, NewVM)
		ParseVM = m
		SaveInAPI()
	}

	http.Redirect(w, r, "/", 302)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	id, _ := StrToInt(r.FormValue("id"))
	M := SimpleVM{}
	for _, v := range *ParseVM {
		if v.ID == id {
			M = v
		}
	}

	fmt.Println(M)
	t.ExecuteTemplate(w, "write", M)

}
func main() {

	preparation()

	Count = len(*ParseVM)
	fmt.Println(Count)

	//client for take the token for API

	fmt.Println("listening on port :3000")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/SavePost", savePostHandler)
	http.ListenAndServe(":3000", nil)

}
