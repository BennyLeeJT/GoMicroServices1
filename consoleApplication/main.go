package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type itemInfo struct {
	category int
	quantity int
	unitCost float64
}

type courseInfo struct {
	Title string `json:"Title"`
	Id    string `json:"CourseResourceId"`
}

var jsonData map[string]string

var itemName map[string]itemInfo

const baseURL = "http://localhost:3300/api/v1/courses"

const Key = "2c78afaf-97da-4816-bbee-9ad239abb296"

func getCourse(courseID string) {
	url := baseURL

	if courseID != "" {
		url = baseURL + "/" + courseID + "?key=" + Key
	}

	if courseID == "" {
		url = baseURL + "?key=" + Key
	}

	response, err := http.Get(url) // this is the GET request

	if err != nil { // if there IS error
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else { // if NO error

		data, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("response.StatusCode : %v\n", response.StatusCode)
		fmt.Printf("%v\n\n", string(data))

		response.Body.Close()
	}
}

// jsonData is a map initially. not yet JSON
// Key is the access token
func addCourse(courseID string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData) // convert object/struct/map into JSON

	// fmt.Printf("JSON : %v\n", string(jsonValue))

	// this is the POST request
	response, err := http.Post(baseURL+"/"+courseID+"?key="+Key,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil { // if there IS error
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else { // if NO error
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("response.StatusCode : %v\n", response.StatusCode)
		fmt.Printf("%v\n\n", string(data))

		response.Body.Close()
	}

	getCourse("")
}

func updateCourse(courseID string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)

	request, err := http.NewRequest(http.MethodPut,
		baseURL+"/"+courseID+"?key="+Key,
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("err : %v\n\n", err)
	}

	request.Header.Set("Content-Type", "application/json") // you can separately set the header's content type as json

	client := &http.Client{} // note http.Client has high level code for GET & POST, but not PUT or DELETE

	response, err := client.Do(request)

	if err != nil { // if there IS error
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else { // if NO error
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("response.StatusCode : %v\n\n", response.StatusCode)
		fmt.Printf("%v\n\n", string(data))

		response.Body.Close()
	}

	getCourse("")
}

func deleteCourse(courseID string) {
	request, err := http.NewRequest(http.MethodDelete,
		baseURL+"/"+courseID+"?key="+Key, nil)

	if err != nil { // if there IS error
		fmt.Printf("err : %v\n", err)
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil { // if there IS error
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else { // if NO error
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Printf("response.StatusCode : %v\n\n", response.StatusCode)
		fmt.Printf("%v\n\n", string(data))

		response.Body.Close()
	}

	getCourse("")
}

func deleteMySQL(db *sql.DB) {
	query := fmt.Sprintln("DELETE FROM courses;")
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// reset so that if nothing is entered, it will not execute based on last entered choice
func resetChoiceMenu(lastEnteredChoice *int) {
	*lastEnteredChoice = 0
}

func resetAutoIncrementIDColumn(db *sql.DB) {
	query := fmt.Sprintln("ALTER TABLE courses AUTO_INCREMENT = 0;")
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func addToMySQL(db *sql.DB, newCourseResourceIdToAddIntoSQL string, newCourseTitleToAddIntoSQL string) {
	query := fmt.Sprintf("INSERT INTO courses (CourseID, CourseTitle) VALUES ('%s', '%s')",
		newCourseResourceIdToAddIntoSQL, newCourseTitleToAddIntoSQL)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func updateMySQL(db *sql.DB, existingCourseResourceIdToModifyInMySQL string, newCourseTitleToModifyInSQL string) {
	query := fmt.Sprintf("UPDATE courses SET CourseTitle='%s' WHERE CourseID='%s';", newCourseTitleToModifyInSQL, existingCourseResourceIdToModifyInMySQL)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func deleteInMySQL(db *sql.DB, existingCourseResourceIdToDeleteInMySQL string) {
	query := fmt.Sprintf("DELETE FROM courses WHERE CourseID='%s';", existingCourseResourceIdToDeleteInMySQL)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	passcode := godotenv.Load()
	if passcode != nil {
		fmt.Printf("passcode : %v\n\n", passcode)
	}
	extractedDataFromEnv := os.Getenv("ACCESS_TOKEN")
	fmt.Printf("Extracted Key : %v\n\n", extractedDataFromEnv)

	key := []byte(extractedDataFromEnv)
	data := []byte("This is Data For aesgcm.Seal decrypted")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	encryptedDataCypherText := aesgcm.Seal(nil, nonce, data, nil)
	fmt.Printf("Encrypted: %x\n\n", encryptedDataCypherText)
	decryptedData, err := aesgcm.Open(nil, nonce, encryptedDataCypherText, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Decrypted: %s\n\n", decryptedData)

	// ####################
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/my_db")
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	fmt.Println("Console App Database connected")
	// ####################

	// when user exit and rerun this console app, to clear all records in MySQL database
	// this will make testing and checking easier
	deleteMySQL(db)

	// previously ID field in MySQL is set to be auto increment, to reset each time user exit and rerun this console app. KIV for future use.
	resetAutoIncrementIDColumn(db)

	var (
		choiceMenu int

		loopMenu  bool = false
		selection      = [5]int{1: 5}

		newCourseTitleToAdd      string
		newCourseResourceIdToAdd string

		existingCourseResourceIdToModify string
		existingCourseResourceIdToDelete string
	)

	fmt.Println("")

	jsonData = map[string]string{}

	allCoursesInConsole := map[string]courseInfo{}

	for {
		if choiceMenu == 0 || loopMenu == true {

			fmt.Println("Courses Console Application\n" +
				"============================\n" +
				"1. View All Current Courses.\n" +
				"2. Add Course.\n" +
				"3. Modify Course.\n" +
				"4. Delete Course.\n" +
				"5. Delete All Data in MySQL Database.\n" +
				"Select your choice by indicating the number :")

			fmt.Scanln(&choiceMenu)

			if int(choiceMenu) == 0 {
				fmt.Println("** Nothing has been entered for menu choice **") // check
			} else {
				fmt.Printf("** You have entered %v as your menu choice **\n", choiceMenu) // check
			}

			// validate choice of user, to go back to courses list menu if invalid
			if choiceMenu > len(selection) || choiceMenu < 1 {
				fmt.Println("**>> Note : Please select option 1 to 5 only <<**")
				fmt.Println("")
			} else {

				if choiceMenu == 1 { // List all courses

					// GET request sent via http.Get(url) in the function
					getCourse("")

					resetChoiceMenu(&choiceMenu)

					// END OF CHOICE MENU 1 LIST ALL COURSES
				} else if choiceMenu == 2 { // Add Course
					fmt.Println("Add New Course")
					fmt.Println("What is the title of your course? eg. Computer Science 101")
					// fmt.Scanln(&newCourseTitleToAdd) // for one word without white spaces
					// fmt.Scanf("%s %s %s", &a, &b, &c) // for multiple single words, also no white spaces / sentences
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan() // use `for scanner.Scan()` to keep reading
					newCourseTitleToAdd := scanner.Text()
					// fmt.Println("captured:", newCourseTitleToAdd) // will show sentences with white spaces

					fmt.Println("What is the ID of your course? eg. CS101")
					fmt.Scanln(&newCourseResourceIdToAdd)

					jsonData = map[string]string{}

					// only store ONE data to add
					jsonData["Title"] = newCourseTitleToAdd
					jsonData["Id"] = newCourseResourceIdToAdd

					allCoursesInConsole[newCourseResourceIdToAdd] = courseInfo{newCourseTitleToAdd, newCourseResourceIdToAdd}

					fmt.Printf("Adding New Course... please wait\n\n")
					addCourse(newCourseResourceIdToAdd, jsonData)
					addToMySQL(db, newCourseResourceIdToAdd, newCourseTitleToAdd)
					fmt.Printf("Addition done successfully!\n\n")

					resetChoiceMenu(&choiceMenu)
					// END OF CHOICE MENU 2 ADD NEW COURSE
				} else if choiceMenu == 3 { // Modify Course
					fmt.Println("Modify Course.")
					getCourse("")

					fmt.Println("Which course would you wish to modify? Please specify course ID. eg. CS101")
					fmt.Scanln(&existingCourseResourceIdToModify)
					idExists := allCoursesInConsole[existingCourseResourceIdToModify].Id == existingCourseResourceIdToModify

					fmt.Println("")

					if idExists == true {
						fmt.Printf("Course ID to edit is %v\n\n", existingCourseResourceIdToModify)
						fmt.Printf("Enter new Course Title. Eg. Macro Economics 203\n\n")
						// fmt.Scanln(&newCourseTitleToAdd)
						scanner := bufio.NewScanner(os.Stdin)
						scanner.Scan() // use `for scanner.Scan()` to keep reading
						var newCourseTitleToModifyTo string
						newCourseTitleToAdd = scanner.Text()
						// fmt.Printf("newCourseTitleToAdd : %v/n/n", newCourseTitleToAdd)
						updateMySQL(db, existingCourseResourceIdToModify, newCourseTitleToModifyTo)
						fmt.Printf("Course Title for ID %v has been updated in MySQL database successfully\n\n", existingCourseResourceIdToModify)
					} else {
						// user input Id not existing, to create new
						newCourseResourceIdToAdd = existingCourseResourceIdToModify
						fmt.Println("No such Course ID. This Course ID will be created.")
						fmt.Printf("Enter new Course Title for this ID\n\n")
						// fmt.Scanln(&newCourseTitleToAdd)
						scanner := bufio.NewScanner(os.Stdin)
						scanner.Scan() // use `for scanner.Scan()` to keep reading
						var newCourseTitleToAdd string
						newCourseTitleToAdd = scanner.Text()
						addToMySQL(db, newCourseResourceIdToAdd, newCourseTitleToAdd)
					}

					jsonData = map[string]string{}
					jsonData["Title"] = newCourseTitleToAdd
					jsonData["Id"] = newCourseResourceIdToAdd
					updateCourse(existingCourseResourceIdToModify, jsonData)

					resetChoiceMenu(&choiceMenu)
					// END OF CHOICE MENU 3 MODIFY COURSE
				} else if choiceMenu == 4 { // Delete Course
					fmt.Println("Delete Course")

					getCourse("")

					fmt.Println("Enter the Course Id to be deleted")
					existingCourseResourceIdToDelete = "" // to reset if run again
					fmt.Scanln(&existingCourseResourceIdToDelete)

					if len(existingCourseResourceIdToDelete) == 0 || // if empty input (empty string) OR input not found
						allCoursesInConsole[existingCourseResourceIdToDelete].Id == existingCourseResourceIdToDelete == false {

						fmt.Println("Course ID not found. Nothing is deleted. Please enter valid Course ID.")
					} else { // input is found
						fmt.Printf("Course ID %v found in database. Proceeding to delete.\n\n", existingCourseResourceIdToDelete)
						deleteCourse(existingCourseResourceIdToDelete)
						deleteInMySQL(db, existingCourseResourceIdToDelete)
						fmt.Printf("Course ID %v has been deleted from MySQL database successfully.\n\n", existingCourseResourceIdToDelete)
					}

					fmt.Println("")

					resetChoiceMenu(&choiceMenu)
					// END OF CHOICE MENU 4 DELETE COURSE
				} else if choiceMenu == 5 { // Delete all rows in MySQL
					fmt.Println("Delete All Rows in MySQL")

					getCourse("")

					fmt.Println("Key in the word 'delete' or 'Delete' without quotes to confirm to proceed deletion immediately :")
					var confirmToDeleteMySQL string
					confirmToDeleteMySQL = ""
					fmt.Scanln(&confirmToDeleteMySQL)

					if confirmToDeleteMySQL == "delete" || confirmToDeleteMySQL == "Delete" {
						fmt.Printf("You have keyed in %v\n\n", confirmToDeleteMySQL)
						deleteMySQL(db)
						fmt.Printf("Deletion in now in progress and will be done soon.\n\n")
						resetAutoIncrementIDColumn(db)
					} else {
						fmt.Printf("Word keyed is not exact. Please try again.\n\n")
					}

					resetChoiceMenu(&choiceMenu)

				} // END OF CHOICE MENU 5 DELETE ALL ROWS IN MYSQL
			}

		}

		loopMenu = true
	}
}
