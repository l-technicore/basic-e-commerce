package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"log"
)

func AllUsers(w http.ResponseWriter, r *http.Request) {
	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "select id, name, mobile, address from users")
	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Close error is not nil:", err.Error())
		}
	}()
	var data []user
	var success bool
	for rows.Next() {
		success = false
		row := make([]interface{}, 4)
		err = rows.Scan(&row[0], &row[1], &row[2], &row[3])
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}

		var temp user

		if v, flag := row[0].(int); flag {
			temp.id = v;
		} else {
			log.Println("Id is not Int!")
			break
		}
		if v, flag := row[1].(string); flag {
			temp.name = v;
		} else {
			log.Println("Name is not String!")
			break
		}
		if v, flag := row[2].(int); flag {
			temp.mobile = v;
		} else {
			log.Println("Mobile is not Int!")
			break
		}
		if v, flag := row[3].(string); flag {
			temp.address = v;
		} else {
			log.Println("Address is not String!")
			break
		}
		data = append(data, temp)
		success = true
	}
	if err := rows.Err(); !success || err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Fetching Data: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}
	body, _ := json.Marshal(output{code: http.StatusOk, message: "Success", data: data})
	fmt.Fprintf(w, "%s", string(body))
}

func ReadUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "select id, name, mobile, address from users where id = :id", user_id)
	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Close error is not nil:", err.Error())
		}
	}()
	var data user
	var success bool
	for rows.Next() {
		success = false
		row := make([]interface{}, 4)
		err = rows.Scan(&row[0], &row[1], &row[2], &row[3])
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}

		if v, flag := row[0].(int); flag {
			data.id = v;
		} else {
			log.Println("Id is not Int!")
			break
		}
		if v, flag := row[1].(string); flag {
			data.name = v;
		} else {
			log.Println("Name is not String!")
			break
		}
		if v, flag := row[2].(int); flag {
			data.mobile = v;
		} else {
			log.Println("Mobile is not Int!")
			break
		}
		if v, flag := row[3].(string); flag {
			data.address = v;
		} else {
			log.Println("Address is not String!")
			break
		}
		success = true
	}
	if err := rows.Err(); !success || err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Fetching Data: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}
	body, _ := json.Marshal(output{code: http.StatusOk, message: "Success", data: data})
	fmt.Fprintf(w, "%s", string(body))
}

// Requires 3 fields name, mobile, address in JSON payload of the POST request
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	// Can implement sanity checks for input Params (not implemented here)
    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "insert into users (name, mobile, address) values (:name, :mobile, :address) returning id", data["name"], data["mobile"], data["address"])
	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Close error is not nil:", err.Error())
		}
	}()
	var user_id int
	var success bool
	for rows.Next() {
		success = false
		row := make([]interface{}, 1)
		err = rows.Scan(&row[0])
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}

		if v, flag := row[0].(int); flag {
			user_id = v;
		} else {
			log.Println("Id is not Int!")
			break
		}
		success = true
	}
	if err := rows.Err(); !success || err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Fetching Data: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	body, _ := json.Marshal(output{code: http.StatusOk, message: "New User Creation Success!", data: "{\"id\":"+fmt.Sprintf("%d", user_id)+"}"})
	fmt.Fprintf(w, "%s", string(body))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

	var data map[string]interface{}

    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "update users SET (name, mobile, address) = (:name, :mobile, :address) where id = :id", data["name"], data["mobile"], data["address"], user_id)
	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Close error is not nil:", err.Error())
		}
	}()

	body, _ := json.Marshal(output{code: http.StatusOk, message: "User Updation Success!"})
	fmt.Fprintf(w, "%s", string(body))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "delete from users where id = :id", user_id)
	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Close error is not nil:", err.Error())
		}
	}()

	body, _ := json.Marshal(output{code: http.StatusOk, message: "User Delete Success!"})
	fmt.Fprintf(w, "%s", string(body))
}

func AllOrders(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

    var api_response string
	response, err 	:= http.Get("localhost:9002/users/"+user_id+"/orders")

	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Fetching Orders: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	} else {
		defer response.Body.Close()
		data, err := 	ioutil.ReadAll(response.Body)
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Reading Orders: " + err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}
		api_response = 	string(data)
	}

	fmt.Fprintf(w, "%s", api_response)
}

func Order(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]
    order_id := vars["order_id"]

    var api_response string
	response, err 	:= http.Get("localhost:9002/users/"+user_id+"/orders/"+order_id)

	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Fetching Orders: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	} else {
		defer response.Body.Close()
		data, err := 	ioutil.ReadAll(response.Body)
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Reading Orders: " + err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}
		api_response = 	string(data)
	}
	
	fmt.Fprintf(w, "%s", api_response)
}

// Requires 2 fields order_details and order_value in JSON payload of the POST request
func CreateOrder(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	// Can implement sanity checks for input Params (not implemented here)
    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

    var api_response string
	response, err 	:= http.Post("localhost:9002/users/"+user_id+"/orders/create", "application/json", reqBody)

	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Generating Order: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	} else {
		defer response.Body.Close()
		data, err := 	ioutil.ReadAll(response.Body)
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Generating Order: " + err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}
		api_response = 	string(data)
	}

	fmt.Fprintf(w, "%s", api_response)
}

// Requires 2 fields order_details and order_value in JSON payload of the POST request
func UpdateOrder(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	// Can implement sanity checks for input Params (not implemented here)
    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]
    order_id := vars["order_id"]

    var api_response string
	response, err 	:= http.Post("localhost:9002/users/"+user_id+"/orders/"+order_id+"/update", "application/json", reqBody)

	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Updating Order: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	} else {
		defer response.Body.Close()
		data, err := 	ioutil.ReadAll(response.Body)
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Updating Order: " + err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}
		api_response = 	string(data)
	}

	fmt.Fprintf(w, "%s", api_response)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	// Can implement sanity checks for input Params (not implemented here)
    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]
    order_id := vars["order_id"]

    var api_response string
	response, err 	:= http.Post("localhost:9002/users/"+user_id+"/orders/"+order_id+"/delete", "application/json", reqBody)

	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error while Deleting Order: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	} else {
		defer response.Body.Close()
		data, err := 	ioutil.ReadAll(response.Body)
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error while Deleting Order: " + err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}
		api_response = 	string(data)
	}

	fmt.Fprintf(w, "%s", api_response)
}

func DeleteAllOrders(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

    var api_response string
	response, err 	:= http.Get("localhost:9002/users/"+user_id+"/delete/orders")

	if err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Unknown Error while Deleting Order: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	} else {
		defer response.Body.Close()
		data, err := 	ioutil.ReadAll(response.Body)
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Unknown Error while Deleting Orders: " + err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}
		api_response = 	string(data)
	}

	fmt.Fprintf(w, "%s", api_response)
}