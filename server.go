package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type PageManager struct {
	Title       string         `json:"title"`
	Desc        string         `json:"desc"`
	Pages       []QuestionPage `json:"pages"`
	PageCount   int            `json:"pagecount"`
	TotalStatus StatusInfo     `json:"totalstatusinfo"`
}

type QuestionPage struct {
	PageId      string          `json:"pageid"`
	PageName    string          `json:"pagename"`
	PageDesc    string          `json:"pagendesc"`
	PageBackImg string          `json:"pagebackimg"`
	Questions   []*QuestionInfo `json:"questions"`
	Status      StatusInfo      `json:"statusinfo"`
	Comments    []string        `json:"comments"`
}

type QuestionInfo struct { // テンプレート展開用データ構造
	Id       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Lpx      string `json:"lpx"`
	Tpx      string `json:"tpx"`
	Ckind1   string `json:"ckind1"`
	Ckind2   string `json:"ckind2"`

	Status string `json:"status"`
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
const QUESTIONPAGES_FILE string = "./DATA/questionpages.json"
const QUESTIONPAGES_BAKUP_FILE string = "./DATA/questionpages.backup.json"

// page manager all data is loaded this object
var pagemanager PageManager

// for get each pagedata when called view. key is pageURL(this is created by PageId).
var pagedatamap map[string]*QuestionPage

// ============================================================================
// Each Page View Creater
// ============================================================================
// どのページもこの処理にくる。ページにあったデータを割り当てる必要がある.
// 対応ページデータ(nowpage)をpageurlとのマップ情報で差し替える
func questionView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	var nowPageID string = (r.RequestURI)[1:]
	var nowpage *QuestionPage = pagedatamap[nowPageID]

	//fmt.Println("url---->", nowPageID)
	//fmt.Println("check nowpage-->:", (nowpage).PageId, (nowpage).PageName)

	funcMap := template.FuncMap{
		"savefunc": func() string {
			fmt.Println("保存が呼ばれたよ!")
			savePageInfo()
			return "OK"
		},
	}
	// テンプレートをパース
	tmpl := template.Must(template.New("tmpl").Funcs(funcMap).ParseFiles("tmpl/question.html"))
	//tmpl := template.Must(template.ParseFiles("tmpl/question.html"))

	if r.Method == "GET" {

	} else {
		// POST 保存処理
		fmt.Println("in POST ")

		err := r.ParseForm()
		if err != nil {
			fmt.Println("form parse err")
		}

		// 取り出せるフォームはtextとか、Buttonから取り出しても何もないよ
		(*nowpage).Status.StatusStr = r.Form["status-info"]

		all := len((*nowpage).Status.StatusStr)
		okcount := 0
		r := regexp.MustCompile("-statusOK")

		for i := 0; i < all; i++ {
			//fmt.Println("--------> now in status check:", (*nowpage).Status.StatusStr[i])
			var state string
			var nowStr string = (*nowpage).Status.StatusStr[i]
			if r.MatchString(nowStr) {
				state = "OK"
				okcount++
			} else {
				state = "NG"
			}

			idsplit := strings.Split(nowStr, "-")

			//QuestionPageの各QuestionInfoのStatusをUPDATE
			var qi *QuestionInfo //pointer出ないと上書きされない
			// todo あとでマップに書きに書き換える
			for i := 0; i < len((*nowpage).Questions); i++ {
				qi = ((*nowpage).Questions)[i]
				if qi.Id == idsplit[0] {
					qi.Status = state
					//fmt.Println("update state---:"+idsplit[0], ":::", qi.Status)
					break
				}
			}
		}
		(*nowpage).Status.ALL = all
		(*nowpage).Status.OK = okcount
		(*nowpage).Status.YET = all - okcount
		(*nowpage).Status.PER = strconv.FormatFloat((float64(okcount)/float64(all))*100, 'f', 2, 64)

	}

	if err := tmpl.ExecuteTemplate(w, "question.html", nowpage); err != nil {
		log.Fatal(err)
	}
}

// ============================================================================
// Top Page View Creater
// ============================================================================
func topView(w http.ResponseWriter, r *http.Request) {
	// テンプレートをパース
	tmpl := template.Must(template.ParseFiles("tmpl/top.html"))

	//status 更新
	pagemanager.TotalStatus.ALL = pagemanager.PageCount
	var okcount int = 0
	for i := 0; i < len(pagemanager.Pages); i++ {
		if (pagemanager.Pages[i]).Status.PER == "100.00" {
			okcount++
		}
	}
	pagemanager.TotalStatus.OK = okcount
	pagemanager.TotalStatus.PER = strconv.FormatFloat((float64(okcount)/float64(pagemanager.PageCount))*100, 'f', 2, 64)

	if r.Method == "GET" {

	} else {
	}
	if err := tmpl.ExecuteTemplate(w, "top.html", pagemanager); err != nil {
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
	jsonStr, err := ioutil.ReadFile(QUESTIONPAGES_FILE)
	if err != nil {
		fmt.Printf("error:%s¥n", err)
		return
	}
	jsonbytes := ([]byte)(jsonStr)

	err = json.Unmarshal(jsonbytes, &pagemanager)
	if err != nil {
		fmt.Printf("error:%s¥n", err)
		return
	}
	//fmt.Println(pagemanager.Pages[0].PageName)
	//fmt.Println(pagemanager.Pages[0].Questions[0].Question)
}
func savePageInfo() {

	jsonBytes, err := json.MarshalIndent(pagemanager, "", "	")

	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}

	err2 := ioutil.WriteFile(QUESTIONPAGES_FILE, jsonBytes, os.ModePerm)
	if err2 != nil {
		fmt.Printf("error:%s¥n", err2)
		return
	} else {
		fmt.Println("## file saved to :", QUESTIONPAGES_FILE)
	}

}

//-----------------------------------------------------------------------------
// main
//-----------------------------------------------------------------------------
func main() {
	loadPageInfo()
	count := pagemanager.PageCount
	fmt.Println("## Qusetion Pages = ", count, " page exist.")
	pagedatamap = make(map[string]*QuestionPage, 1) //1は初期キャパシティ　追加により増える

	// Topハンドラ追加
	http.HandleFunc("/top", topView)
	fmt.Println("## Add Page Handlers=> top")
	// Pageハンドラ追加
	for i := 0; i < count; i++ {
		now := pagemanager.Pages[i]
		pageid := now.PageId
		//page data map作成 viewで利用
		pagedatamap[pageid] = &(pagemanager.Pages[i])

		http.HandleFunc("/"+pageid, questionView)
		fmt.Println("## Add Page Handlers=> pageID:", pageid, " Title:", now.PageName, " as /", pageid)
	}

	// Access可能なURLを設定
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	// Server起動
	PORT := ":8080"
	fmt.Println("### Server Start ### " + PORT)
	http.ListenAndServe(PORT, nil)
}
