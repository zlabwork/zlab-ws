package zlabws

var Cfg = App{}

type App struct {
    Db struct {
        Mysql struct {
            Host string
            Port string
            User string
            Pass string
        }
        Redis struct {
            Host string
            Port string
        }
    }
}
