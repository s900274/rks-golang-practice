package define

type ServiceConfig struct {
    LogFile          string     `cfg:"logfile" toml:"logfile"`
    HttpServerIp     string     `cfg:"httpserverip" toml:"httpserverip"`
    HttpServerPort   int        `cfg:"httpserverport" toml:"httpserverport"`
}