package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type PageManager struct {
	Pages     []QuestionPage `json:"pages"`
	PageCount int            `json:"pagecount"`
}

type QuestionPage struct {
	PageId      string          `json:"pageid"`
	PageName    string          `json:"pagename"`
	PageDisc    string          `json:"pagendisc"`
	PageBackImg string          `json:"pagebackimg"`
	Questions   []*QuestionInfo `json:"questions"`
	Status      StatusInfo      `json:"statusinfo"`
}

type QuestionInfo struct { // テンプレート展開用データ構造
	Id       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Lpx      string `json:"lpx"`
	Tpx      string `json:"tpx"`
	Rpos     string `json:"rpos"`
	Status   string `json:"status"`
}

type StatusInfo struct {
	ALL       int      `json:"statusall"`
	OK        int      `json:"statusok"`
	YET       int      `json:"statuyet"`
	PER       string   `json:"statuper"`
	StatusStr []string `json:"statusstr"`
}

// ============================================================================
// GLOBAL VARIABLES
// ============================================================================
// Question Difinition
const TOOLINFO_FILE string = "./DATA/questionpages.json"

// page manager
var pagemanager PageManager

// current page for work
var thispage QuestionPage

// current page questions
var messeges []*QuestionInfo

// ============================================================================
// Each Page View Creater
// ============================================================================
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
		// POST 保存処理
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

// ============================================================================
// Main
// ============================================================================
//-----------------------------------------------------------------------------
// init
//-----------------------------------------------------------------------------
func loadPageInfo() {
	pagemanager.PageCount = 1
	p := make([]QuestionPage, 1)
	pagemanager.Pages = p
	pagemanager.Pages[0] = thispage
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
