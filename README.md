# VianuEdu-Server

[![Build Status](https://travis-ci.org/CNITV/VianuEdu-Server.svg?branch=master)](https://travis-ci.org/CNITV/VianuEdu-Server)
[![Go Report Card](https://goreportcard.com/badge/github.com/CNITV/VianuEdu-Server)](https://goreportcard.com/report/github.com/CNITV/VianuEdu-Server)
[![GoDoc](https://godoc.org/github.com/CNITV/VianuEdu-Server?status.svg)](https://godoc.org/github.com/CNITV/VianuEdu-Server)

Componenta server-side a VianuEdu.

## Cum functioneaza?

In esenta, VianuEdu-Server este un REST API scris in Golang care
comunica cu o baza de date MongoDB cu o schema predefinita. Acest API
incearca sa fie usor de folosit si, in acelasi timp, sigur.

Momentan, nu este complet, dar este la un stadiu utilizabil.

Server-ul are doua componente:
- O componenta de server pur HTTP, prin care poti distribui un site
static la alegerea ta. Acest site este salvat direct in folder-ul
"static". Orice fisier HTML inserat acolo este adaugat la site.
- Componenta de API care comunica cu baza de date pentru VianuEdu.
Aceasta are niste functii predefinite care pot fi accesate in:
```
http://www.example.com/api/*
```

In curand, documentatia de la API va fi disponibila ca si functie in
server.

## Instalare

- Downloadeaza ultimul release al VianuEdu-Server
- Despacheteaza zip-ul
- Insereaza site-ul static, daca este dorit, in folder-ul "static"
- Schimba toate configurarile din folder-ul "config" dupa necesitati
- Porneste o instanta MongoDB care sa aibe username si parola din
configuratiile anterior completate (user care are acces doar la baza de
date pentru proiect), si sa aibe urmatoarea structura:
```
[dbName]
│
├───[MATERIE]Edu.Grades
│   ├───{ ... }
│   └───{ ... }
├───[MATERIE]Edu.Tests
│   ├───{ ... }
│   └───{ ... }
├───Students.Accounts
│   ├───{ ... }
│   └───{ ... }
├───Students.SubmittedAnswers
│   ├───{ ... }
│   └───{ ... }
├───Teachers.Accounts
│   ├───{ ... }
│   └───{ ... }
└───[dbName].TestList
    ├───{ ... }
    └───{ ... }
```
- Porneste server-ul folosind comanda corespunzatoare sistemului tau de
operare.

## Rulare dupa instalare

### Linux & Mac OS X:
```
root@localhost:~# ./VianuEdu-Server &
```
### Windows:
```
C:\Users\elev> START /B VianuEdu-Server.exe
```
Programele vor fi rulate in background.

## Functionalitati extra

- Log-urile pentru server se pot gasi in folder-ul "log" (compatibile
cu ELK Stack)
- Repozitoriu de lectii
- TLS