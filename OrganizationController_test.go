package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go-rest-api/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	router := SetupRouter()

	orgName := "Autotest Org"
	jsonData := map[string]string{"name": orgName}
	jsonValue, _ := json.Marshal(jsonData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(jsonValue))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var organization models.Organization
	if err := json.Unmarshal(w.Body.Bytes(), &organization); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, orgName, organization.Name)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/organizations/"+strconv.Itoa(int(organization.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	if err := json.Unmarshal(w.Body.Bytes(), &organization); err != nil {
		assert.Fail(t, err.Error())
	}

	assert.Equal(t, orgName, organization.Name)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/organizations/"+strconv.Itoa(int(organization.ID)), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
