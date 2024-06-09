package handlers

import (
	"encoding/json"
	"iin_check/internal/db"
	"iin_check/internal/iin"
	"iin_check/internal/models"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "web/index.html")
	case "POST":
		iinID := r.FormValue("iin")
		http.Redirect(w, r, "/iin_check/"+iinID, http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func IinCheckHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iinID := vars["iin"]
	info, err := iin.ValidateIIN(iinID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func AddPersonHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "web/html/add_person.html")
	case "POST":
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		iinA := r.FormValue("iin")
		info, err := iin.ValidateIIN(iinA)
		if err != nil || !info.Correct {
			response := map[string]interface{}{
				"success": false,
				"errors":  "Invalid IIN",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		person := models.Person{
			Name:  name,
			IIN:   iinA,
			Phone: phone,
		}

		people, err := db.LoadDatabase()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		people = append(people, person)

		err = db.SaveDatabase(people)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
		}
		json.NewEncoder(w).Encode(response)
	}

}

func GetPersonByIINHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iinID := vars["iin"]

	info, err := iin.ValidateIIN(iinID)
	if err != nil || !info.Correct {
		response := map[string]interface{}{
			"success": false,
			"errors":  "Invalid IIN",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	people, err := db.LoadDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, person := range people {
		if person.IIN == iinID {
			json.NewEncoder(w).Encode(person)
			return
		}
	}

	http.Error(w, "Person not found", http.StatusNotFound)
}

func GetPersonsByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namePart := vars["name"]

	people, err := db.LoadDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []models.Person
	for _, person := range people {
		if containsIgnoreCase(person.Name, namePart) {
			result = append(result, person)
		}
	}

	json.NewEncoder(w).Encode(result)
}

func containsIgnoreCase(str, substr string) bool {
	str = strings.ToLower(str)
	substr = strings.ToLower(substr)
	return strings.Contains(str, substr)
}
