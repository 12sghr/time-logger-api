package main

import (
    "fmt"
    //"html/template"
    "log"
    "net/http"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
    "time"
    "strconv"
    "reflect"
    "encoding/hex"
    "golang.org/x/crypto/scrypt"
    "encoding/binary"
    "crypto/rand"
    "os"
)

// type Task struct {
//     UserId int      `json:"userid"`
//     Title string     `json:"title"`
//     Begin int    `json:"begin"`
//     End int      `json:"end"`
// }

type DisplayTask struct {
  Title string        `json:"title"`
  Begin int           `json:"begin"`
  End int              `json:"end"`
  TaskId int          `json:"taskId"`
  LongDays int      `json:"longDays"`
  LongHours int     `json:"longHours"`
  LongMinutes int  `json:"longMinutes"`
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

// func login(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("method:", r.Method) //リクエストを取得するメソッド
//     if r.Method == "GET" {
//         t, _ := template.ParseFiles("login.gtpl")
//         t.Execute(w, nil)
//     } else {
//         //ログインデータがリクエストされ、ログインのロジック判断が実行されます。
//         r.ParseForm()
//         fmt.Println("username:", r.Form["username"])
//         fmt.Println("password:", r.Form["password"])
//     }
// }

func createAccount(w http.ResponseWriter, r*http.Request) {
    fmt.Println("method:", r.Method)

    if r.Method == "GET" {
        fmt.Println("GET createAccount")


    } else {
        fmt.Println("POST createAccount")
        r.ParseForm()

        userName := r.Form["user_name"][0]
        password := r.Form["password"][0]

        hashedPassword, salt := hashPassword(password)

        db, err := sql.Open("mysql", "bbbf544d0f318f:d1d1d9ba@tcp(us-cdbr-iron-east-04.cleardb.net:3306)/heroku_2e3348d8916b23d?interpolateParams=true")
        if err != nil {
            panic(err.Error())
        }

        defer db.Close() // 関数がリターンする直前に呼び出される

        _, insertErr := db.Exec("INSERT INTO `users` (`name`, `hash`, `salt`) VALUES (?, ?, ?);", userName, hashedPassword, salt) //
        if insertErr != nil {
            panic(insertErr.Error())
        }


    }
}

func login(w http.ResponseWriter, r*http.Request) {
    fmt.Println("method:", r.Method)

    if r.Method == "GET" {
        fmt.Println("GET login")


    } else {
        fmt.Println("POST login")
        r.ParseForm()

        //userName := r.Form["user_name"][0]


    }
}

func hashPassword(pass string) (string, string) {
    salt := random()
    fmt.Println(salt)
    byteSalt := []byte(salt)
    converted, _ := scrypt.Key([]byte(pass), byteSalt, 16384, 8, 1, 32)
    return hex.EncodeToString(converted[:]), salt
}

func random() string {
    var n uint64
    binary.Read(rand.Reader, binary.LittleEndian, &n)
    return strconv.FormatUint(n, 36)
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
        db, err := sql.Open("mysql", "bbbf544d0f318f:d1d1d9ba@tcp(us-cdbr-iron-east-04.cleardb.net:3306)/heroku_2e3348d8916b23d?interpolateParams=true")
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

        // task := {
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

        db, err := sql.Open("mysql", "bbbf544d0f318f:d1d1d9ba@tcp(us-cdbr-iron-east-04.cleardb.net:3306)/heroku_2e3348d8916b23d?interpolateParams=true")
        if err != nil {
            panic(err.Error())
        }

        defer db.Close() // 関数がリターンする直前に呼び出される

        userId := 3

        rows, qerr := db.Query("SELECT title, begin, end, task_id, `longDays`, `longHours`, `longMinutes` FROM tasks WHERE user_id = ? ORDER BY task_id DESC", userId)
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
            var longDays int
            var longHours int
            var longMinutes int
            if berr := rows.Scan(&title, &begin, &end, &taskId, &longDays, &longHours, &longMinutes); berr != nil {
                log.Fatal("scan error: %v", berr)
            }
            displayTask := DisplayTask{
                Title: title,
                Begin: begin,
                End: end,
                TaskId: taskId,
                LongDays: longDays,
                LongHours: longHours,
                LongMinutes: longMinutes,
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
    db, err := sql.Open("mysql", "bbbf544d0f318f:d1d1d9ba@tcp(us-cdbr-iron-east-04.cleardb.net:3306)/heroku_2e3348d8916b23d?interpolateParams=true")
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
    hours0 := int(duration.Hours())
    longDays := hours0 / 24
    longHours := hours0 % 24
    longMinutes := int(duration.Minutes()) % 60


    _, updateErr := db.Exec("UPDATE tasks SET `end` = ?, `longDays` = ?, `longHours` = ?, `longMinutes` = ? WHERE task_id = ?", endTime, longDays, longHours, longMinutes, taskId) //
    //_, updateErr := db.Exec("UPDATE tasks SET end = ? WHERE task_id = ?", endTime, taskId) //
    if updateErr != nil {
        panic(updateErr.Error())
    }

    return

}

func main() {
    fmt.Println("Server starting......")
    http.HandleFunc("/", sayhelloName)       //アクセスのルーティングを設定します
    //http.HandleFunc("/login", login)         //アクセスのルーティングを設定します
    http.HandleFunc("/createAccount", createAccount)
    http.HandleFunc("/login", login)
    http.HandleFunc("/mainPage", mainPage)
    http.HandleFunc("/mainPage/end", mainEnd)
    listen := os.Getenv("PORT")
    if listen == "" {
        listen = "9090"
    }
    fmt.Println("Listening on " + listen + " port")
    err := http.ListenAndServe(":" + listen, nil) //監視するポートを設定します
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
