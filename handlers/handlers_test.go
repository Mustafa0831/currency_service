package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"currency_service/handlers"
	"currency_service/models"
	"currency_service/utils"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *sqlx.DB

type Config struct {
	Port         string `json:"port"`
	DBConnection string `json:"db_connection"`
}

func setupTest() {
	var err error
	fmt.Println("Starting setupTest")

	config := Config{
		Port:         "8080",
		DBConnection: "sqlserver://kursUser:KursPswd@123@localhost:1433?database=TEST",
	}
	fmt.Printf("Config loaded: %+v\n", config)

	testDB, err = utils.ConnectDB(config.DBConnection)
	if err != nil {
		log.Fatalf("Could not connect to test database: %v", err)
	}

	_, err = testDB.Exec(`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='R_CURRENCY' and xtype='U')
    CREATE TABLE R_CURRENCY (
        ID INT IDENTITY(1,1) PRIMARY KEY,
        TITLE NVARCHAR(255) NOT NULL,
        CODE NVARCHAR(10) NOT NULL,
        VALUE FLOAT NOT NULL,
        A_DATE DATE NOT NULL
    )`)
	if err != nil {
		log.Fatalf("Could not create test table schema: %v", err)
	}

	handlers.InitHandlers(testDB)
}

func clearTestDB(t *testing.T) {
	_, err := testDB.Exec("DELETE FROM R_CURRENCY")
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	setupTest()
	code := m.Run()
	os.Exit(code)
}

func TestLogOutput(t *testing.T) {
	fmt.Println("This is a test log output")
	log.Println("This is another test log output")
}

func TestSaveCurrency(t *testing.T) {
	clearTestDB(t)

	router := mux.NewRouter()
	router.HandleFunc("/currency/save/{date}", handlers.SaveCurrency).Methods("GET")

	req, err := http.NewRequest("GET", "/currency/save/01.01.2024", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]bool
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"])

	time.Sleep(2 * time.Second)

	var count int
	err = testDB.Get(&count, "SELECT COUNT(*) FROM R_CURRENCY")
	require.NoError(t, err)
	assert.Greater(t, count, 0)
}

func TestGetCurrency(t *testing.T) {
	clearTestDB(t)

	testCurrency := models.Currency{
		Title: "US Dollar",
		Code:  "USD",
		Value: 1.0,
		ADate: "2024-01-01",
	}
	_, err := testDB.NamedExec(`INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (:TITLE, :CODE, :VALUE, :A_DATE)`, &testCurrency)
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/currency/{date}/{code}", handlers.GetCurrency).Methods("GET")

	req, err := http.NewRequest("GET", "/currency/01.01.2024/USD", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var currencies []models.Currency
	err = json.NewDecoder(rr.Body).Decode(&currencies)
	require.NoError(t, err)

	assert.Len(t, currencies, 1)
	assert.Equal(t, "US Dollar", currencies[0].Title)
	assert.Equal(t, "USD", currencies[0].Code)
	assert.Equal(t, 1.0, currencies[0].Value)

	expectedDate := "2024-01-01"
	actualDate := currencies[0].ADate[:10]
	assert.Equal(t, expectedDate, actualDate)
}
