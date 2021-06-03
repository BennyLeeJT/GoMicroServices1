package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// this will be used for POST, PUT,
type courseInfo struct {
	Title string `json:"Title"`
}

// used for storing courses on the REST API
var courses map[string]courseInfo

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the REST API!")
}

// this function is to add checking for access token/key
// this is literal checking for simplicity
func validKey(w http.ResponseWriter, r *http.Request) bool { // also can return flag value, the value to be checked in the func course
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else { // if "key" is provided in query string but is not found to be matching any key stored, invalid key
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("401 - Please provide valid access token/key"))
			return false
		}
	} else { // this else is if "key" is not provided in the query string
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Please provide access token/key2"))
		return false
	}
}

func allcourses(w http.ResponseWriter, r *http.Request) {

	// first thing first, check for access token/key when user wants to see all current courses
	if !validKey(w, r) {
		return
	}

	keyValue := r.URL.Query()

	for key, value := range keyValue {
		fmt.Printf("key, value : %v %v\n", key, value)
	}

	fmt.Fprintf(w, "courses : %v\n\n", courses)
	// returns all the courses in JSON
	fmt.Fprintf(w, "List of all current courses updated as follows :\n")
	json.NewEncoder(w).Encode(courses)
}

func course(w http.ResponseWriter, r *http.Request) {

	// first thing first, check for access token/key when user wants to acccess a specific course id
	if !validKey(w, r) {
		return
	}

	params := mux.Vars(r)

	// Retrieving Courses Info
	if r.Method == "GET" {
		fmt.Fprintf(w, "Displaying All Current Courses : %v\n", courses)
		fmt.Fprintf(w, "Course ID requested to be search in params : %v\n", params)
		if _, ok := courses[params["courseid"]]; ok {
			fmt.Fprintf(w, "Course Info for %v : ", params["courseid"])
			json.NewEncoder(w).Encode(courses[params["courseid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	} else if r.Method == "DELETE" { // Deleting courses, if course id exists
		if _, ok := courses[params["courseid"]]; ok {
			delete(courses, params["courseid"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Course deleted: " +
				params["courseid"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new course
		if r.Method == "POST" {
			// read the string sent to the service
			var newCourse courseInfo

			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newCourse)

				if newCourse.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply course " +
						"information " + "in JSON format")) // 3 strings to make it dynamic, if we want to reply info or json with other content, also can replace with a function to make it dynamic request from user.
					return
				}

				// check if course exists;
				if _, ok := courses[params["courseid"]]; !ok {
					courses[params["courseid"]] = newCourse
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["courseid"]))
				} else { // add only if course does not exist
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte("409 - Duplicate course ID. Course not added."))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information " +
					"in JSON format"))
			}
		} else if r.Method == "PUT" { // PUT is for creating or updating existing course
			var newCourse courseInfo

			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newCourse) // convert JSON to struct/object

				if newCourse.Title == "" {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply course " +
						" information " + "in valid JSON format"))
					return
				}
				// check if course exists; add only if
				// course does not exist
				if _, ok := courses[params["courseid"]]; !ok {
					courses[params["courseid"]] =
						newCourse
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["courseid"]))
				} else {
					// update course
					courses[params["courseid"]] = newCourse
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Course updated: " +
						params["courseid"]))
				}
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " +
					"course information " +
					"in JSON format"))
			}
		}

	}

}

// note REST API URL : domain/api/versionnumber/service/resource
func main() {

	// clear out all the courses stored in the json when user exits and rerun this REST API server
	courses = map[string]courseInfo{}

	// instantiate courses
	courses = make(map[string]courseInfo)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/courses", allcourses)
	router.HandleFunc("/api/v1/courses/{courseid}", course).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 3300")
	log.Fatal(http.ListenAndServe(":3300", router))

}
