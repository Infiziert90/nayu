package main

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/argon2"
	_"golang.org/x/sys/cpu"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"src/database"
	"strings"
)

type Template struct {
	templates *template.Template
}

type RespStatus struct {
	Status string
}

type Resp struct {
	URL       string
	DeleteURL string
}

var (
	StatusOK     = RespStatus{"Success"}
	StatusFailed = RespStatus{"Failed"}
	db           = database.Dao
	hashedPW     = []byte{0, 6, 77, 79, 43, 92, 249, 222, 5, 108, 128, 233, 198, 50, 52, 30, 99, 130, 88, 232, 54, 252, 94, 214, 33, 172, 46, 156, 253, 215, 195, 171}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	db.Connect()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	t := &Template{templates: template.Must(template.ParseGlob("template/*.html"))}
	e.Renderer = t

	e.GET("/", index)
	e.GET("/*", find)
	e.GET("/upload", getUpload)
	e.POST("/upload", postUpload)
	e.GET("/delete", getDelete)
	e.POST("/delete", postDelete)
	e.GET("/delete/*", postDelete)
	e.File("/favicon.ico", "static/favicon.ico")
	e.Logger.Fatal(e.Start(":6776"))
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func postUpload(c echo.Context) error {
	pw := c.FormValue("password")
	if pw == "" || string(hashedPW) != string(argon2.IDKey([]byte(pw), []byte(""), 1, 64*1024, 4, 32)) {
		return c.JSON(http.StatusBadRequest, StatusFailed)
	}

	upload := database.Upload{
		ID:         bson.NewObjectId(),
		UniqueCode: database.CreateUCode(),
	}
	upload.DeleteCode = upload.UniqueCode + database.CreateDCode()

	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, StatusFailed)
	}

	src, err := file.Open()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, StatusFailed)
	}
	defer src.Close()

	end := strings.Split(file.Filename, ".")
	upload.File = fmt.Sprintf("files/%s.%s", upload.UniqueCode, end[len(end)-1])
	dst, err := os.Create(upload.File)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, StatusFailed)
	}
	defer dst.Close()

	io.Copy(dst, src)
	if _, err = io.Copy(dst, src); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, StatusFailed)
	}

	if err = db.Insert(&upload); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, StatusFailed)
	}

	r := Resp{
		fmt.Sprintf("%s/%s", c.Request().Host, upload.UniqueCode),
		fmt.Sprintf("%s/delete/%s", c.Request().Host, upload.DeleteCode),
	}

	return c.JSONPretty(http.StatusOK, &r, "  ")
}

func getUpload(c echo.Context) error {
	return c.Render(http.StatusOK, "upload.html", nil)
}

func getDelete(c echo.Context) error {
	return c.Render(http.StatusOK, "delete.html", nil)
}

func postDelete(c echo.Context) error {
	dcode := c.FormValue("dcode")
	if dcode == "" || len(dcode) != 12 {
		dcode = c.Request().RequestURI[8:]
		if len(dcode) != 12 {
			return c.JSON(http.StatusBadRequest, StatusFailed)
		}
	}

	upload, err := db.FindByDC(dcode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusFailed)
	}

	db.Delete(&upload)
	os.Remove(upload.File)

	return c.JSON(http.StatusOK, StatusOK)
}

func find(c echo.Context) error {
	url := c.Request().RequestURI[1:]
	if len(url) != 6 {
		return c.JSON(http.StatusBadRequest, StatusFailed)
	}

	item, err := db.FindByUC(url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, StatusFailed)
	}

	return c.File(item.File)
}
