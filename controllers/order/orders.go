package order

import (
	"encoding/json"
	"fmt"
	"log"
)

func AllOrders(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
    user_id := vars["id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// rows, err = db.QueryContext(ctx, "select o.order_id as order_id, p.product_name as product, o.quantity as quantity, o.value as value from orders o, products p where o.fk_product_id = p. product_id and fk_user_id = :user_id", user_id)
	rows, err = db.QueryContext(ctx, "select order_id, details, value from orders where fk_user_id = :user_id", user_id)
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
	var data []order
	var success bool
	for rows.Next() {
		success = false
		row := make([]interface{}, 3)
		err = rows.Scan(&row[0], &row[1], &row[2])
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}

		var temp order

		if v, flag := row[0].(int); flag {
			temp.id = v;
		} else {
			log.Println("Order Id is not Int!")
			break
		}
		if v, flag := row[1].(string); flag {
			temp.details = v;
		} else {
			log.Println("Product name is not String!")
			break
		}
		if v, flag := row[2].(int); flag {
			temp.value = v;
		} else {
			log.Println("Order value is not Int!")
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

func Order(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
    user_id := vars["id"]
    order_id := vars["order_id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "select order_id, details, value from orders where fk_user_id = :user_id and order_id = :order_id", user_id, order_id)
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
	var data order
	var success bool
	for rows.Next() {
		success = false
		row := make([]interface{}, 3)
		err = rows.Scan(&row[0], &row[1], &row[2])
		if err != nil {
			body, _ := json.Marshal(output{code: http.StatusBadRequest, message: err.Error()})
			fmt.Fprintf(w, "%s", string(body))
			return
		}

		if v, flag := row[0].(int); flag {
			data.id = v;
		} else {
			log.Println("Order Id is not Int!")
			break
		}
		if v, flag := row[1].(string); flag {
			data.details = v;
		} else {
			log.Println("Product name is not String!")
			break
		}
		if v, flag := row[2].(int); flag {
			data.value = v;
		} else {
			log.Println("Order value is not Int!")
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

// Requires 2 fields order_details and order_value in JSON payload of the POST request
func CreateOrder(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	// Can implement sanity checks for input Params (not implemented here)
    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "insert into orders (fk_user_id, details, value) values (:user_id, :order_details, :order_value) returning order_id", user_id, data["order_details"], data["order_value"])
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
	var order_id int
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
			order_id = v;
		} else {
			log.Println("Order Id is not Int!")
			break
		}
		success = true
	}
	if err := rows.Err(); !success || err != nil {
		body, _ := json.Marshal(output{code: http.StatusFailedDependency, message: "Error Fetching Data: " + err.Error()})
		fmt.Fprintf(w, "%s", string(body))
		return
	}

	body, _ := json.Marshal(output{code: http.StatusOk, message: "New Order Creation Success!", data: "{\"order_id\":"+fmt.Sprintf("%d", order_id)+"}"})
	fmt.Fprintf(w, "%s", string(body))
}

// Requires 2 fields order_details and order_value in JSON payload of the POST request
func UpdateOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]
    order_id := vars["order_id"]

	var data map[string]interface{}

    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &data)

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "update orders SET (details, value) = (:order_details, :value) where fk_user_id = :id and order_id = :order_id", data["order_details"], data["order_value"], user_id, order_id)
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

	body, _ := json.Marshal(output{code: http.StatusOk, message: "Order Updation Success!"})
	fmt.Fprintf(w, "%s", string(body))
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]
    order_id := vars["order_id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "delete from orders where fk_user_id = :id and order_id = :order_id", user_id, order_id)
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

	body, _ := json.Marshal(output{code: http.StatusOk, message: "Order Deletion Success!"})
	fmt.Fprintf(w, "%s", string(body))
}

func DeleteAllOrders(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Can implement sanity checks for input Params (not implemented here)
    user_id := vars["id"]

	var rows *sql.Rows
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	rows, err = db.QueryContext(ctx, "delete from orders where fk_user_id = :id", user_id)
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

	body, _ := json.Marshal(output{code: http.StatusOk, message: "All User Orders Deleted Successfully!"})
	fmt.Fprintf(w, "%s", string(body))
}