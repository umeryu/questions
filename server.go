package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Page struct { // テンプレート展開用データ構造
	Title string
	Count int
}

type QuestionPage struct {
	PageId      string
	PageName    string
	PageDisc    string
	PageBackImg string
	Questions   []*QuestionInfo
	Status      StatusInfo
}

type StatusInfo struct {
	ALL       int
	OK        int
	YET       int
	PER       string
	StatusStr []string
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

		//  id-status:OK id-status:NG の文字列を取得できる
		thispage.Status.StatusStr = r.Form["status-info"]
		//fmt.Println("=> ", thispage.Status.StatusStr, "len:", len(thispage.Status.StatusStr))
		all := len(thispage.Status.StatusStr)
		okcount := 0
		r := regexp.MustCompile("-status:OK")

		for i := 0; i < all; i++ {
			if r.MatchString(thispage.Status.StatusStr[i]) {
				okcount++
			}
		}
		thispage.Status.ALL = all
		thispage.Status.OK = okcount
		thispage.Status.YET = all - okcount
		thispage.Status.PER = strconv.FormatFloat((float64(okcount)/float64(all))*100, 'f', 2, 64)

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
	thispage.PageName = "背中：僧帽筋"
	thispage.PageBackImg = "/img/page1.png"
	thispage.PageDisc = "僧帽筋の起点、停止点を述べよ。"
	createQuestionInfo()
	thispage.Questions = messeges

}

func createQuestionInfo() int {
	var count int = 0
	count = 6
	messeges = make([]*QuestionInfo, count)
	messeges[0] = &QuestionInfo{
		"myid1",
		"起点1 2つ",
		"仙骨、腸骨稜",
		"300px",
		"450px",
		"L",
		"NG",
	}
	messeges[1] = &QuestionInfo{
		"myid2",
		"起点2　1",
		"腰椎L1-L5",
		"300px",
		"400px",
		"R",
		"NG",
	}
	messeges[2] = &QuestionInfo{
		"myid3",
		"起点3 １つ",
		"肋骨 9-12",
		"310px",
		"300px",
		"R",
		"NG",
	}
	messeges[3] = &QuestionInfo{
		"myid4",
		"起点4 1つ",
		"肩甲骨下角",
		"310px",
		"250px",
		"L",
		"NG",
	}
	messeges[4] = &QuestionInfo{
		"myid5",
		"起点5 1つ",
		"大後頭骨隆起",
		"250px",
		"80px",
		"L",
		"NG",
	}
	messeges[5] = &QuestionInfo{
		"myid6",
		"終点 3つ",
		"肩峰、肩甲骨突起、鎖骨外側1/3",
		"330px",
		"200px",
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
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	http.ListenAndServe(":8080", nil)
}
