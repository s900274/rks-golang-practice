package define

type ServiceConfig struct {
    LogFile          string     `flag:"logfile" cfg:"logfile" toml:"logfile"`
    HttpServerIp     string     `cfg:"httpserverport" toml:"httpserverport"`
    HttpServerPort   int        `cfg:"httpserverip" toml:"httpserverip"`
}