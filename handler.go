package main
import (
	"log"
	"net/http"
        "github.com/gorilla/mux"
        "gopkg.in/yaml.v2"
        "io/ioutil"
	"fmt"
        "reflect"
)
type conf struct{
  
  R []Routes `yaml:"Routes"`
}


type Commander struct{}


type Mystruct struct{
 id func(http.ResponseWriter,*http.Request)
}


type Routes struct{
    Path string `yaml:"Path"`
    Callback string `yaml:"Callback"`
    Method string `yaml:"Method"`
}



func handleRequests() {
        var i int
	myRouter := mux.NewRouter().StrictSlash(true)
        rout := getconfig()
        c := &Commander{}
        for i=0 ;i<len(rout.R);i++ {
        fmt.Println(rout.R[i].Path,rout.R[i].Callback,rout.R[i].Method)
        
        m := reflect.ValueOf(c).MethodByName(rout.R[i].Callback)
        Call := m.Interface().(func(http.ResponseWriter,*http.Request))
      if rout.R[i].Method == "GET"{
       myRouter.HandleFunc(rout.R[i].Path, Call ).Methods(rout.R[i].Method)
      }else{
      myRouter.HandleFunc(rout.R[i].Path, Call )
      }
   }


	/*myRouter.HandleFunc("/api/v1/project/putdata", putdata)
        myRouter.HandleFunc("/api/v1/project/getdataByManager/{id}", getdataByManager).Methods("GET")
        myRouter.HandleFunc("/api/v1/project/getdataByProject/{id}", getdataByProject).Methods("GET")
        myRouter.HandleFunc("/api/v1/project/getProjectName", getProjectName).Methods("GET")*/


	log.Fatal(http.ListenAndServe(":8000", myRouter))
}



func  getconfig()(c conf){
    yamlFile, err := ioutil.ReadFile("routes.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal([]byte(yamlFile), &c)
fmt.Println("1")
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}
