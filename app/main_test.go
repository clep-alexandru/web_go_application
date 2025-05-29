package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestPaginaCreareContHandler_Success(t *testing.T) {

	mockDB, mock, _ := sqlmock.New()
	db = mockDB
	defer db.Close()

	username := "testuser"
	password := "parolatest"

	mock.ExpectExec("INSERT INTO users").
		WithArgs(username, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	form := "username=" + username + "&parola=" + password
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	PaginaCreareContHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Status incorect: %d, expected 303", resp.StatusCode)
	}
}

func TestPaginaLogareHandler_Success(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	db = mockDB
	defer db.Close()

	username := "testuser"
	password := "parolatest"

	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mock.ExpectQuery("SELECT password FROM users").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(string(hashed)))

	form := "username=" + username + "&password=" + password
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	PaginaLogareHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Status incorect la login: %d, expected 303", resp.StatusCode)
	}
}
