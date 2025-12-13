env "local" {
  src = ["./database/schema"]
  url = "mysql://root:password@localhost:3306/app-db"
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://./database/migrations"
  }
}
