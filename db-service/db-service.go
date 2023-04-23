package dbsvc

type DBService interface{
	SetErrorResp(status int, msg string)
	SetResult(result map[string]interface{}, rows ...<-chan interface{})
	ReadJSON(params interface{}) (status int, err error)
}

type FnAction func(DBService)
var actions = map[string]FnAction {
	"open-db":    OpenDB,
	"close-db":   CloseDB,
	"begin-tx":   BeginTx,
	"commit":     Commit,
	"rollback":   Rollback,
	"prepare":    Prepare,
	"close-stmt": CloseStmt,
	"exec":       Exec,
	"query":      Query,
	"ping":       Ping,
}

func GetAction(name string) (action FnAction, ok bool) {
	action, ok = actions[name]
	return
}

func GetActionNames() <-chan string {
	res := make(chan string)
	go func() {
		for name, _ := range actions {
			res <- name
		}
		close(res)
	}()
	return res
}
