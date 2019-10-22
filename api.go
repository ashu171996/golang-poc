package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
        "strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)
type update struct{
     Pm string  `json:"manager_name"`
     Pn string  `json:"name"`
}

type Project struct {
        Id             string `json:"Id"`
	ProjectName    string `json:"name"`
	ManagerName    string `json:"manager_name"`
	ManagerEmailID string `json:"manager_email_id"`
	Flag           string `json:"flag"`
        
}


type Recive struct {
       // Id             string `json:"Id"`
	ProjectName    string `json:"name"`
	ManagerName    string `json:"manager_name"`
	ManagerEmailID string `json:"manager_email_id"`
	//Flag           string `json:"flag"`
}

type Error struct {
        Errtype string `json:"error type"`
        Message string `json:"error message"`
}


func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}


func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root987"
	dbName := "weekly_update"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(172.21.234.63:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}




func (c *Commander) Putdata(w http.ResponseWriter, r *http.Request) {
       var erro Error
       db := dbConn()

	reqToken := r.Header.Get("Authorization")

	splitToken := strings.Split(reqToken, "Bearer")
        reqToken = splitToken[1]
        user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
	if user.Next() != false {
		var dat Recive

		json.NewDecoder(r.Body).Decode(&dat)
		Pn := dat.ProjectName
		Mn := dat.ManagerName
		Email := dat.ManagerEmailID
		Flag := 1

		insForm, err := db.Prepare("INSERT INTO Project(name, manager_name,manager_email_id,flag)VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		result,_ := insForm.Exec(Pn, Mn, Email,Flag)
                _,er := result.RowsAffected()
                if er != nil{
                 erro.Errtype = "Dupelicate insertion"
                 erro.Message = "Error in insertion"
                 json.NewEncoder(w).Encode(erro)
                    }
		defer db.Close()
                setupResponse(&w, r)
                w.WriteHeader(http.StatusCreated)
	}else{
       w.WriteHeader(http.StatusUnauthorized)
       }


}



func (c *Commander) GetdataByManager(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	reqToken := r.Header.Get("Authorization")
        offset,_ :=strconv.ParseInt( r.Header.Get("OFFSET"),10,64)

	splitToken := strings.Split(reqToken, "Bearer")
        reqToken = splitToken[1]
        user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
	if user.Next() != false {
		p := mux.Vars(r)
	        key := p["id"]
        
         rows, err := db.Query("SELECT * FROM Project WHERE manager_name=? AND flag = 1 LIMIT 10 OFFSET ?", key,offset)
      	if err != nil {
               fmt.Println("error")
      		log.Fatal(err)
      	}
      	defer rows.Close()
        var pro Project
        var Proj []Project
     
        for rows.Next() {

      		 rows.Scan(&pro.ProjectName,&pro.ManagerName,&pro.ManagerEmailID,&pro.Flag,&pro.Id)
         Proj = append(Proj,pro)
      	}
       setupResponse(&w, r)
       w.Header().Set("Content-Type","application/json")
       w.WriteHeader(http.StatusOK)
       json.NewEncoder(w).Encode(Proj)

	}else{
       w.WriteHeader(http.StatusUnauthorized)
       }
	

}

func (c *Commander) GetdataByProject(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	reqToken := r.Header.Get("Authorization")
        offset,_ :=strconv.ParseInt( r.Header.Get("OFFSET"),10,64)
	splitToken := strings.Split(reqToken, "Bearer")
        reqToken = splitToken[1]
        user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
	if user.Next() != false {
		p := mux.Vars(r)
	        key := p["id"]
        
         rows, err := db.Query("SELECT * FROM Project WHERE name=? AND flag = 1 LIMIT 10 OFFSET ?", key,offset)
      	if err != nil {
               fmt.Println("error")
      		log.Fatal(err)
      	}
      	defer rows.Close()
        var pro Project
        var Proj []Project
     
        for rows.Next() {

      		 rows.Scan(&pro.ProjectName,&pro.ManagerName,&pro.ManagerEmailID,&pro.Flag,&pro.Id)
         Proj = append(Proj,pro)
      	}
       setupResponse(&w, r)
       w.Header().Set("Content-Type","application/json")
       w.WriteHeader(http.StatusOK)
       json.NewEncoder(w).Encode(Proj)

	}else{
       w.WriteHeader(http.StatusUnauthorized)
       }
}


func (c *Commander) GetProjectName(w http.ResponseWriter,r *http.Request) {

      db := dbConn()
      reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
        reqToken = splitToken[1]
        user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
	if user.Next() != false {
        rows, err := db.Query("SELECT DISTINCT name FROM Project")
      	if err != nil {
               fmt.Println("error")
      		log.Fatal(err)
      	}
      	defer rows.Close()
     var Nam []string
     var Name string
        
         for rows.Next() {

      		 rows.Scan(&Name)
         Nam = append(Nam,Name)
      	}
       setupResponse(&w, r)
       w.Header().Set("Content-Type","application/json")
       w.WriteHeader(http.StatusOK)
       json.NewEncoder(w).Encode(Nam)

}else{
       w.WriteHeader(http.StatusUnauthorized)
       }

}


func (c *Commander) UpdateData(w http.ResponseWriter,r *http.Request) {
     db := dbConn()
      reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
        reqToken = splitToken[1]
        user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
	if user.Next() != false {
         var dat Project

		json.NewDecoder(r.Body).Decode(&dat)
                ID := dat.Id
		Pn := dat.ProjectName
		Mn := dat.ManagerName
		Email := dat.ManagerEmailID
		//Flag := dat.Flag

         update,_ := db.Query("UPDATE Project name = ?, manager_name = ?,manager_email_id = ? WHERE Id =?",Pn,Mn,Email,ID)
         defer update.Close() 

        }     



}


func (c *Commander) DeleteData(w http.ResponseWriter,r *http.Request) {
       db := dbConn()
       var dat update
        reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
        reqToken = splitToken[1]
        user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
	if user.Next() != false {
         
          json.NewDecoder(r.Body).Decode(&dat)

         del,_ := db.Query("UPDATE Project flag = 0 WHERE  Project name = ?, manager_name = ?",dat.Pn,dat.Pm)
         defer del.Close()



}


}
