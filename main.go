package main

import (
  "time"
  "encoding/json"
  "path"
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

func ImageFetcher( w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  filePath := path.Join("images", "jump_dog.jpg")
  http.ServeFile(w, r, filePath)
}

func FetchPGJSON( w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  start := time.Now()
  db, err := sql.Open("postgres", "user=charlie dbname=eventdb sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }
  rows, err := db.Query("SELECT json_agg(test) FROM test")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()
  result := ""
  var jsonRecord string
  for rows.Next() {
    err := rows.Scan(&jsonRecord)
    if err != nil {
      log.Fatal(err)
    }
    result += jsonRecord
  }
  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Fprintf(w, result)
  elapsed := time.Since(start)
  log.Printf("fetch with PG serialization took %s", elapsed)
}

func FetchGoJSON( w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  start := time.Now()
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
  fmt.Fprintf(w, string(j))
  elapsed := time.Since(start)
  log.Printf("fetch with Go serialization took %s", elapsed)
}

func main() {
  router := httprouter.New()
  router.GET("/", Index)
  router.GET("/hello/:name", Hello)
  router.GET("/fetchJson1", FetchPGJSON)
  router.GET("/fetchJson2", FetchGoJSON)
  router.GET("/image", ImageFetcher)

  log.Fatal(http.ListenAndServe(":3000", router))
}
