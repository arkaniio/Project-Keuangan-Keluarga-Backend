package main

import (
    "fmt"
    "log"
    "project-keuangan-keluarga/config"
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load()
    dbCfg := config.DefaultDatabaseConfig()
    db, err := config.InitDB(dbCfg)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var userCount int
    db.Get(&userCount, "SELECT COUNT(*) FROM users")
    fmt.Printf("Users count: %d\n", userCount)

    var familyCount int
    db.Get(&familyCount, "SELECT COUNT(*) FROM families")
    fmt.Printf("Families count: %d\n", familyCount)

    var memberCount int
    db.Get(&memberCount, "SELECT COUNT(*) FROM family_members")
    fmt.Printf("Family Members count: %d\n", memberCount)

    fmt.Println("\nDetail Members:")
    var members []struct {
        Username string `db:"username"`
        FamilyId string `db:"family_id"`
        Role     string `db:"role"`
    }
    db.Select(&members, "SELECT u.username, fm.family_id::text, fm.role FROM family_members fm JOIN users u ON fm.user_id = u.id")
    for _, m := range members {
        fmt.Printf("User: %s, Family: %s, Role: %s\n", m.Username, m.FamilyId, m.Role)
    }

    fmt.Println("\nDetail Families:")
    var families []struct {
        Name string `db:"name"`
        Id   string `db:"id"`
    }
    db.Select(&families, "SELECT name, id::text FROM families")
    for _, f := range families {
        fmt.Printf("Family: %s, ID: %s\n", f.Name, f.Id)
    }
}
