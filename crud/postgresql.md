# How to avoid defualt running of database in background

- Since you're developing your Go backend, a good setup is to manually start PostgreSQL when you work and stop it when you're done.

## Option 1 — Start/stop manually from CMD

```sql
sc config postgresql-x64-18 start= demand ---> (now it will start/stop on demand)
net start postgresql-x64-18 ---> start work
net stop postgresql-x64-18 --->  stop work
```

## Option 2 — From Powershell

```sql
Start-Service postgresql-x64-18
Stop-Service postgresql-x64-18
```

## Check version By CMD

```sql
psql --version
psql (PostgreSQL) 18.6
```

## Connect to PostgreSQL By CMD
```sql
psql -U postgres
```

## 