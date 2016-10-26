package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
    "time"
    "strconv"
    "reflect"
)

type Task struct {
    UserId int      `json:"userid"`
    Title string     `json:"title"`
    Begin int    `json:"begin"`
    End int      `json:"end"`
}

type DisplayTask struct {
  Title string    `json:"title"`
  Begin int       `json:"begin"`
  End int          `json:"end"`
  TaskId int      `json:"taskId"`
  Long int        `json:"long"`
}

type DisplayTasks []DisplayTask

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()       //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
    //注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。
    fmt.Println(r.Form) //これらのデータはサーバのプリント情報に出力されます
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") //ここでwに書き込まれたものがクライアントに出力されます。
}

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //リクエストを取得するメソッド
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, nil)
    } else {
        //ログインデータがリクエストされ、ログインのロジック判断が実行されます。
        r.ParseForm()
        fmt.Println("username:", r.Form["username"])
        fmt.Println("password:", r.Form["password"])
    }
}

func mainPage(w http.ResponseWriter, r*http.Request) {
    fmt.Println("method:", r.Method)

    var res string

    if r.Method == "POST" {
        fmt.Println("POST mainPage")

        // templates.WritePageTemplate(w)
        r.ParseForm()
        fmt.Println("入力された値: ", r.Form["doing_thing"])
        res = r.Form["doing_thing"][0]
        fmt.Println(res)
        db, err := sql.Open("mysql", "root@tcp(localhost:3306)/time_logger?interpolateParams=true")
        if err != nil {
            panic(err.Error())
        }

        defer db.Close() // 関数がリターンする直前に呼び出される

        i := 3
        t := time.Now()
        beginTime := t.Format("200601021504")
        //e := 201610302300
        _, insertErr := db.Exec("INSERT INTO tasks (user_id, title, begin) VALUES (?, ?, ?);", i, res, beginTime) //
        if insertErr != nil {
            panic(insertErr.Error())
        }

        // task := Task{
        //     UserId: 1,
        //     Title: res,
        //     Begin: 201610102345,
        //     End: 201610110010,
        // }
        // fmt.Println(task)
        //
        // jsonBytes, err := json.Marshal(task)
        // if err != nil {
        //     fmt.Println("JSON Marshal error:", err)
        //     return
        // }
        //
        // w.Header().Set("Content-Type", "application/json")
        // fmt.Fprint(w, string(jsonBytes))

        //json.NewEncoder(w).Encode(task)
    } else {
        fmt.Println("GET mainPage")

        db, err := sql.Open("mysql", "root@tcp(localhost:3306)/time_logger?interpolateParams=true")
        if err != nil {
            panic(err.Error())
        }

        defer db.Close() // 関数がリターンする直前に呼び出される

        userId := 3

        rows, qerr := db.Query("SELECT title, begin, end, task_id, `long` FROM tasks WHERE user_id = ? ORDER BY task_id DESC", userId)
        // db.Query("SELECT * FROM user ",nil)はダメだった。
        defer rows.Close()
        if qerr != nil {
            log.Fatal("query error: %v", qerr)
        }

        var displayTasks DisplayTasks

        for rows.Next() {
            var title string
            var begin int
            var end int
            var taskId int
            var long int
            if berr := rows.Scan(&title, &begin, &end, &taskId, &long); berr != nil {
                log.Fatal("scan error: %v", berr)
            }
            displayTask := DisplayTask{
                Title: title,
                Begin: begin,
                End: end,
                TaskId: taskId,
                Long: long,
            }
            displayTasks = append(displayTasks, displayTask)


        }
        fmt.Println(displayTasks)
        jsonBytes, err := json.Marshal(displayTasks)
        if err != nil {
            fmt.Println("JSON Marshal error:", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, string(jsonBytes))
    }
    fmt.Println("きた")
}

func mainEnd(w http.ResponseWriter, r*http.Request) {
    fmt.Println("method:", r.Method)
    fmt.Println("POST mainEnd")
    db, err := sql.Open("mysql", "root@tcp(localhost:3306)/time_logger?interpolateParams=true")
    if err != nil {
        panic(err.Error())
    }

    defer db.Close() // 関数がリターンする直前に呼び出される

    r.ParseForm()
    var taskIdStr = r.Form["task_id"][0]
    taskId, _ := strconv.Atoi(taskIdStr)


    t := time.Now()
    endTime := t.Format("200601021504")
    fmt.Println(reflect.TypeOf(endTime))

    rows, qerr := db.Query("SELECT begin FROM tasks WHERE task_id = ?", taskId)
    // db.Query("SELECT * FROM user ",nil)はダメだった。
    if qerr != nil {
        log.Fatal("query error: %v", qerr)
    }

    var begin string
    for rows.Next() {
        if berr := rows.Scan(&begin); berr != nil {
            log.Fatal("scan error: %v", berr)
        }
    }

    fmt.Println(begin)

    var timeformat = "200601021504"
    beginT, _ := time.Parse(timeformat, begin)
    endT, _ := time.Parse(timeformat, endTime)
    fmt.Println(beginT)
    fmt.Println(endT)
    fmt.Println(endT.Sub(beginT))
    duration := endT.Sub(beginT)
    long := int(duration.Minutes()) % 60
    fmt.Println(reflect.TypeOf(long))


    _, updateErr := db.Exec("UPDATE tasks SET `end` = ?, `long` = ? WHERE task_id = ?", endTime, long, taskId) //
    //_, updateErr := db.Exec("UPDATE tasks SET end = ? WHERE task_id = ?", endTime, taskId) //
    if updateErr != nil {
        panic(updateErr.Error())
    }

    return

}

func main() {
    http.HandleFunc("/", sayhelloName)       //アクセスのルーティングを設定します
    http.HandleFunc("/login", login)         //アクセスのルーティングを設定します
    http.HandleFunc("/mainPage", mainPage)
    http.HandleFunc("/mainPage/end", mainEnd)
    err := http.ListenAndServe(":9090", nil) //監視するポートを設定します
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
