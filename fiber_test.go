package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"
)

var app = fiber.New()

func TestRouteHelloWorld(t *testing.T) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	request := httptest.NewRequest("GET", "/", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World", string(bytes))

}

func TestCtx(t *testing.T) {
	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name", "Guest")
		return c.SendString("Hello " + name)
	})

	request := httptest.NewRequest("GET", "/hello?name=Dion", nil)
	response, err := app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Dion", string(bytes))
	request = httptest.NewRequest("GET", "/hello", nil)
	response, err = app.Test(request)

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Guest", string(bytes))
}

func TestHttpRequest(t *testing.T) {
	app.Get("/req", func(c *fiber.Ctx) error {
		first := c.Get("firstname")
		last := c.Cookies("lastname")
		return c.SendString("Hello " + first + " " + last)
	})
	request := httptest.NewRequest("GET", "/req", nil)
	request.Header.Set("firstname", "Eko")
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "Soegianto"})
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Eko Soegianto", string(bytes))
}

func TestRouteParameter(t *testing.T) {
	app.Get("/users/:userId/orders/:orderId", func(c *fiber.Ctx) error {
		userId := c.Params("userId")
		orderId := c.Params("orderId")
		return c.SendString("Order " + orderId + " from " + userId)
	})
	request := httptest.NewRequest("GET", "/users/eko/orders/2", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Order 2 from eko", string(bytes))
}
func TestFormReq(t *testing.T) {
	app.Post("/hi", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		return c.SendString("Hi " + name)
	})
	body := strings.NewReader("name=Eko")
	request := httptest.NewRequest("POST", "/hi", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hi Eko", string(bytes))
}

//go:embed source/contoh.txt
var contohFile []byte

func TestFormUpload(t *testing.T) {
	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		err = c.SaveFile(file, "./target/"+file.Filename)
		if err != nil {
			return err
		}
		return c.SendString("Upload Success")
	})
	body := new(bytes.Buffer)
	write := multipart.NewWriter(body)
	file, err := write.CreateFormFile("file", "contoh.txt")
	assert.Nil(t, err)

	file.Write(contohFile)
	write.Close()

	request := httptest.NewRequest("POST", "/upload", body)
	request.Header.Set("Content-Type", write.FormDataContentType())
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	by, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Upload Success", string(by))
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRequestBody(t *testing.T) {
	app.Post("/login", func(c *fiber.Ctx) error {
		body := c.Body()

		req := new(LoginReq)
		err := json.Unmarshal(body, req)

		if err != nil {
			return err
		}
		return c.SendString("Hi " + req.Username)
	})
	body := strings.NewReader(`{"username":"Eric","password":"rahasia"}`)
	request := httptest.NewRequest("POST", "/login", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hi Eric", string(bytes))
}

type RegisterReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password"  form:"password"`
	Name     string `json:"name" form:"name"`
}

func TestBodyParser(t *testing.T) {
	app.Post("/register", func(c *fiber.Ctx) error {
		req := new(RegisterReq)
		err := c.BodyParser(req)
		if err != nil {
			return err
		}
		return c.SendString("Register Success " + req.Username)
	})
}
func TestBodyParserJSON(t *testing.T) {
	TestBodyParser(t)
	body := strings.NewReader(`{"username":"Eric","password":"rahasia","name":"Eric Kunthady"}`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Eric", string(bytes))
}
func TestBodyParserForm(t *testing.T) {
	TestBodyParser(t)
	body := strings.NewReader(`username=Eric&password=rahasia&name=Eric+Kunthady`)
	request := httptest.NewRequest("POST", "/register", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Register Success Eric", string(bytes))
}
func TestResponseJSON(t *testing.T) {
	app.Get("/user", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"username": "khan",
			"name":     "Eko Khan",
		})
	})
	request := httptest.NewRequest("GET", "/user", nil)
	request.Header.Set("Accept", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"Eko Khan","username":"khan"}`, string(bytes))
}
func TestDownloadFile(t *testing.T) {
	app.Get("/download", func(c *fiber.Ctx) error {
		return c.Download("./source/contoh.txt", "contoh.txt")
	})
	request := httptest.NewRequest("GET", "/download", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, `attachment; filename="contoh.txt"`, response.Header.Get("Content-Disposition"))

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "this is sample text file for upload", string(bytes))
}

func TestRouteGroup(t *testing.T) {
	helloW := func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	}

	api := app.Group("/api")

	api.Get("/hello", helloW)
	api.Get("/world", helloW)

	web := app.Group("/web")

	web.Get("/hello", helloW)
	web.Get("/world", helloW)

	request := httptest.NewRequest("GET", "/web/hello", nil)
	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World", string(bytes))

}
