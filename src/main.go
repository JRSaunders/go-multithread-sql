package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

type NodeQueries struct {
	NodeQueries []*NodeQuery `json:"node_queries"`
	Auth        Auth         `json:"auth"`
	Kill        int          `json:"kill" default0:"0"`
}

type Auth struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

type Dsn struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
	CharSet  string `json:"charset"`
}

type DockerNetwork struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Node struct {
	Name          string        `json:"name"`
	Dsn           Dsn           `json:"dsn"`
	Geo           string        `json:"geo"`
	DockerNetwork DockerNetwork `json:"docker_network"`
}

type NodeQuery struct {
	Node            Node         `json:"node"`
	Sql             string       `json:"sql"`
	Binds           []BindString `json:"binds"`
	JsonReturnBytes []interface{}
	Error           string
}
type BindString struct {
	Value string `json:"value"`
	Key   string `json:"key"`
}
type ReturnDataNodes struct {
	Nodes []ReturnData `json:"nodes"`
}

type ReturnData struct {
	NodeName string        `json:"node_name"`
	Data     []interface{} `json:"data"`
	Error    string        `json:"error"`
}

var username string
var password string
var conns int
var connections *int
var debug bool = false
var docker bool = false

func main() {
	dbug := os.Getenv("GOTHREADED_DEBUG")
	if dbug != "" {
		fmt.Println("Debug mode!")
		debug = true
	}
	dcker := os.Getenv("GOTHREADED_DOCKER")
	if dcker != "" {
		fmt.Println("Running in Docker container")
		docker = true
	}
	user := os.Getenv("GOTHREADED_USER")
	if user == "" {
		user = "gothreaded"
	}
	pass := os.Getenv("GOTHREADED_PASS")
	if pass == "" {
		pass = "password"
	}
	host := os.Getenv("GOTHREADED_HOST")

	port := os.Getenv("GOTHREADED_PORT")
	if port == "" {
		port = "1534"
	}
	username = user
	password = pass
	conns = 0
	connections = &conns

	ln, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(ln.Addr().String() + `: Ready to receive connections`)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print(err)
		}

		fmt.Println(conn.RemoteAddr().String() + `: Connected`)
		go handleConnection(conn, connections)
	}
}

func runNodeQuery(db *sql.DB, nodeQuery *NodeQuery, driver string) (*sql.Rows, error) {

	params := make([]interface{}, len(nodeQuery.Binds))

	for i, v := range nodeQuery.Binds {
		params[i] = v.Value
	}
	sql := nodeQuery.Sql
	re := regexp.MustCompile(`(<|and|>|>=|<=|like|between|^|=|in \(|,) :\w+ (,|order|and|or|limit|\))|( |\():\w+(,|\))|( ):\w+($)`)

	sql = re.ReplaceAllString(sql, "$1 ? $2")

	if driver == "postgres" {

		re = regexp.MustCompile(`(<|and|>|>=|<=|like|between|^|=|in \(|,) \? ( |,|order|and|or|limit|\))|( |\()\?(,|\))|( )\?($)`)
		i := 0
		sql = re.ReplaceAllStringFunc(sql, func(s string) string {
			i++
			qmark := fmt.Sprintf("$%d", i)
			s = strings.Replace(s, "?", qmark, -1)
			return s
		})
	}
	if debug {
		fmt.Println("QUERY=")
		f := fmt.Sprintf("%v", params)
		fmt.Println(sql + f)
	}
	return db.Query(sql, params...)

}

func handleConnection(conn net.Conn, connections *int) {
	*connections = *connections + 1
	fmt.Println(*connections, " Total connections")

	var response [4096000]byte
	n, _ := conn.Read(response[0:])
	s := string(response[0:n])

	var nq NodeQueries

	json.Unmarshal([]byte(s), &nq)
	if debug {
		fmt.Println(s)
	}

	data := ReturnDataNodes{}
	if nq.Auth.User != username || nq.Auth.Pass != password {
		fmt.Println(conn.RemoteAddr().String() + " Auth Failed!")
		for _, nodeQuery := range nq.NodeQueries {
			data.Nodes = append(data.Nodes, ReturnData{
				NodeName: nodeQuery.Node.Name,
				Data:     nodeQuery.JsonReturnBytes,
				Error:    "Failed GoThreaded Auth",
			})
		}

	} else {
		if nq.Kill > 0 {
			fmt.Println("Killed!")
			fmt.Println("1")
			os.Exit(1)
		}
		var wg sync.WaitGroup
		for _, nodeQuery := range nq.NodeQueries {
			wg.Add(1)
			go runQuery(&wg, nodeQuery)

		}
		wg.Wait()
		fmt.Println(conn.RemoteAddr().String() + " Done")

		for _, nodeQuery := range nq.NodeQueries {
			data.Nodes = append(data.Nodes, ReturnData{
				NodeName: nodeQuery.Node.Name,
				Data:     nodeQuery.JsonReturnBytes,
				Error:    nodeQuery.Error,
			})
		}

	}
	y, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	conn.Write(y)
	close(conn)
}

func runQuery(wg *sync.WaitGroup, nodeQuery *NodeQuery) {

	finalRows := []interface{}{}
	nodeDsn := nodeQuery.Node.Dsn
	nodeDocker := nodeQuery.Node.DockerNetwork
	driver := ""
	dsn := ""
	host := nodeDsn.Host
	port := nodeDsn.Port

	if docker && nodeDocker.Host != "" {
		host = nodeDocker.Host
	}
	if docker && nodeDocker.Port != "" {
		port = nodeDocker.Port
	}
	switch nodeDsn.Driver {
	case "mysql":
		driver = "mysql"
		dsn = nodeDsn.User + ":" + nodeDsn.Password +
			"@(" + host + ":" + port +
			")/" + nodeDsn.Dbname + "?charset=" + nodeDsn.CharSet
	case "pgsql", "postgres":
		driver = "postgres"
		dsn = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host,
			port,
			nodeDsn.User,
			nodeDsn.Password,
			nodeDsn.Dbname)
	}
	if driver == "" {
		fmt.Println("No Sql Driver")
		nodeQuery.Error = "No Sql Driver"
		wg.Done()
		return
	}

	db, err := sql.Open(driver, dsn)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := runNodeQuery(db, nodeQuery, driver)
	if err != nil {
		fmt.Println(err.Error())
		nodeQuery.Error = err.Error()
		db.Close()
		wg.Done()
		return
	}
	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		fmt.Println(err.Error())
		nodeQuery.Error = err.Error()
		db.Close()
		wg.Done()
		return
	}

	count := len(columnTypes)

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP", "DATETIME", "DATE":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			fmt.Println(err.Error())
			db.Close()
			wg.Done()
			return
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}
	db.Close()
	if debug {
		results := fmt.Sprintf("%d Results", len(finalRows))
		fmt.Println(nodeQuery.Node.Name + ": " + results)
	}
	nodeQuery.JsonReturnBytes = finalRows

	wg.Done()
}

func close(conn net.Conn) {
	conn.Close()
	*connections = *connections - 1
}
