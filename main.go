package main

import (
  "encoding/json"
  "fmt"
  "github.com/julienschmidt/httprouter"
  "database/sql"
  _ "github.com/lib/pq"
  "net/http"
  "log"
)

const (
  DB_USER     = "charlie"
  DB_NAME     = "eventdb"
)

func Index( w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "Welcome!\n")
}

func Hello( w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  log.Println("Hmmmm:")
  log.Println(ps.ByName("name"))
  log.Println(ps)
  fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func TestFetch( w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  type Record struct {
    Id int
    Name string
    Data string
  }
  type RecordSlice struct {
    Records []Record
  }
  db, err := sql.Open("postgres", "user=charlie dbname=eventdb sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }
  rows, err := db.Query("SELECT * FROM test")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()
  var record RecordSlice
  var id int
  var name string
  var data string
  for rows.Next() {
    err := rows.Scan(&id, &name, &data)
    if err != nil {
      log.Fatal(err)
    }
    log.Println(id)
    log.Println(name)
    log.Println(data)
    record.Records = append(record.Records, Record{id, name, data})
  }
  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }
  j, err := json.Marshal(record)
  if err != nil {
    fmt.Println("json err:", err)
  }
  fmt.Println(string(j))
  log.Println(id, name, data)
  fmt.Fprintf(w, string(j))
}

func main() {
  router := httprouter.New()
  router.GET("/", Index)
  router.GET("/hello/:name", Hello)
  router.GET("/fetch", TestFetch)

  log.Fatal(http.ListenAndServe(":8080", router))
}
