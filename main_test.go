package main

import (
	"OMUS/server/controllers"
	helper "OMUS/server/helpers"
	"OMUS/server/seed"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

var server_test = controllers.Server{}

func TestMain(m *testing.M) {
	// INIT TEST SERVER DB
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server_test.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	seed.Load(server_test.DB)
	server_test.Run(":8080")

	code := m.Run()

	clearTable()
	os.Exit(code)

}

func clearTable() {
	server_test.DB.Exec("DELETE FROM public.urls")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/urls", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addURLs(3)

	req, _ := http.NewRequest("GET", "/urls", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server_test.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addURLs(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		server_test.DB.Exec("INSERT INTO urls(original_url, encoded_url, visits_counter) VALUES($1, $2, $3)", "URL"+strconv.Itoa(i), helper.Encode(uint64(i)), i)
	}
}
