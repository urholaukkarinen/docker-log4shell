# docker-log4shell

Simple Go app / Docker image for playing with the [CVE-2021-44228](https://www.cve.org/CVERecord?id=CVE-2021-44228) vulnerability. Hosts a simple file server and an ldap server that provides any classes hosted in the file server: `ldap://ip:1389/Hello -> http://ip:8080/Hello.class`

```
docker build -t log4shell .
docker run -it -v "path-to-classes:/files" -p 1389:1389 -p 8080:8080 -e HOST=localhost log4shell
```