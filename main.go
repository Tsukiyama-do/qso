package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	task "./controllers"

	"./env"
	"github.com/gin-gonic/gin"

	//    simplejson "go-simplejson"
	"fmt"
	"strconv"

	"github.com/bitly/go-simplejson"
)

func init() {
	log.SetFlags(log.Lshortfile)

}

const DOWNLOADS_PATH = "./downloads/"
const UPLOADS_PATH = "./uploads/"

const TLS_CRT = "./private/https-jj1pow-com-003.crt"
const TLS_KEY = "./private/https-jj1pow-com-003.key"

func main() {

	//  パラメータチェック
	if len(os.Args) == 2 {
		if os.Args[1] == "release" {
			gin.SetMode(gin.ReleaseMode)
			env.S_mset(false) // means release mode
		}
	}
	//

	router := gin.Default()
	router.LoadHTMLGlob("views/*.tmpl")
	router.Static("/assets", "./public/assets")
	router.Static("/public", "./public")
	router.Static("/favicon.ico", "./favicon.ico")
	router.Static("/private", "./.ssh")

	// 初期画面の出力のためのGET処理
	router.GET("/", func(c *gin.Context) {
		controller := task.NewTask()
		tasks := controller.GetAll()

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Welcome to JJ1POW page!", //　追記 テンプレートにタスクを渡す
			"task":  tasks,
		})
	})

	// ダウンロード時のQSLカードのファイルダウンロード処理
	router.GET("/downloads/:filename", func(c *gin.Context) {
		fileName := c.Param("filename")
		//      log.Printf("Hi, %s", fileName)
		targetPath := DOWNLOADS_PATH + fileName
		_, err := os.Stat(targetPath)
		if err != nil {
			c.String(403, "I am sorry, QSL card no created yet!. Please send the message to my gmail.")
			log.Println("I am sorry, QSL card no created yet!. Please send the message to my gmail.")
			return
		}
		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "application/pdf")
		c.File(targetPath)
	})

	// ダウンロード時のQSLカードの一覧を出力する検索
	router.POST("/", func(c *gin.Context) {
		buf := make([]byte, 2048)
		num, _ := c.Request.Body.Read(buf)
		//        bjson := string(buf[0:n])
		//        log.Println(bjson)
		//        if len(bjson) >= 0 {      // dummy if branch
		//        }
		if env.S_mode() {
			log.Printf("request - num of buf : %d\n", num)
		}
		bjson, err := simplejson.NewJson(buf)
		if err != nil {
			log.Println("Error occured at JSON conversion.")
			c.JSON(http.StatusBadRequest, gin.H{"results": fmt.Sprintf("Error occured at JSON conversion")})
			return
		}

		// control operations
		controller := task.NewTask()
		sl_callsign := controller.SearchDB(bjson.Get("callsign").MustString(), bjson.Get("frdate").MustString(), bjson.Get("todate").MustString())

		var s_json_cal string
		s_json_cal = `[ `
		for _, items := range *sl_callsign {
			s_json_cal = s_json_cal + `{ "No." : "` + strconv.Itoa(items.ID) + `" , "Callsign":  "` + items.Callsign + `" , "Date":  "` + items.Datetime + `" , "File":  "` + items.Files + `" },`
		}

		for _, items := range *sl_callsign {
			if env.S_mode() {
				log.Printf(`{ "No." : "` + strconv.Itoa(items.ID) + `" , "Callsign":  "` + items.Callsign + `" , "Date":  "` + items.Datetime + `" , "File":  "` + items.Files + `" },`)
			}
		}

		if len(s_json_cal) > 0 { //  文字列末尾の　, を削除している。
			s_json_cal = string(s_json_cal[:(len(s_json_cal) - 1)])
		}

		s_json_cal = s_json_cal + `] `

		c.JSON(200, gin.H{"results": s_json_cal})

	})

	// QSLカードのアップロード時のPOST
	router.POST("/uploads/", func(c *gin.Context) {
		//    log.Println("POST uploads")
		// Source
		var targetPath string
		file, err := c.FormFile("qslupload")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"results": fmt.Sprintf("File Transfer Error: %s", err.Error())})
			log.Printf("File Transfer Error: %s", err.Error())
			return
		}

		filename := filepath.Base(file.Filename)
		targetPath = UPLOADS_PATH + filename

		// ファイル存在チェック　IsExist check to filename.
		for n := 0; n < 10; n++ {
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				break //  ループを抜ける
			} else {
				s_n := strconv.Itoa(n)
				targetPath = targetPath + ".bk-" + string(s_n)
				if env.S_mode() {
					log.Printf("targetPath : %s\n　", targetPath)
				}
			}
		}

		// Save an uploaded file to designated directory.
		if err = c.SaveUploadedFile(file, targetPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"results": fmt.Sprintf("File Saving Error: %s", err.Error())})
			log.Printf("File Saving Error: %s", err.Error())
			return
		}
		// register callsign to database
		s_callsignup, _ := c.GetPostForm("callsignup")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"results": fmt.Sprintf("Callsign Registration Error : %s", err.Error())})
			log.Printf("Callsign Registration Error 1: %s", err.Error())
			return
		}

		log.Println(s_callsignup)

		controller := task.NewTask()                      // control operations
		err = controller.UploadDB(s_callsignup, filename) // コールサインをDB登録する。
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"results": fmt.Sprintf("Callsign Registration DB Error : %s", err.Error())})
			log.Printf("Callsign Registration DB Error : %s", err.Error())
			return
		}

		// Success reply
		c.JSON(http.StatusOK, gin.H{"results": fmt.Sprintf("File %s uploaded successfully.", file.Filename)})
		log.Printf("POST file upload operations of %s is completed.", filename)

	}) // END of router.POST

	// router.Run(":8080")
	err := router.RunTLS(":8443", TLS_CRT, TLS_KEY)
	if err != nil {
		log.Fatal(err)
	}

}
