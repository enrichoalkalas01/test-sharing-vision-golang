# Uber Fx Dependency Injection - Implementation Guide

## Apa itu Uber Fx?

[Uber Fx](https://github.com/uber-go/fx) adalah framework dependency injection untuk Go. Fx menggantikan proses manual wiring (membuat dan menghubungkan dependency satu per satu) dengan sistem otomatis yang:

- **Auto-resolve dependencies** — cukup daftarkan constructor, Fx otomatis tahu urutan pembuatannya
- **Lifecycle management** — startup dan shutdown dikelola otomatis (graceful shutdown via SIGINT/SIGTERM)
- **Modular** — setiap concern dipisah ke module tersendiri

---

## Apa yang Berubah?

### File yang Dibuat (6 file baru)

```
internal/modules/
├── config_module.go      # Provider untuk *viper.Viper
├── logger_module.go      # Provider untuk *logger.Logger dan *zap.Logger
├── database_module.go    # Provider untuk *gorm.DB
├── article_module.go     # Provider untuk Repository, Usecase, Handler
├── server_module.go      # Provider untuk *fiber.App + lifecycle
└── route_module.go       # Registrasi middleware dan routes
```

### File yang Diubah (1 file)

- **`cmd/api/server.go`** — Manual wiring diganti dengan `fx.New(...).Run()`

### File yang TIDAK Berubah

Semua business logic tetap sama:

- `configs/viper.go`
- `pkg/logger/logger.go`
- `pkg/database/mysql.go`
- `internal/handler/article_handler.go`
- `internal/usecase/article/*`
- `internal/repository/article/*`
- `internal/domain/article.go`
- `pkg/middlewares/*`

---

## Sebelum vs Sesudah

### Sebelum (Manual Wiring)

```go
func main() {
    // Step 1 - buat logger manual
    log, err := logger.NewLogger(logger.Config{...})
    if err != nil { panic(err) }
    defer log.Sync()

    // Step 2 - buat config manual
    configEnv, err := configs.NewViper(".env", "env", ".", "../../")
    if err != nil { log.Fatal(...) }

    // Step 3 - buat database manual
    mysql, err := database.NewMySQL(configEnv, log.Logger)
    if err != nil { log.Fatal(...) }

    // Step 4 - buat server, setup middleware, routes, start
    server := servers.NewFiberServer(configEnv, log.Logger, mysql)
    server.SetupMiddlewares()
    server.SetupRoutes()
    server.Start()
}
```

**Masalah:**
- Urutan harus benar (logger → config → database → server)
- Error handling boilerplate di setiap langkah
- Graceful shutdown harus di-handle manual
- Menambah dependency baru = ubah `main()` + pastikan urutan benar

### Sesudah (Uber Fx)

```go
func main() {
    fx.New(
        modules.LoggerModule,
        modules.ConfigModule,
        modules.DatabaseModule,
        modules.ArticleModule,
        modules.ServerModule,
        modules.RouteModule,
    ).Run()
}
```

**Keuntungan:**
- Urutan otomatis (Fx resolve dependency graph sendiri)
- Error handling terpusat di Fx
- Graceful shutdown otomatis (SIGINT/SIGTERM)
- Menambah dependency baru = buat module baru, tambahkan ke list

---

## Cara Kerja Setiap Module

### 1. Config Module (`config_module.go`)

```go
var ConfigModule = fx.Module("config",
    fx.Provide(NewViper),  // mendaftarkan constructor
)

func NewViper() (*viper.Viper, error) {
    return configs.NewViper(".env", "env", ".", "../../")
}
```

- `fx.Provide` mendaftarkan fungsi sebagai **constructor**
- Fx melihat return type `*viper.Viper` dan tahu: "jika ada yang butuh `*viper.Viper`, panggil fungsi ini"
- Hanya dipanggil **sekali** (singleton)

### 2. Logger Module (`logger_module.go`)

```go
var LoggerModule = fx.Module("logger",
    fx.Provide(NewLogger),         // → *logger.Logger
    fx.Provide(ExtractZapLogger),  // → *zap.Logger (dari *logger.Logger)
    fx.Invoke(registerLoggerLifecycle),
)
```

- Menyediakan 2 type: `*logger.Logger` dan `*zap.Logger`
- `ExtractZapLogger` mengambil embedded `*zap.Logger` dari wrapper `*logger.Logger`
- Kenapa perlu keduanya? Karena semua constructor existing (`NewMySQL`, `NewArticleRepository`, dll) menggunakan `*zap.Logger`, bukan `*logger.Logger`
- **Lifecycle hook**: `Sync()` dipanggil saat shutdown untuk flush log buffer

### 3. Database Module (`database_module.go`)

```go
var DatabaseModule = fx.Module("database",
    fx.Provide(NewMySQL),                    // → *gorm.DB
    fx.Invoke(registerDatabaseLifecycle),
)

func NewMySQL(config *viper.Viper, log *zap.Logger) (*gorm.DB, error) {
    return database.NewMySQL(config, log)
}
```

- Fx melihat parameter `*viper.Viper` dan `*zap.Logger` → otomatis resolve dari ConfigModule dan LoggerModule
- **Lifecycle hook**: close database connection saat shutdown

### 4. Article Module (`article_module.go`)

```go
var ArticleModule = fx.Module("article",
    fx.Provide(repository.NewArticleRepository),  // → ArticleRepository
    fx.Provide(usecase.NewArticleUsecase),         // → ArticleUsecase
    fx.Provide(handler.NewArticleHandler),         // → *ArticleHandler
)
```

- Constructor didaftarkan langsung tanpa wrapper karena signature sudah kompatibel
- Fx otomatis resolve chain: `*gorm.DB` + `*zap.Logger` → `ArticleRepository` → `ArticleUsecase` → `*ArticleHandler`

### 5. Server Module (`server_module.go`)

```go
var ServerModule = fx.Module("server",
    fx.Provide(NewFiberApp),                   // → *fiber.App
    fx.Invoke(registerServerLifecycle),
)
```

- Membuat `*fiber.App` standalone (tidak pakai `FiberServer` wrapper)
- **Lifecycle hooks**:
  - `OnStart`: jalankan `app.Listen(address)` di goroutine
  - `OnStop`: panggil `app.ShutdownWithContext(ctx)` untuk graceful shutdown

### 6. Route Module (`route_module.go`)

```go
var RouteModule = fx.Module("route",
    fx.Invoke(registerMiddlewares),  // setup semua middleware
    fx.Invoke(registerRoutes),       // setup semua routes
)
```

- `fx.Invoke` berbeda dari `fx.Provide` — fungsi dipanggil untuk **side effect**, bukan menyediakan type baru
- Middleware order tetap sama: Recovery → CORS → Request ID → Fiber Logger → Zap Logger
- Route structure tetap sama: `/api/`, `/api/v1/`, `/api/v1/article/...`

---

## Dependency Graph

Fx otomatis me-resolve dependency dalam urutan yang benar:

```
*logger.Logger ──→ *zap.Logger
                        │
*viper.Viper ───────────┤
                        │
                   *gorm.DB
                        │
               ArticleRepository
                        │
               ArticleUsecase
                        │
               *ArticleHandler
                        │
*fiber.App ─────────────┤
        │               │
   middlewares      routes (registered)
        │
   server lifecycle (Listen / Shutdown)
```

Tidak peduli urutan module di `fx.New(...)` — Fx resolve berdasarkan type dependency, bukan urutan deklarasi.

---

## Lifecycle (Startup & Shutdown)

### Startup Order

Fx menjalankan `OnStart` hooks dalam urutan dependency (dependency pertama):

1. Logger dibuat
2. Config dibuat
3. Database dibuat (depends on config + logger)
4. Repository → Usecase → Handler dibuat
5. Fiber app dibuat
6. Middlewares diregister
7. Routes diregister
8. Server mulai listen (`OnStart` hook)

### Shutdown Order

Fx menjalankan `OnStop` hooks dalam **urutan terbalik**:

1. Server shutdown (stop accepting requests)
2. Database connection ditutup
3. Logger di-sync (flush buffer)

Shutdown di-trigger otomatis oleh:
- `SIGINT` (Ctrl+C)
- `SIGTERM` (docker stop, kubernetes termination)

---

## Cara Menambah Feature Baru

Contoh: menambah feature `User`

### 1. Buat domain, repository, usecase, handler seperti biasa

### 2. Buat module baru

```go
// internal/modules/user_module.go
var UserModule = fx.Module("user",
    fx.Provide(userRepo.NewUserRepository),
    fx.Provide(userUsecase.NewUserUsecase),
    fx.Provide(userHandler.NewUserHandler),
)
```

### 3. Tambahkan routes di `route_module.go`

```go
func registerRoutes(
    app *fiber.App,
    articleHandler *handler.ArticleHandler,
    userHandler *handler.UserHandler,  // tambah parameter
    config *viper.Viper,
    log *zap.Logger,
) {
    // ... existing routes ...

    // User routes
    users := apiV1.Group("/user")
    users.Get("/", userHandler.GetList)
    // ...
}
```

### 4. Daftarkan di `cmd/api/server.go`

```go
fx.New(
    modules.LoggerModule,
    modules.ConfigModule,
    modules.DatabaseModule,
    modules.ArticleModule,
    modules.UserModule,     // tambah ini
    modules.ServerModule,
    modules.RouteModule,
)
```

Selesai. Tidak perlu mengubah urutan, tidak perlu handle error manual.

---

## Konsep Penting Fx

| Konsep | Penjelasan |
|--------|------------|
| `fx.Provide` | Mendaftarkan constructor. Fx panggil saat ada yang butuh return type-nya. Singleton — hanya dipanggil sekali. |
| `fx.Invoke` | Mendaftarkan fungsi yang dipanggil untuk side effect (register routes, setup middleware). Selalu dijalankan saat startup. |
| `fx.Module` | Mengelompokkan `Provide` dan `Invoke` dalam satu unit logis. Purely organizational. |
| `fx.Lifecycle` | Hook untuk `OnStart` dan `OnStop`. Digunakan untuk start server, close connections, dll. |
| `fx.New(...).Run()` | Membuat container, resolve semua dependency, jalankan lifecycle, dan block sampai shutdown signal diterima. |

---

## Troubleshooting

### "missing type" error saat startup
Fx akan menampilkan error jelas jika ada dependency yang belum terdaftar. Contoh:
```
fx.Provide(handler.NewArticleHandler) needs usecase.ArticleUsecase
```
Solusi: pastikan module yang menyediakan type tersebut sudah didaftarkan.

### "already provided" error
Dua module menyediakan type yang sama. Solusi: gunakan `fx.Annotate` dengan `fx.ResultTags` untuk membedakan.

### Server tidak graceful shutdown
Pastikan menggunakan `.Run()` (blocking), bukan `.Start()` (non-blocking). `.Run()` otomatis menunggu shutdown signal.
