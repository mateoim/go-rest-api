package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-rest-api/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var router *gin.Engine
var organization models.Organization

func TestMain(m *testing.M) {
	router = SetupRouter()

	jsonData := map[string]string{"name": "Autotest Org"}
	jsonValue, _ := json.Marshal(jsonData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &organization)

	exitVal := m.Run()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/organizations/"+strconv.Itoa(int(organization.ID)), nil)
	router.ServeHTTP(w, req)

	os.Exit(exitVal)
}

func TestCreateUser(t *testing.T) {
	jsonData := `{"first-name": "Auto", "last-name": "Tester", "email": "test@mail.com", "organization": ` + strconv.Itoa(int(organization.ID)) + `}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(jsonData)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var user models.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		assert.Fail(t, err.Error())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/"+strconv.Itoa(int(user.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var userGet models.User
	if err := json.Unmarshal(w.Body.Bytes(), &userGet); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, user.ID, userGet.ID)
	assert.Equal(t, user.FirstName, userGet.FirstName)
	assert.Equal(t, user.OrganizationID, userGet.OrganizationID)
	assert.Equal(t, user.Email, userGet.Email)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/users/"+strconv.Itoa(int(user.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCreateUserInvalidEmail(t *testing.T) {
	jsonData := `{"first-name": "Auto", "last-name": "Tester", "email": "test.com", "organization": ` + strconv.Itoa(int(organization.ID)) + `}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(jsonData)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserInvalidOrganization(t *testing.T) {
	jsonData := `{"first-name": "Auto", "last-name": "Tester", "email": "test@mail.com", "organization": -1}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(jsonData)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUsers(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/organizations/"+strconv.Itoa(int(organization.ID)), nil)
	router.ServeHTTP(w, req)
}
