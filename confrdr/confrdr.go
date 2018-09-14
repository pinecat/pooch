/* package */
package confrdr

/* imports */
import (
    "log"
    "os"
    "net"
    "bufio"
    "strings"
)

/* gloabals */
var (
    PoochConf   PoochConfData
    DBConf      DBConfData
)

/* structs */
type PoochConfData struct {
    Port    string
    NetIf   string
    IP      string
    SSLCert string
    SSLKey  string
}

type DBConfData struct {
    Url     string
    Port    string
    User    string
    Pass    string
}

func ReadConfFile(filepath string) {
    file, err := os.Open(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    c := make(chan string, 20)
    defer close(c)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if (scanner.Text() != "[pooch]" && scanner.Text() != "[db]" && scanner.Text() != "" && scanner.Text() != "\n") {
            s := strings.Split(scanner.Text(), ":")
            c <- s[1]
        }
    }

    PoochConf.Port = <-c;
    PoochConf.NetIf = <-c;
    PoochConf.IP = getIP(PoochConf.NetIf)
    PoochConf.SSLCert = <-c;
    PoochConf.SSLKey = <-c;

    DBConf.Url = <-c;
    DBConf.Port = <-c;
    DBConf.User = <-c;
    DBConf.Pass = <-c;

    getIP(PoochConf.NetIf)
}

func getIP(netif string) (string) {
    ifaces, err := net.Interfaces()
    if err != nil {
        log.Fatal(err)
        return ""
    }

    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            log.Fatal(err)
            return ""
        }

        for _, addr := range addrs {
            if iface.Name == netif {
                s := strings.Split(addr.String(), "/")
                return s[0]
            }
        }
    }
    return ""
}
