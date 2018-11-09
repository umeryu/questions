package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Page struct { // テンプレート展開用データ構造
	Title string
	Count int
}

type QuestionPage struct {
	PageId      string
	PageName    string
	PageBackImg string
	Questions   []*QuestionInfo
	Status      map[string]string
}

type QuestionInfo struct { // テンプレート展開用データ構造
	Id       string
	Question string
	Answer   string
	Lpx      string
	Tpx      string
	Rpos     string
	Status   string
}

var thispage QuestionPage
var messeges []*QuestionInfo

func questionView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	funcMap := template.FuncMap{
		"statussetfunc": func(id string) string {
			fmt.Println(id)
			return id
		},
	}
	// テンプレートをパース
	tmpl := template.Must(template.New("tmpl").Funcs(funcMap).ParseFiles("tmpl/question.html"))
	//tmpl := template.Must(template.ParseFiles("tmpl/hello.html"))

	if r.Method == "GET" {

	} else {
		fmt.Println("in POST ")

		err := r.ParseForm()
		if err != nil {
			fmt.Println("form parse erro")
		}

		cont := r.Form["status-info"]
		fmt.Println("=> ", cont)
		//  id, status
		// 現状保存し、Statusに保存
		// 現状の表示情報を取り出す。
		// 保存状態からデータを作り直しmessegeエンハンス

	}

	//messegesは外部から読み込み
	// テンプレートを描画 こっちはディレクトリ構成いらない
	if err := tmpl.ExecuteTemplate(w, "question.html", thispage); err != nil {
		log.Fatal(err)
	}

}

//-----------------------------------------------------------------------------
// init
//-----------------------------------------------------------------------------
func loadPageInfo() {
	thispage.PageId = "page1"
	thispage.PageName = "first page"
	thispage.PageBackImg = "./img/page1.png"

	createQuestionInfo()
	thispage.Questions = messeges

}

func createQuestionInfo() int {
	var count int = 0
	count = 4
	messeges = make([]*QuestionInfo, count)
	messeges[0] = &QuestionInfo{
		"myid1",
		"Question1",
		"Answer1",
		"100px",
		"200px",
		"L",
		"OK",
	}
	messeges[1] = &QuestionInfo{
		"myid2",
		"Question2",
		"世界",
		"300px",
		"300px",
		"R",
		"OK",
	}
	messeges[2] = &QuestionInfo{
		"myid3",
		"Question3",
		"世界",
		"500px",
		"500px",
		"R",
		"NG",
	}
	messeges[3] = &QuestionInfo{
		"myid4",
		"Question4",
		"世界",
		"100px",
		"500px",
		"L",
		"NG",
	}
	return count
}

//-----------------------------------------------------------------------------
// main
//-----------------------------------------------------------------------------
func main() {
	loadPageInfo()
	http.HandleFunc("/q", questionView)
	http.ListenAndServe(":8080", nil)
}
